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
	errNotStarted   = errors.New("线程池未启动")
	errStopped   = errors.New("线程池已停止")
	errNilTask   = errors.New("任务为空")
	errWaitTimeout = errors.New("等待任务超时")
	errTaskCanceled = errors.New("任务已取消")
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
	done   chan struct{}
	result   any
	err   error
	doneOnce sync.Once
}

func (t *D任务) finish(result any, err error) {
	t.doneOnce.Do(func() {
		t.result = result
		t.err = err
		close(t.done)
	})
}

// D等待 阻塞直到任务完成。
func (t *D任务) D等待() (any, error) {
	if t == nil {
		return nil, errNilTask
	}
	<-t.done
	return t.result, t.err
}

// D等待超时 最多等待秒数；超时返回错误。
func (t *D任务) D等待超时(sec int) (any, error) {
	if t == nil {
		return nil, errNilTask
	}
	if sec < 0 {
		return t.D等待()
	}
	choice := time.NewTimer(time.Duration(sec) * time.Second)
	defer choice.Stop()
	select {
	case <-t.done:
		return t.result, t.err
	case <-choice.C:
		return nil, errWaitTimeout
	}
}

type internalTask struct {
	noResult func()
	withResult func() (any, error)
	handle   *D任务
}

// D线程池 可伸缩 worker 池。
type D线程池 struct {
	mu sync.Mutex

	maxWorkers int
	liveWorkers int
	running   int
	queueCap int

	queue   chan *internalTask
	quitCh chan struct{}

	started bool
	stopped bool
	paused bool
	force  bool

	pauseCond *sync.Cond

	submitCtx context.Context
	submitCancel  context.CancelFunc
	taskCtx context.Context
	taskCanceled  context.CancelFunc

	workerWait sync.WaitGroup

	Submitted atomic.Int64
	Completed atomic.Int64
	Failed   atomic.Int64
}

// D新建 创建线程池；最大线程数 <=0 时为 1。默认队列容量 = 最大*4。
func newClient(maxWorkers int) *D线程池 {
	if maxWorkers <= 0 {
		maxWorkers = 1
	}
	submitCtx, submitCancel := context.WithCancel(context.Background())
	taskCtx, taskCanceled := context.WithCancel(context.Background())
	p := &D线程池{
		maxWorkers: maxWorkers,
		queueCap: maxWorkers * 4,
		submitCtx: submitCtx,
		submitCancel:  submitCancel,
		taskCtx: taskCtx,
		taskCanceled:  taskCanceled,
	}
	p.pauseCond = sync.NewCond(&p.mu)
	return p
}

// D设置队列容量 仅启动前有效。
func (p *D线程池) D设置队列容量(capacity int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.started {
		return
	}
	if capacity < 1 {
		capacity = 1
	}
	p.queueCap = capacity
}

// D任务上下文 返回任务 context（仅强制停止时 cancel）。
func (p *D线程池) D任务上下文() context.Context {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.taskCtx
}

// D启动 启动 worker。
func (p *D线程池) D启动() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stopped {
		return errStopped
	}
	if p.started {
		return nil
	}
	p.queue = make(chan *internalTask, p.queueCap)
	p.quitCh = make(chan struct{})
	p.started = true
	for i := 0; i < p.maxWorkers; i++ {
		p.startWorkersLocked()
	}
	return nil
}

func (p *D线程池) startWorkersLocked() {
	p.liveWorkers++
	p.workerWait.Add(1)
	go p.workerLoop()
}

// D设置最大线程数 动态扩缩。
func (p *D线程池) D设置最大线程数(count int) {
	if count <= 0 {
		count = 1
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.maxWorkers = count
	if !p.started || p.stopped {
		return
	}
	for p.liveWorkers < p.maxWorkers {
		p.startWorkersLocked()
	}
	p.pauseCond.Broadcast()
}

func (p *D线程池) D暂停() {
	p.mu.Lock()
	p.paused = true
	p.mu.Unlock()
}

func (p *D线程池) D继续() {
	p.mu.Lock()
	p.paused = false
	p.pauseCond.Broadcast()
	p.mu.Unlock()
}

func (p *D线程池) canSubmitLocked() error {
	if !p.started {
		return errNotStarted
	}
	if p.stopped {
		return errStopped
	}
	return nil
}

func (p *D线程池) submitInternal(item *internalTask) error {
	p.mu.Lock()
	if err := p.canSubmitLocked(); err != nil {
		p.mu.Unlock()
		return err
	}
	queue := p.queue
	submitCtx := p.submitCtx
	p.mu.Unlock()

	select {
	case queue <- item:
		p.Submitted.Add(1)
		return nil
	case <-submitCtx.Done():
		return errStopped
	}
}

// D投递 投递无返回值任务；队列满时阻塞。
func (p *D线程池) D投递(task func()) error {
	if task == nil {
		return errNilTask
	}
	return p.submitInternal(&internalTask{noResult: task})
}

// D投递任务 投递有返回值任务，返回可等待句柄。
func (p *D线程池) D投递任务(task func() (any, error)) (*D任务, error) {
	if task == nil {
		return nil, errNilTask
	}
	handle := &D任务{done: make(chan struct{})}
	if err := p.submitInternal(&internalTask{withResult: task, handle: handle}); err != nil {
		handle.finish(nil, err)
		return handle, err
	}
	return handle, nil
}

func (p *D线程池) D取执行中数量() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.running
}

func (p *D线程池) D取空闲数量() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	idle := p.liveWorkers - p.running
	if idle < 0 {
		return 0
	}
	return idle
}

func (p *D线程池) D取排队数量() int {
	p.mu.Lock()
	queue := p.queue
	p.mu.Unlock()
	if queue == nil {
		return 0
	}
	return len(queue)
}

func (p *D线程池) D取统计() D线程池统计 {
	p.mu.Lock()
	max := p.maxWorkers
	running := p.running
	alive := p.liveWorkers
	queue := p.queue
	p.mu.Unlock()
	queued := 0
	if queue != nil {
		queued = len(queue)
	}
	idle := alive - running
	if idle < 0 {
		idle = 0
	}
	return D线程池统计{
		D最大线程数: max,
		D执行中数量: running,
		D空闲数量:   idle,
		D排队数量:   queued,
		D已完成数:   p.Completed.Load(),
		D失败数:     p.Failed.Load(),
		D已投递数:   p.Submitted.Load(),
	}
}

func (p *D线程池) waitDrain() {
	for {
		p.mu.Lock()
		queued := 0
		if p.queue != nil {
			queued = len(p.queue)
		}
		running := p.running
		p.mu.Unlock()
		if queued == 0 && running == 0 {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func (p *D线程池) signalWorkersAndWait() {
	p.mu.Lock()
	if p.quitCh != nil {
		choice := p.quitCh
		p.quitCh = nil
		close(choice)
	}
	p.pauseCond.Broadcast()
	p.mu.Unlock()
	p.workerWait.Wait()
}

// D停止 拒绝新投递，排空队列并等待执行中结束。
func (p *D线程池) D停止() {
	p.mu.Lock()
	if !p.started || p.stopped {
		p.mu.Unlock()
		return
	}
	p.stopped = true
	p.force = false
	p.paused = false
	p.submitCancel()
	p.pauseCond.Broadcast()
	p.mu.Unlock()

	p.waitDrain()
	p.signalWorkersAndWait()
}

// D强制停止 拒绝新投递，丢弃排队任务，等待当前执行中结束。
func (p *D线程池) D强制停止() {
	p.mu.Lock()
	if !p.started || p.stopped {
		p.mu.Unlock()
		return
	}
	p.stopped = true
	p.force = true
	p.paused = false
	queue := p.queue
	p.submitCancel()
	p.taskCanceled()
	p.pauseCond.Broadcast()
	p.mu.Unlock()

	// 丢弃排队中的任务
	for {
		choice := false
		select {
		case item := <-queue:
			if item.handle != nil {
				item.handle.finish(nil, errTaskCanceled)
			}
			choice = true
		default:
		}
		if !choice {
			break
		}
	}

	// 等当前执行中结束
	for {
		p.mu.Lock()
		running := p.running
		p.mu.Unlock()
		if running == 0 {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	p.signalWorkersAndWait()
}

func (p *D线程池) workerLoop() {
	defer p.workerWait.Done()
	for {
		p.mu.Lock()
		for p.paused && !p.stopped {
			p.pauseCond.Wait()
		}
		if p.liveWorkers > p.maxWorkers {
			p.liveWorkers--
			p.mu.Unlock()
			return
		}
		force := p.force
		queue := p.queue
		quitCh := p.quitCh
		p.mu.Unlock()

		if force {
			p.mu.Lock()
			p.liveWorkers--
			p.mu.Unlock()
			return
		}

		if quitCh == nil {
			p.mu.Lock()
			p.liveWorkers--
			p.mu.Unlock()
			return
		}

		select {
		case <-quitCh:
			p.mu.Lock()
			p.liveWorkers--
			p.mu.Unlock()
			return
		case item := <-queue:
			p.mu.Lock()
			if p.force {
				p.mu.Unlock()
				if item.handle != nil {
					item.handle.finish(nil, errTaskCanceled)
				}
				continue
			}
			p.running++
			p.mu.Unlock()

			fail := p.runOne(item)

			p.mu.Lock()
			p.running--
			p.mu.Unlock()
			p.Completed.Add(1)
			if fail {
				p.Failed.Add(1)
			}
		case <-time.After(50 * time.Millisecond):
			// 空闲时检查缩容
			continue
		}
	}
}

func (p *D线程池) runOne(item *internalTask) (fail bool) {
	defer func() {
		if resume := recover(); resume != nil {
			fail = true
			err := fmt.Errorf("任务panic: %v", resume)
			if item.handle != nil {
				item.handle.finish(nil, err)
			}
		}
	}()

	if item.withResult != nil {
		result, err := item.withResult()
		if item.handle != nil {
			item.handle.finish(result, err)
		}
		return err != nil
	}
	if item.noResult != nil {
		item.noResult()
	}
	return false
}
