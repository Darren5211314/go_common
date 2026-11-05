package dsync

import (
	"sync"
	"testing"
	"time"

	"github.com/Darren5211314/go_common/internal/dlog"
)

func Test互斥锁累加(t *testing.T) {
	锁 := 新建互斥锁()
	var 合计 int64
	var 等待 sync.WaitGroup
	并发, 每协程 := 32, 1000
	等待.Add(并发)
	for i := 0; i < 并发; i++ {
		go func() {
			defer 等待.Done()
			for j := 0; j < 每协程; j++ {
				锁.D加锁()
				合计++
				锁.D解锁()
			}
		}()
	}
	等待.Wait()
	期望 := int64(并发 * 每协程)
	dlog.Info("互斥累加", "合计", 合计, "期望", 期望)
	if 合计 != 期望 {
		t.Fatalf("互斥累加不符")
	}
	if 锁.D尝试加锁() {
		锁.D解锁()
	} else {
		t.Fatal("空闲时应能尝试加锁")
	}
}

func Test读写锁(t *testing.T) {
	锁 := 新建读写锁()
	锁.D写加锁()
	if 锁.D尝试读加锁() {
		t.Fatal("持有写锁时不应拿到读锁")
	}
	if 锁.D尝试写加锁() {
		t.Fatal("持有写锁时不应再拿到写锁")
	}
	锁.D写解锁()

	锁.D读加锁()
	锁.D读加锁()
	if !锁.D尝试读加锁() {
		t.Fatal("读锁可重入共享")
	}
	锁.D读解锁()
	锁.D读解锁()
	锁.D读解锁()

	开始 := make(chan struct{})
	完成 := make(chan struct{})
	go func() {
		<-开始
		锁.D写加锁()
		close(完成)
		锁.D写解锁()
	}()
	锁.D读加锁()
	close(开始)
	select {
	case <-完成:
		t.Fatal("读锁未释放前写锁不应成功")
	case <-time.After(50 * time.Millisecond):
	}
	锁.D读解锁()
	select {
	case <-完成:
	case <-time.After(time.Second):
		t.Fatal("释放读锁后写锁应成功")
	}
}

func Test原子计数(t *testing.T) {
	计数 := 新建原子计数()
	var 等待 sync.WaitGroup
	并发, 每协程 := 32, 1000
	等待.Add(并发)
	for i := 0; i < 并发; i++ {
		go func() {
			defer 等待.Done()
			for j := 0; j < 每协程; j++ {
				计数.D加(1)
			}
		}()
	}
	等待.Wait()
	期望 := int64(并发 * 每协程)
	dlog.Info("原子累加", "值", 计数.D取(), "期望", 期望)
	if 计数.D取() != 期望 {
		t.Fatalf("原子累加不符")
	}

	初 := 新建原子计数初值(10)
	if 初.D取() != 10 {
		t.Fatal("初值不符")
	}
	if 初.D减(3) != 7 {
		t.Fatal("减不符")
	}
	初.D设置(100)
	if 初.D取() != 100 {
		t.Fatal("设置不符")
	}
	if !初.D比较交换(100, 200) || 初.D取() != 200 {
		t.Fatal("CAS 成功路径不符")
	}
	if 初.D比较交换(100, 300) {
		t.Fatal("CAS 失败路径应返回 false")
	}
	if 初.D取() != 200 {
		t.Fatal("CAS 失败后值应不变")
	}
}
