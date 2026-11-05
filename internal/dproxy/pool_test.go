package dproxy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Darren5211314/go_common/internal/dlog"
)

func Test池无验证获取与统计(t *testing.T) {
	var 拉取次数 atomic.Int64
	源 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := 拉取次数.Add(1)
		fmt.Fprintf(w, "10.0.0.%d:8000\n10.0.0.%d:8001\n", n, n)
	}))
	defer 源.Close()

	代理池 := 新建()
	代理池.D设置代理链接(源.URL)
	代理池.D设置代理池容量(2)
	代理池.D设置是否验证(false)
	代理池.D设置获取间隔(1)
	代理池.D设置存活时间(0)
	if err := 代理池.D启动(); err != nil {
		t.Fatalf("启动失败: %v", err)
	}
	defer 代理池.D关闭()

	代理1, err := 代理池.D获取代理(5)
	dlog.Info("获取1", "代理", 代理1, "错误", err)
	if err != nil || 代理1 == "" {
		t.Fatalf("获取1失败: %v", err)
	}
	代理2, err := 代理池.D获取代理(5)
	dlog.Info("获取2", "代理", 代理2, "错误", err)
	if err != nil || 代理2 == "" {
		t.Fatalf("获取2失败: %v", err)
	}
	if 代理1 == 代理2 {
		t.Errorf("两次取出不应相同")
	}

	统计 := 代理池.D取统计()
	dlog.Info("统计", "从链接", 统计.D从链接获取数, "已取出", 统计.D已取出数, "可用", 统计.D池中可用数)
	if 统计.D已取出数 < 2 || 统计.D从链接获取数 < 2 {
		t.Errorf("统计不符: %+v", 统计)
	}
}

func Test池过期剔除(t *testing.T) {
	源 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "10.1.1.1:9000\n")
	}))
	defer 源.Close()

	代理池 := 新建()
	代理池.D设置代理链接(源.URL)
	代理池.D设置代理池容量(1)
	代理池.D设置是否验证(false)
	代理池.D设置存活时间(1)
	代理池.D设置获取间隔(1)
	if err := 代理池.D启动(); err != nil {
		t.Fatalf("启动失败: %v", err)
	}
	defer 代理池.D关闭()

	时间起点 := time.Now()
	for time.Since(时间起点) < 3*time.Second {
		统计 := 代理池.D取统计()
		if 统计.D池中可用数 >= 1 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	time.Sleep(1100 * time.Millisecond)
	统计 := 代理池.D取统计()
	dlog.Info("过期后统计", "可用", 统计.D池中可用数, "失效", 统计.D失效数)
	if 统计.D失效数 < 1 {
		t.Errorf("期望有失效计数")
	}
}

func Test池获取超时(t *testing.T) {
	源 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		fmt.Fprint(w, "10.2.2.2:9000\n")
	}))
	defer 源.Close()

	代理池 := 新建()
	代理池.D设置代理链接(源.URL)
	代理池.D设置代理池容量(1)
	代理池.D设置是否验证(false)
	代理池.D设置获取间隔(1)
	if err := 代理池.D启动(); err != nil {
		t.Fatalf("启动失败: %v", err)
	}
	defer 代理池.D关闭()

	_, err := 代理池.D获取代理(0)
	dlog.Info("超时获取", "错误", err)
	if err != 错误获取超时 {
		t.Errorf("期望超时错误，得到 %v", err)
	}
}
