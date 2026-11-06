package dproxy

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Darren5211314/go_common/internal/dlog"
)

var (
	errNotStarted   = errors.New("代理池未启动")
	errClosed   = errors.New("代理池已关闭")
	errGetTimeout = errors.New("获取代理超时")
	errNoLink = errors.New("未设置代理链接")
)

// ProxyStats 代理池运行统计。
type ProxyStats struct {
	LinkCount int64
	VerifyOK   int64
	VerifyFail   int64
	Verifying   int64
	PoolCount   int64
	GetCount     int64
	Expired       int64
}

type poolEntry struct {
	proxy   string
	expireAt time.Time // 零值表示不过期
}

// ProxyPool 有状态代理池。
type ProxyPool struct {
	mu   sync.Mutex
	cond  *sync.Cond
	closeCh chan struct{}
	refillCh chan struct{}
	verifyCh chan string

	started bool
	closed bool

	proxyLink    string
	capacity      int
	verifyEnabled    bool
	verifyTimeout    time.Duration
	ttlSec    time.Duration
	verifyWorkers   int
	fetchInterval    time.Duration
	verifyURL    string
	defaultScheme    string
	httpClient *http.Client

	entries []poolEntry
	seen   map[string]struct{} // 池中或验证中

	linkCount  atomic.Int64
	VerifyOK   atomic.Int64
	VerifyFail atomic.Int64
	Verifying  atomic.Int64
	getCount   atomic.Int64
	Expired    atomic.Int64
	logEnabled atomic.Bool

	wg sync.WaitGroup
}

// newPool 创建未启动的代理池（带默认配置）。
func newPool() *ProxyPool {
	p := &ProxyPool{
		closeCh:     make(chan struct{}),
		refillCh:     make(chan struct{}, 1),
		capacity:       20,
		verifyEnabled:   true,
		verifyTimeout:   5 * time.Second,
		ttlSec:   0,
		verifyWorkers: 5,
		fetchInterval:   3 * time.Second,
		verifyURL:   "https://httpbin.org/ip",
		defaultScheme:   "http",
		seen:       make(map[string]struct{}),
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
	p.logEnabled.Store(true)
	p.cond = sync.NewCond(&p.mu)
	return p
}

// SetLog 控制是否输出拉取/验证/取出等运行日志；默认 true。
func (p *ProxyPool) SetLog(enable bool) {
	p.logEnabled.Store(enable)
}

func (p *ProxyPool) logInfow(msg string, keysAndValues ...any) {
	if !p.logEnabled.Load() {
		return
	}
	dlog.Sugar().Infow(msg, keysAndValues...)
}

func (p *ProxyPool) logWarnw(msg string, keysAndValues ...any) {
	if !p.logEnabled.Load() {
		return
	}
	dlog.Sugar().Warnw(msg, keysAndValues...)
}

func (p *ProxyPool) SetProxy(link string) {
	p.mu.Lock()
	p.proxyLink = strings.TrimSpace(link)
	p.mu.Unlock()
}

func (p *ProxyPool) SetPoolNum(capacity int) {
	if capacity < 1 {
		capacity = 1
	}
	p.mu.Lock()
	p.capacity = capacity
	p.mu.Unlock()
	p.notifyRefill()
}

func (p *ProxyPool) SetVerify(verify bool) {
	p.mu.Lock()
	p.verifyEnabled = verify
	p.mu.Unlock()
}

func (p *ProxyPool) SetVerifyTimeout(sec int) {
	if sec <= 0 {
		sec = 5
	}
	p.mu.Lock()
	p.verifyTimeout = time.Duration(sec) * time.Second
	p.mu.Unlock()
}

func (p *ProxyPool) SetTTL(sec int) {
	p.mu.Lock()
	if sec <= 0 {
		p.ttlSec = 0
	} else {
		p.ttlSec = time.Duration(sec) * time.Second
	}
	p.mu.Unlock()
}

func (p *ProxyPool) SetVerifyWorkers(count int) {
	if count < 1 {
		count = 1
	}
	p.mu.Lock()
	p.verifyWorkers = count
	p.mu.Unlock()
}

func (p *ProxyPool) SetGetInterval(sec int) {
	if sec <= 0 {
		sec = 3
	}
	p.mu.Lock()
	p.fetchInterval = time.Duration(sec) * time.Second
	p.mu.Unlock()
}

func (p *ProxyPool) SetVerifyURL(url string) {
	p.mu.Lock()
	p.verifyURL = strings.TrimSpace(url)
	p.mu.Unlock()
}

func (p *ProxyPool) SetProxyType(scheme string) {
	scheme, err := normalizeScheme(scheme)
	if err != nil {
		scheme = "http"
	}
	p.mu.Lock()
	p.defaultScheme = scheme
	p.mu.Unlock()
}

func (p *ProxyPool) notifyRefill() {
	select {
	case p.refillCh <- struct{}{}:
	default:
	}
}

// Start 启动拉取与验证协程。
func (p *ProxyPool) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return errClosed
	}
	if p.started {
		return nil
	}
	if p.proxyLink == "" {
		return errNoLink
	}
	workers := p.verifyWorkers
	p.verifyCh = make(chan string, workers*2)
	p.started = true

	p.wg.Add(1)
	go p.fetchLoop()

	for i := 0; i < workers; i++ {
		p.wg.Add(1)
		go p.verifyLoop()
	}
	p.notifyRefill()
	return nil
}

// Close 停止工作协程并唤醒等待中的获取。
func (p *ProxyPool) Close() {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.closed = true
	close(p.closeCh)
	p.cond.Broadcast()
	p.mu.Unlock()
	p.wg.Wait()
}

func (p *ProxyPool) purgeExpiredLocked(now time.Time) {
	newList := p.entries[:0]
	for _, item := range p.entries {
		if !item.expireAt.IsZero() && !now.Before(item.expireAt) {
			delete(p.seen, proxyKey(item.proxy))
			p.Expired.Add(1)
			continue
		}
		newList = append(newList, item)
	}
	p.entries = newList
}

func (p *ProxyPool) availableLocked() int {
	return len(p.entries)
}

func (p *ProxyPool) needRefillLocked() bool {
	return p.availableLocked()+int(p.Verifying.Load()) < p.capacity
}

// GetProxy 取出一个代理；waitTime=-1 一直等，>=0 为超时秒数。
func (p *ProxyPool) GetProxy(waitTime int) (string, error) {
	deadline := time.Time{}
	if waitTime >= 0 {
		deadline = time.Now().Add(time.Duration(waitTime) * time.Second)
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.started {
		return "", errNotStarted
	}

	for {
		if p.closed {
			return "", errClosed
		}
		p.purgeExpiredLocked(time.Now())
		if len(p.entries) > 0 {
			item := p.entries[0]
			p.entries = p.entries[1:]
			delete(p.seen, proxyKey(item.proxy))
			taken := p.getCount.Add(1)
			remain := len(p.entries)
			if p.needRefillLocked() {
				p.notifyRefill()
			}
			p.logInfow("取出代理", "proxy", item.proxy, "已取出", taken, "池剩余", remain)
			return item.proxy, nil
		}
		p.notifyRefill()
		if waitTime >= 0 && !deadline.IsZero() && !time.Now().Before(deadline) {
			return "", errGetTimeout
		}
		if waitTime >= 0 && !deadline.IsZero() {
			remain := time.Until(deadline)
			if remain <= 0 {
				return "", errGetTimeout
			}
			timer := time.AfterFunc(remain, func() {
				p.mu.Lock()
				p.cond.Broadcast()
				p.mu.Unlock()
			})
			p.cond.Wait()
			timer.Stop()
			continue
		}
		p.cond.Wait()
	}
}

func (p *ProxyPool) Stats() ProxyStats {
	p.mu.Lock()
	p.purgeExpiredLocked(time.Now())
	ok := int64(len(p.entries))
	p.mu.Unlock()
	return ProxyStats{
		LinkCount:  p.linkCount.Load(),
		VerifyOK:   p.VerifyOK.Load(),
		VerifyFail: p.VerifyFail.Load(),
		Verifying:  p.Verifying.Load(),
		PoolCount:  ok,
		GetCount:   p.getCount.Load(),
		Expired:    p.Expired.Load(),
	}
}

func (p *ProxyPool) fetchLoop() {
	defer p.wg.Done()
	failStreak := 0
	var lastFetch time.Time

	for {
		select {
		case <-p.closeCh:
			return
		case <-p.refillCh:
		}

		for {
			select {
			case <-p.closeCh:
				return
			default:
			}

			p.mu.Lock()
			need := p.needRefillLocked()
			link := p.proxyLink
			interval := p.fetchInterval
			defaultScheme := p.defaultScheme
			needVerify := p.verifyEnabled
			p.mu.Unlock()

			if !need {
				break
			}

			now := time.Now()
			if !lastFetch.IsZero() {
				wait := interval - now.Sub(lastFetch)
				if wait > 0 {
					choice := time.NewTimer(wait)
					select {
					case <-p.closeCh:
						choice.Stop()
						return
					case <-choice.C:
					}
				}
			}

			err := p.fetchOnce(link, defaultScheme, needVerify)
			lastFetch = time.Now()
			if err != nil {
				failStreak++
				p.logWarnw("从链接拉取代理失败", "err", err)
				backoff := time.Duration(1<<minInt(failStreak-1, 5)) * time.Second
				if backoff > 60*time.Second {
					backoff = 60 * time.Second
				}
				choice := time.NewTimer(backoff)
				select {
				case <-p.closeCh:
					choice.Stop()
					return
				case <-choice.C:
				}
				continue
			}
			failStreak = 0
		}
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (p *ProxyPool) fetchOnce(link, defaultScheme string, needVerify bool) error {
	resp, err := p.httpClient.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New("代理链接返回错误状态码")
	}

	added := 0
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		proxy, err := parseProxyLine(line, defaultScheme)
		if err != nil {
			continue
		}
		k := proxyKey(proxy)
		p.mu.Lock()
		if _, exists := p.seen[k]; exists || p.closed {
			p.mu.Unlock()
			continue
		}
		if !p.needRefillLocked() {
			p.mu.Unlock()
			break
		}
		p.seen[k] = struct{}{}
		p.linkCount.Add(1)
		added++
		if !needVerify {
			p.enqueueLocked(proxy)
			p.mu.Unlock()
			continue
		}
		p.Verifying.Add(1)
		p.mu.Unlock()

		select {
		case <-p.closeCh:
			p.Verifying.Add(-1)
			p.mu.Lock()
			delete(p.seen, k)
			p.mu.Unlock()
			return errClosed
		case p.verifyCh <- proxy:
		}
	}
	if added > 0 {
		p.mu.Lock()
		avail := len(p.entries)
		p.mu.Unlock()
		p.logInfow("从链接获取代理",
			"本轮新增", added,
			"验证", needVerify,
			"池可用", avail,
			"正在验证", p.Verifying.Load(),
		)
	}
	return nil
}

func (p *ProxyPool) enqueueLocked(proxy string) {
	expired := time.Time{}
	if p.ttlSec > 0 {
		expired = time.Now().Add(p.ttlSec)
	}
	p.entries = append(p.entries, poolEntry{proxy: proxy, expireAt: expired})
	p.cond.Broadcast()
}

func (p *ProxyPool) verifyLoop() {
	defer p.wg.Done()
	for {
		select {
		case <-p.closeCh:
			return
		case proxy, ok := <-p.verifyCh:
			if !ok {
				return
			}
			p.mu.Lock()
			timeout := p.verifyTimeout
			url := p.verifyURL
			alive := p.ttlSec
			closed := p.closed
			p.mu.Unlock()
			if closed {
				p.Verifying.Add(-1)
				p.mu.Lock()
				delete(p.seen, proxyKey(proxy))
				p.mu.Unlock()
				continue
			}

			p.logInfow("正在验证代理", "proxy", proxy)
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err := verifyProxy(ctx, proxy, url, timeout)
			cancel()
			p.Verifying.Add(-1)

			p.mu.Lock()
			if err != nil || p.closed {
				delete(p.seen, proxyKey(proxy))
				if err != nil {
					p.VerifyFail.Add(1)
				}
				if p.needRefillLocked() {
					p.notifyRefill()
				}
				p.mu.Unlock()
				if err != nil {
					p.logWarnw("代理验证失败", "proxy", proxy, "err", err)
				}
				continue
			}
			expired := time.Time{}
			if alive > 0 {
				expired = time.Now().Add(alive)
			}
			p.entries = append(p.entries, poolEntry{proxy: proxy, expireAt: expired})
			p.VerifyOK.Add(1)
			avail := len(p.entries)
			p.cond.Broadcast()
			if p.needRefillLocked() {
				p.notifyRefill()
			}
			p.mu.Unlock()
			p.logInfow("代理验证成功", "proxy", proxy, "池可用", avail)
		}
	}
}
