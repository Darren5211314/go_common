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
)

var (
	错误未启动   = errors.New("代理池未启动")
	错误已关闭   = errors.New("代理池已关闭")
	错误获取超时 = errors.New("获取代理超时")
	错误未配置链接 = errors.New("未设置代理链接")
)

// D代理统计 代理池运行统计。
type D代理统计 struct {
	D从链接获取数 int64
	D验证成功数   int64
	D验证失败数   int64
	D正在验证数   int64
	D池中可用数   int64
	D已取出数     int64
	D失效数       int64
}

type 池条目 struct {
	代理   string
	过期时间 time.Time // 零值表示不过期
}

// D代理池 有状态代理池。
type D代理池 struct {
	锁   sync.Mutex
	条件  *sync.Cond
	关闭通 chan struct{}
	补充通 chan struct{}
	验证通 chan string

	已启动 bool
	已关闭 bool

	代理链接    string
	容量      int
	是否验证    bool
	验证超时    time.Duration
	存活时间    time.Duration
	验证线程数   int
	获取间隔    time.Duration
	验证地址    string
	默认协议    string
	http客户端 *http.Client

	条目列表 []池条目
	已知   map[string]struct{} // 池中或验证中

	从链接获取数 atomic.Int64
	验证成功数   atomic.Int64
	验证失败数   atomic.Int64
	正在验证数   atomic.Int64
	已取出数     atomic.Int64
	失效数       atomic.Int64

	等待组 sync.WaitGroup
}

// D新建 创建未启动的代理池（带默认配置）。
func 新建() *D代理池 {
	p := &D代理池{
		关闭通:     make(chan struct{}),
		补充通:     make(chan struct{}, 1),
		容量:       20,
		是否验证:   true,
		验证超时:   5 * time.Second,
		存活时间:   0,
		验证线程数: 5,
		获取间隔:   3 * time.Second,
		验证地址:   "https://httpbin.org/ip",
		默认协议:   "http",
		已知:       make(map[string]struct{}),
		http客户端: &http.Client{Timeout: 15 * time.Second},
	}
	p.条件 = sync.NewCond(&p.锁)
	return p
}

func (p *D代理池) D设置代理链接(链接 string) {
	p.锁.Lock()
	p.代理链接 = strings.TrimSpace(链接)
	p.锁.Unlock()
}

func (p *D代理池) D设置代理池容量(容量 int) {
	if 容量 < 1 {
		容量 = 1
	}
	p.锁.Lock()
	p.容量 = 容量
	p.锁.Unlock()
	p.通知补充()
}

func (p *D代理池) D设置是否验证(验证 bool) {
	p.锁.Lock()
	p.是否验证 = 验证
	p.锁.Unlock()
}

func (p *D代理池) D设置验证超时(秒 int) {
	if 秒 <= 0 {
		秒 = 5
	}
	p.锁.Lock()
	p.验证超时 = time.Duration(秒) * time.Second
	p.锁.Unlock()
}

func (p *D代理池) D设置存活时间(秒 int) {
	p.锁.Lock()
	if 秒 <= 0 {
		p.存活时间 = 0
	} else {
		p.存活时间 = time.Duration(秒) * time.Second
	}
	p.锁.Unlock()
}

func (p *D代理池) D设置验证线程数(数量 int) {
	if 数量 < 1 {
		数量 = 1
	}
	p.锁.Lock()
	p.验证线程数 = 数量
	p.锁.Unlock()
}

func (p *D代理池) D设置获取间隔(秒 int) {
	if 秒 <= 0 {
		秒 = 3
	}
	p.锁.Lock()
	p.获取间隔 = time.Duration(秒) * time.Second
	p.锁.Unlock()
}

func (p *D代理池) D设置验证地址(地址 string) {
	p.锁.Lock()
	p.验证地址 = strings.TrimSpace(地址)
	p.锁.Unlock()
}

func (p *D代理池) D设置默认协议(协议 string) {
	协议, 错误 := 规范化协议(协议)
	if 错误 != nil {
		协议 = "http"
	}
	p.锁.Lock()
	p.默认协议 = 协议
	p.锁.Unlock()
}

func (p *D代理池) 通知补充() {
	select {
	case p.补充通 <- struct{}{}:
	default:
	}
}

// D启动 启动拉取与验证协程。
func (p *D代理池) D启动() error {
	p.锁.Lock()
	defer p.锁.Unlock()
	if p.已关闭 {
		return 错误已关闭
	}
	if p.已启动 {
		return nil
	}
	if p.代理链接 == "" {
		return 错误未配置链接
	}
	线程数 := p.验证线程数
	p.验证通 = make(chan string, 线程数*2)
	p.已启动 = true

	p.等待组.Add(1)
	go p.拉取循环()

	for i := 0; i < 线程数; i++ {
		p.等待组.Add(1)
		go p.验证循环()
	}
	p.通知补充()
	return nil
}

// D关闭 停止工作协程并唤醒等待中的获取。
func (p *D代理池) D关闭() {
	p.锁.Lock()
	if p.已关闭 {
		p.锁.Unlock()
		return
	}
	p.已关闭 = true
	close(p.关闭通)
	p.条件.Broadcast()
	p.锁.Unlock()
	p.等待组.Wait()
}

func (p *D代理池) 剔除过期锁定(现在 time.Time) {
	新列表 := p.条目列表[:0]
	for _, 项 := range p.条目列表 {
		if !项.过期时间.IsZero() && !现在.Before(项.过期时间) {
			delete(p.已知, 代理键(项.代理))
			p.失效数.Add(1)
			continue
		}
		新列表 = append(新列表, 项)
	}
	p.条目列表 = 新列表
}

func (p *D代理池) 可用数锁定() int {
	return len(p.条目列表)
}

func (p *D代理池) 需要补充锁定() bool {
	return p.可用数锁定()+int(p.正在验证数.Load()) < p.容量
}

// D获取代理 取出一个代理；等待秒=-1 一直等，>=0 为超时秒数。
func (p *D代理池) D获取代理(等待秒 int) (string, error) {
	截止 := time.Time{}
	if 等待秒 >= 0 {
		截止 = time.Now().Add(time.Duration(等待秒) * time.Second)
	}

	p.锁.Lock()
	defer p.锁.Unlock()
	if !p.已启动 {
		return "", 错误未启动
	}

	for {
		if p.已关闭 {
			return "", 错误已关闭
		}
		p.剔除过期锁定(time.Now())
		if len(p.条目列表) > 0 {
			项 := p.条目列表[0]
			p.条目列表 = p.条目列表[1:]
			delete(p.已知, 代理键(项.代理))
			p.已取出数.Add(1)
			if p.需要补充锁定() {
				p.通知补充()
			}
			return 项.代理, nil
		}
		p.通知补充()
		if 等待秒 >= 0 && !截止.IsZero() && !time.Now().Before(截止) {
			return "", 错误获取超时
		}
		if 等待秒 >= 0 && !截止.IsZero() {
			剩余 := time.Until(截止)
			if 剩余 <= 0 {
				return "", 错误获取超时
			}
			计时器 := time.AfterFunc(剩余, func() {
				p.锁.Lock()
				p.条件.Broadcast()
				p.锁.Unlock()
			})
			p.条件.Wait()
			计时器.Stop()
			continue
		}
		p.条件.Wait()
	}
}

func (p *D代理池) D取统计() D代理统计 {
	p.锁.Lock()
	p.剔除过期锁定(time.Now())
	可用 := int64(len(p.条目列表))
	p.锁.Unlock()
	return D代理统计{
		D从链接获取数: p.从链接获取数.Load(),
		D验证成功数:   p.验证成功数.Load(),
		D验证失败数:   p.验证失败数.Load(),
		D正在验证数:   p.正在验证数.Load(),
		D池中可用数:   可用,
		D已取出数:     p.已取出数.Load(),
		D失效数:       p.失效数.Load(),
	}
}

func (p *D代理池) 拉取循环() {
	defer p.等待组.Done()
	连续失败 := 0
	var 上次拉取 time.Time

	for {
		select {
		case <-p.关闭通:
			return
		case <-p.补充通:
		}

		for {
			select {
			case <-p.关闭通:
				return
			default:
			}

			p.锁.Lock()
			需要 := p.需要补充锁定()
			链接 := p.代理链接
			间隔 := p.获取间隔
			默认协议 := p.默认协议
			要验证 := p.是否验证
			p.锁.Unlock()

			if !需要 {
				break
			}

			现在 := time.Now()
			if !上次拉取.IsZero() {
				等待 := 间隔 - 现在.Sub(上次拉取)
				if 等待 > 0 {
					选择 := time.NewTimer(等待)
					select {
					case <-p.关闭通:
						选择.Stop()
						return
					case <-选择.C:
					}
				}
			}

			错误 := p.执行一次拉取(链接, 默认协议, 要验证)
			上次拉取 = time.Now()
			if 错误 != nil {
				连续失败++
				退避 := time.Duration(1<<minInt(连续失败-1, 5)) * time.Second
				if 退避 > 60*time.Second {
					退避 = 60 * time.Second
				}
				选择 := time.NewTimer(退避)
				select {
				case <-p.关闭通:
					选择.Stop()
					return
				case <-选择.C:
				}
				continue
			}
			连续失败 = 0
		}
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (p *D代理池) 执行一次拉取(链接, 默认协议 string, 要验证 bool) error {
	响应, 错误 := p.http客户端.Get(链接)
	if 错误 != nil {
		return 错误
	}
	defer 响应.Body.Close()
	正文, 错误 := io.ReadAll(io.LimitReader(响应.Body, 1<<20))
	if 错误 != nil {
		return 错误
	}
	if 响应.StatusCode >= 400 {
		return errors.New("代理链接返回错误状态码")
	}

	行列表 := strings.Split(string(正文), "\n")
	for _, 行 := range 行列表 {
		代理, 错误 := 解析代理行(行, 默认协议)
		if 错误 != nil {
			continue
		}
		键 := 代理键(代理)
		p.锁.Lock()
		if _, 存在 := p.已知[键]; 存在 || p.已关闭 {
			p.锁.Unlock()
			continue
		}
		if !p.需要补充锁定() {
			p.锁.Unlock()
			break
		}
		p.已知[键] = struct{}{}
		p.从链接获取数.Add(1)
		if !要验证 {
			p.入池锁定(代理)
			p.锁.Unlock()
			continue
		}
		p.正在验证数.Add(1)
		p.锁.Unlock()

		select {
		case <-p.关闭通:
			p.正在验证数.Add(-1)
			p.锁.Lock()
			delete(p.已知, 键)
			p.锁.Unlock()
			return 错误已关闭
		case p.验证通 <- 代理:
		}
	}
	return nil
}

func (p *D代理池) 入池锁定(代理 string) {
	过期 := time.Time{}
	if p.存活时间 > 0 {
		过期 = time.Now().Add(p.存活时间)
	}
	p.条目列表 = append(p.条目列表, 池条目{代理: 代理, 过期时间: 过期})
	p.条件.Broadcast()
}

func (p *D代理池) 验证循环() {
	defer p.等待组.Done()
	for {
		select {
		case <-p.关闭通:
			return
		case 代理, 成功 := <-p.验证通:
			if !成功 {
				return
			}
			p.锁.Lock()
			超时 := p.验证超时
			地址 := p.验证地址
			存活 := p.存活时间
			已关闭 := p.已关闭
			p.锁.Unlock()
			if 已关闭 {
				p.正在验证数.Add(-1)
				p.锁.Lock()
				delete(p.已知, 代理键(代理))
				p.锁.Unlock()
				continue
			}

			ctx, 取消 := context.WithTimeout(context.Background(), 超时)
			错误 := 验证代理(ctx, 代理, 地址, 超时)
			取消()
			p.正在验证数.Add(-1)

			p.锁.Lock()
			if 错误 != nil || p.已关闭 {
				delete(p.已知, 代理键(代理))
				if 错误 != nil {
					p.验证失败数.Add(1)
				}
				if p.需要补充锁定() {
					p.通知补充()
				}
				p.锁.Unlock()
				continue
			}
			过期 := time.Time{}
			if 存活 > 0 {
				过期 = time.Now().Add(存活)
			}
			p.条目列表 = append(p.条目列表, 池条目{代理: 代理, 过期时间: 过期})
			p.验证成功数.Add(1)
			p.条件.Broadcast()
			if p.需要补充锁定() {
				p.通知补充()
			}
			p.锁.Unlock()
		}
	}
}
