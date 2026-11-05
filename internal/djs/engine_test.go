package djs

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/Darren5211314/go_common/internal/dlog"
)

func TestGoja执行与变量(t *testing.T) {
	引擎, err := 新建引擎("goja")
	if err != nil {
		t.Fatal(err)
	}
	defer 引擎.D关闭()

	if err := 引擎.D设置变量("a", 10); err != nil {
		t.Fatal(err)
	}
	if err := 引擎.D设置变量("name", "张三"); err != nil {
		t.Fatal(err)
	}
	结果, err := 引擎.D执行(`a + 2`)
	dlog.Info("执行", "结果", 结果, "错误", err)
	if err != nil || 结果 != "12" {
		t.Fatalf("算术不符: %q %v", 结果, err)
	}
	_, err = 引擎.D执行(`var ok = true; var score = 3.5;`)
	if err != nil {
		t.Fatal(err)
	}
	整数, err := 引擎.D取整数("a")
	if err != nil || 整数 != 10 {
		t.Fatalf("取整数不符: %v %v", 整数, err)
	}
	文本, err := 引擎.D取文本("name")
	if err != nil || 文本 != "张三" {
		t.Fatalf("取文本不符: %v %v", 文本, err)
	}
	逻辑, err := 引擎.D取逻辑("ok")
	if err != nil || !逻辑 {
		t.Fatalf("取逻辑不符: %v %v", 逻辑, err)
	}
	小数, err := 引擎.D取小数("score")
	if err != nil || 小数 != 3.5 {
		t.Fatalf("取小数不符: %v %v", 小数, err)
	}
	_, err = 引擎.D执行(`var empty = null;`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = 引擎.D取文本("empty")
	if err == nil || err.Error() != 错误变量不存在.Error() {
		t.Fatalf("null 变量应视为不存在，得到 %v", err)
	}
}

func TestGoja执行文件(t *testing.T) {
	目录 := t.TempDir()
	路径 := filepath.Join(目录, "demo.js")
	if err := os.WriteFile(路径, []byte("1 + 2 + 3"), 0o644); err != nil {
		t.Fatal(err)
	}
	引擎, err := 新建引擎("")
	if err != nil {
		t.Fatal(err)
	}
	defer 引擎.D关闭()
	结果, err := 引擎.D执行文件(路径)
	if err != nil || 结果 != "6" {
		t.Fatalf("执行文件不符: %q %v", 结果, err)
	}
}

func TestGoja单实例并发不崩(t *testing.T) {
	引擎, err := 新建引擎("goja")
	if err != nil {
		t.Fatal(err)
	}
	defer 引擎.D关闭()
	_ = 引擎.D设置变量("n", int64(0))

	var 等待 sync.WaitGroup
	var 失败 atomic.Int64
	并发 := 32
	每协程 := 50
	等待.Add(并发)
	for i := 0; i < 并发; i++ {
		go func() {
			defer 等待.Done()
			for j := 0; j < 每协程; j++ {
				if _, err := 引擎.D执行(`n = n + 1`); err != nil {
					失败.Add(1)
					return
				}
			}
		}()
	}
	等待.Wait()
	if 失败.Load() != 0 {
		t.Fatalf("并发执行出现错误数 %d", 失败.Load())
	}
	最终, err := 引擎.D取整数("n")
	dlog.Info("并发累加", "最终", 最终, "期望", 并发*每协程)
	if err != nil || 最终 != int64(并发*每协程) {
		t.Fatalf("累加不符: %v %v", 最终, err)
	}
}

func Test多实例并发(t *testing.T) {
	var 等待 sync.WaitGroup
	var 失败 atomic.Int64
	并发 := 16
	等待.Add(并发)
	for i := 0; i < 并发; i++ {
		go func(序号 int) {
			defer 等待.Done()
			引擎, err := 新建引擎("goja")
			if err != nil {
				失败.Add(1)
				return
			}
			defer 引擎.D关闭()
			_ = 引擎.D设置变量("x", 序号)
			结果, err := 引擎.D执行(`x * 2`)
			if err != nil {
				失败.Add(1)
				return
			}
			期望 := strconv.Itoa(序号 * 2)
			if 结果 != 期望 {
				失败.Add(1)
			}
		}(i + 1)
	}
	等待.Wait()
	if 失败.Load() != 0 {
		t.Fatalf("多实例并发失败 %d", 失败.Load())
	}
}

func TestV8可选(t *testing.T) {
	引擎, err := 新建引擎("v8")
	if err != nil {
		t.Log("V8 不可用（无 CGO 或未装库），跳过:", err)
		return
	}
	defer 引擎.D关闭()
	_ = 引擎.D设置变量("a", 7)
	结果, err := 引擎.D执行(`a + 3`)
	dlog.Info("V8执行", "结果", 结果, "错误", err)
	if err != nil || 结果 != "10" {
		t.Fatalf("V8 算术不符: %q %v", 结果, err)
	}
}
