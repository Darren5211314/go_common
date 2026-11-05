package dthread

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var (
	错误未启动   = errors.New("线程池未启动")
	错误已停止   = errors.New("线程池已停止")
	错误任务空   = errors.New("任务为空")
	错误等待超时 = errors.New("等待任务超时")
	错误任务取消 = errors.New("任务已取消")
)

// D线程池统计 瞬时与累计统计。
type D线程池统计 struct {
	D最大线程数 int
	D执行中数量 int
	D空闲数量   int
	D排队数量   int
	D已完成数   int64
	D失败数     int64
	D已投递数   int64
}

// D任务 可等待结果的任务句柄。
type D任务 struct {
	完成   chan struct{}
	结果   any
	错误   error
	完成一次 sync.Once
}

func (t *D任务) 完成结果(结果 any, 错误 error) {
	t.完成一次.Do(func() {
		t.结果 = 结果
		t.错误 = 错误
		close(t.完成)
	})
}

// D等待 阻塞直到任务完成。
func (t *D任务) D等待() (any, error) {
	if t == nil {
		return nil, 错误任务空
	}
	<-t.完成
	return t.结果, t.错误
}

// D等待超时 最多等待秒数；超时返回错误。
func (t *D任务) D等待超时(秒 int) (any, error) {
	if t == nil {
		return nil, 错误任务空
	}
	if 秒 < 0 {
		return t.D等待()
	}
	选择 := time.NewTimer(time.Duration(秒) * time.Second)
	defer 选择.Stop()
	select {
	case <-t.完成:
		return t.结果, t.错误
	case <-选择.C:
		return nil, 错误等待超时
	}
}

type 内部任务 struct {
	无返回 func()
	有返回 func() (any, error)
	句柄   *D任务
}

// D线程池 可伸缩 worker 池。
type D线程池 struct {
	锁 sync.Mutex

	最大线程数 int
	存活工人数 int
	执行中   int
	队列容量 int

	队列   chan *内部任务
	退出通 chan struct{}

	已启动 bool
	已停止 bool
	已暂停 bool
	强制  bool

	暂停条件 *sync.Cond

	投递上下文 context.Context
	投递取消  context.CancelFunc
	任务上下文 context.Context
	任务取消  context.CancelFunc

	工人等待 sync.WaitGroup

	已投递数 atomic.Int64
	已完成数 atomic.Int64
	失败数   atomic.Int64
}

// D新建 创建线程池；最大线程数 <=0 时为 1。默认队列容量 = 最大*4。
func 新建(最大线程数 int) *D线程池 {
	if 最大线程数 <= 0 {
		最大线程数 = 1
	}
	投递上下文, 投递取消 := context.WithCancel(context.Background())
	任务上下文, 任务取消 := context.WithCancel(context.Background())
	p := &D线程池{
		最大线程数: 最大线程数,
		队列容量: 最大线程数 * 4,
		投递上下文: 投递上下文,
		投递取消:  投递取消,
		任务上下文: 任务上下文,
		任务取消:  任务取消,
	}
	p.暂停条件 = sync.NewCond(&p.锁)
	return p
}

// D设置队列容量 仅启动前有效。
func (p *D线程池) D设置队列容量(容量 int) {
	p.锁.Lock()
	defer p.锁.Unlock()
	if p.已启动 {
		return
	}
	if 容量 < 1 {
		容量 = 1
	}
	p.队列容量 = 容量
}

// D任务上下文 返回任务 context（仅强制停止时 cancel）。
func (p *D线程池) D任务上下文() context.Context {
	p.锁.Lock()
	defer p.锁.Unlock()
	return p.任务上下文
}

// D启动 启动 worker。
func (p *D线程池) D启动() error {
	p.锁.Lock()
	defer p.锁.Unlock()
	if p.已停止 {
		return 错误已停止
	}
	if p.已启动 {
		return nil
	}
	p.队列 = make(chan *内部任务, p.队列容量)
	p.退出通 = make(chan struct{})
	p.已启动 = true
	for i := 0; i < p.最大线程数; i++ {
		p.启动工人锁定()
	}
	return nil
}

func (p *D线程池) 启动工人锁定() {
	p.存活工人数++
	p.工人等待.Add(1)
	go p.工人循环()
}

// D设置最大线程数 动态扩缩。
func (p *D线程池) D设置最大线程数(数量 int) {
	if 数量 <= 0 {
		数量 = 1
	}
	p.锁.Lock()
	defer p.锁.Unlock()
	p.最大线程数 = 数量
	if !p.已启动 || p.已停止 {
		return
	}
	for p.存活工人数 < p.最大线程数 {
		p.启动工人锁定()
	}
	p.暂停条件.Broadcast()
}

func (p *D线程池) D暂停() {
	p.锁.Lock()
	p.已暂停 = true
	p.锁.Unlock()
}

func (p *D线程池) D继续() {
	p.锁.Lock()
	p.已暂停 = false
	p.暂停条件.Broadcast()
	p.锁.Unlock()
}

func (p *D线程池) 可投递锁定() error {
	if !p.已启动 {
		return 错误未启动
	}
	if p.已停止 {
		return 错误已停止
	}
	return nil
}

func (p *D线程池) 投递内部(项 *内部任务) error {
	p.锁.Lock()
	if 错误 := p.可投递锁定(); 错误 != nil {
		p.锁.Unlock()
		return 错误
	}
	队列 := p.队列
	投递上下文 := p.投递上下文
	p.锁.Unlock()

	select {
	case 队列 <- 项:
		p.已投递数.Add(1)
		return nil
	case <-投递上下文.Done():
		return 错误已停止
	}
}

// D投递 投递无返回值任务；队列满时阻塞。
func (p *D线程池) D投递(任务 func()) error {
	if 任务 == nil {
		return 错误任务空
	}
	return p.投递内部(&内部任务{无返回: 任务})
}

// D投递任务 投递有返回值任务，返回可等待句柄。
func (p *D线程池) D投递任务(任务 func() (any, error)) (*D任务, error) {
	if 任务 == nil {
		return nil, 错误任务空
	}
	句柄 := &D任务{完成: make(chan struct{})}
	if 错误 := p.投递内部(&内部任务{有返回: 任务, 句柄: 句柄}); 错误 != nil {
		句柄.完成结果(nil, 错误)
		return 句柄, 错误
	}
	return 句柄, nil
}

func (p *D线程池) D取执行中数量() int {
	p.锁.Lock()
	defer p.锁.Unlock()
	return p.执行中
}

func (p *D线程池) D取空闲数量() int {
	p.锁.Lock()
	defer p.锁.Unlock()
	空闲 := p.存活工人数 - p.执行中
	if 空闲 < 0 {
		return 0
	}
	return 空闲
}

func (p *D线程池) D取排队数量() int {
	p.锁.Lock()
	队列 := p.队列
	p.锁.Unlock()
	if 队列 == nil {
		return 0
	}
	return len(队列)
}

func (p *D线程池) D取统计() D线程池统计 {
	p.锁.Lock()
	最大 := p.最大线程数
	执行中 := p.执行中
	存活 := p.存活工人数
	队列 := p.队列
	p.锁.Unlock()
	排队 := 0
	if 队列 != nil {
		排队 = len(队列)
	}
	空闲 := 存活 - 执行中
	if 空闲 < 0 {
		空闲 = 0
	}
	return D线程池统计{
		D最大线程数: 最大,
		D执行中数量: 执行中,
		D空闲数量:   空闲,
		D排队数量:   排队,
		D已完成数:   p.已完成数.Load(),
		D失败数:     p.失败数.Load(),
		D已投递数:   p.已投递数.Load(),
	}
}

func (p *D线程池) 等待排空() {
	for {
		p.锁.Lock()
		排队 := 0
		if p.队列 != nil {
			排队 = len(p.队列)
		}
		执行中 := p.执行中
		p.锁.Unlock()
		if 排队 == 0 && 执行中 == 0 {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func (p *D线程池) 通知工人退出并等待() {
	p.锁.Lock()
	if p.退出通 != nil {
		选择 := p.退出通
		p.退出通 = nil
		close(选择)
	}
	p.暂停条件.Broadcast()
	p.锁.Unlock()
	p.工人等待.Wait()
}

// D停止 拒绝新投递，排空队列并等待执行中结束。
func (p *D线程池) D停止() {
	p.锁.Lock()
	if !p.已启动 || p.已停止 {
		p.锁.Unlock()
		return
	}
	p.已停止 = true
	p.强制 = false
	p.已暂停 = false
	p.投递取消()
	p.暂停条件.Broadcast()
	p.锁.Unlock()

	p.等待排空()
	p.通知工人退出并等待()
}

// D强制停止 拒绝新投递，丢弃排队任务，等待当前执行中结束。
func (p *D线程池) D强制停止() {
	p.锁.Lock()
	if !p.已启动 || p.已停止 {
		p.锁.Unlock()
		return
	}
	p.已停止 = true
	p.强制 = true
	p.已暂停 = false
	队列 := p.队列
	p.投递取消()
	p.任务取消()
	p.暂停条件.Broadcast()
	p.锁.Unlock()

	// 丢弃排队中的任务
	for {
		选择 := false
		select {
		case 项 := <-队列:
			if 项.句柄 != nil {
				项.句柄.完成结果(nil, 错误任务取消)
			}
			选择 = true
		default:
		}
		if !选择 {
			break
		}
	}

	// 等当前执行中结束
	for {
		p.锁.Lock()
		执行中 := p.执行中
		p.锁.Unlock()
		if 执行中 == 0 {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	p.通知工人退出并等待()
}

func (p *D线程池) 工人循环() {
	defer p.工人等待.Done()
	for {
		p.锁.Lock()
		for p.已暂停 && !p.已停止 {
			p.暂停条件.Wait()
		}
		if p.存活工人数 > p.最大线程数 {
			p.存活工人数--
			p.锁.Unlock()
			return
		}
		强制 := p.强制
		队列 := p.队列
		退出通 := p.退出通
		p.锁.Unlock()

		if 强制 {
			p.锁.Lock()
			p.存活工人数--
			p.锁.Unlock()
			return
		}

		if 退出通 == nil {
			p.锁.Lock()
			p.存活工人数--
			p.锁.Unlock()
			return
		}

		select {
		case <-退出通:
			p.锁.Lock()
			p.存活工人数--
			p.锁.Unlock()
			return
		case 项 := <-队列:
			p.锁.Lock()
			if p.强制 {
				p.锁.Unlock()
				if 项.句柄 != nil {
					项.句柄.完成结果(nil, 错误任务取消)
				}
				continue
			}
			p.执行中++
			p.锁.Unlock()

			失败 := p.执行一项(项)

			p.锁.Lock()
			p.执行中--
			p.锁.Unlock()
			p.已完成数.Add(1)
			if 失败 {
				p.失败数.Add(1)
			}
		case <-time.After(50 * time.Millisecond):
			// 空闲时检查缩容
			continue
		}
	}
}

func (p *D线程池) 执行一项(项 *内部任务) (失败 bool) {
	defer func() {
		if 恢复 := recover(); 恢复 != nil {
			失败 = true
			错误 := fmt.Errorf("任务panic: %v", 恢复)
			if 项.句柄 != nil {
				项.句柄.完成结果(nil, 错误)
			}
		}
	}()

	if 项.有返回 != nil {
		结果, 错误 := 项.有返回()
		if 项.句柄 != nil {
			项.句柄.完成结果(结果, 错误)
		}
		return 错误 != nil
	}
	if 项.无返回 != nil {
		项.无返回()
	}
	return false
}
