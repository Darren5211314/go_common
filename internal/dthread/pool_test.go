package dthread

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Darren5211314/go_common/internal/dlog"
)

func Test投递与统计(t *testing.T) {
	池 := 新建(4)
	池.D设置队列容量(16)
	if err := 池.D启动(); err != nil {
		t.Fatal(err)
	}
	defer 池.D停止()

	var 计数 atomic.Int64
	var 等待 sync.WaitGroup
	等待.Add(20)
	for i := 0; i < 20; i++ {
		if err := 池.D投递(func() {
			计数.Add(1)
			等待.Done()
		}); err != nil {
			t.Fatalf("投递失败: %v", err)
		}
	}
	等待.Wait()
	统计 := 池.D取统计()
	dlog.Info("统计", "完成", 统计.D已完成数, "投递", 统计.D已投递数, "计数", 计数.Load())
	if 计数.Load() != 20 || 统计.D已投递数 != 20 {
		t.Errorf("统计不符")
	}
}

func Test有返回任务(t *testing.T) {
	池 := 新建(2)
	_ = 池.D启动()
	defer 池.D停止()

	任务, err := 池.D投递任务(func() (any, error) {
		return "ok", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	结果, err := 任务.D等待()
	dlog.Info("任务结果", "结果", 结果, "错误", err)
	if err != nil || 结果 != "ok" {
		t.Errorf("等待结果不符")
	}
}

func Test暂停继续(t *testing.T) {
	池 := 新建(2)
	_ = 池.D启动()
	defer 池.D停止()

	池.D暂停()
	var 执行 atomic.Int64
	for i := 0; i < 4; i++ {
		_ = 池.D投递(func() {
			执行.Add(1)
			time.Sleep(50 * time.Millisecond)
		})
	}
	time.Sleep(100 * time.Millisecond)
	if 执行.Load() != 0 {
		t.Errorf("暂停期间不应执行，得到 %d", 执行.Load())
	}
	池.D继续()
	截止 := time.Now().Add(2 * time.Second)
	for time.Now().Before(截止) && 执行.Load() < 4 {
		time.Sleep(20 * time.Millisecond)
	}
	if 执行.Load() != 4 {
		t.Errorf("继续后应执行完，得到 %d", 执行.Load())
	}
}

func Test动态扩容(t *testing.T) {
	池 := 新建(1)
	_ = 池.D启动()
	defer 池.D停止()

	阻塞 := make(chan struct{})
	_ = 池.D投递(func() { <-阻塞 })
	截止 := time.Now().Add(time.Second)
	for time.Now().Before(截止) && 池.D取执行中数量() < 1 {
		time.Sleep(5 * time.Millisecond)
	}
	池.D设置最大线程数(3)
	var 额外 atomic.Int64
	_ = 池.D投递(func() { 额外.Add(1) })
	_ = 池.D投递(func() { 额外.Add(1) })
	截止 = time.Now().Add(time.Second)
	for time.Now().Before(截止) && 额外.Load() < 2 {
		time.Sleep(5 * time.Millisecond)
	}
	close(阻塞)
	if 额外.Load() < 2 {
		t.Errorf("扩容后应能并行执行额外任务")
	}
	dlog.Info("扩容测试", "额外", 额外.Load(), "最大", 池.D取统计().D最大线程数)
}

func Test缩容空闲统计(t *testing.T) {
	池 := 新建(3)
	_ = 池.D启动()
	defer 池.D停止()

	截止 := time.Now().Add(time.Second)
	for time.Now().Before(截止) && 池.D取统计().D最大线程数 < 3 {
		time.Sleep(5 * time.Millisecond)
	}
	阻塞 := make(chan struct{})
	_ = 池.D投递(func() { <-阻塞 })
	截止 = time.Now().Add(time.Second)
	for time.Now().Before(截止) && 池.D取执行中数量() < 1 {
		time.Sleep(5 * time.Millisecond)
	}
	池.D设置最大线程数(1)
	空闲 := 池.D取空闲数量()
	统计 := 池.D取统计()
	dlog.Info("缩容空闲", "空闲", 空闲, "统计空闲", 统计.D空闲数量, "执行中", 统计.D执行中数量)
	if 空闲 != 统计.D空闲数量 {
		t.Fatalf("D取空闲数量 与统计不一致")
	}
	// 缩容过渡：存活工人可能仍 > 最大；空闲应为 存活-执行中，不应按最大失真为 0
	if 空闲 < 1 {
		t.Fatalf("缩容过渡期空闲应反映多余存活工人，得到 %d", 空闲)
	}
	close(阻塞)
}

func TestPanic计失败(t *testing.T) {
	池 := 新建(1)
	_ = 池.D启动()
	defer 池.D停止()

	_ = 池.D投递(func() { panic("x") })
	截止 := time.Now().Add(time.Second)
	for time.Now().Before(截止) && 池.D取统计().D失败数 < 1 {
		time.Sleep(5 * time.Millisecond)
	}
	if 池.D取统计().D失败数 < 1 {
		t.Errorf("panic 应计入失败")
	}
}

func Test停止后拒绝投递(t *testing.T) {
	池 := 新建(1)
	_ = 池.D启动()
	池.D停止()
	if err := 池.D投递(func() {}); err == nil {
		t.Errorf("停止后应拒绝投递")
	}
}
