package dhttp

import (
	"context"
	"errors"
	"net"
	"sync/atomic"
	"testing"
	"time"
)

type 慢拨号 struct {
	延迟   time.Duration
	关闭数 *atomic.Int32
}

func (d *慢拨号) Dial(network, address string) (net.Conn, error) {
	time.Sleep(d.延迟)
	客户端, 服务端 := net.Pipe()
	go func() {
		_ = 服务端.Close()
	}()
	return &计数连接{Conn: 客户端, 关闭数: d.关闭数}, nil
}

type 计数连接 struct {
	net.Conn
	关闭数  *atomic.Int32
	已关  atomic.Bool
}

func (c *计数连接) Close() error {
	if c.已关.CompareAndSwap(false, true) {
		c.关闭数.Add(1)
	}
	return c.Conn.Close()
}

func TestSOCKS5取消后关闭连接(t *testing.T) {
	关闭数 := &atomic.Int32{}
	拨号 := 包装SOCKS5拨号(&慢拨号{延迟: 200 * time.Millisecond, 关闭数: 关闭数})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	连接, err := 拨号(ctx, "tcp", "example.com:80")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("期望 Canceled，得到 %v", err)
	}
	if 连接 != nil {
		t.Fatal("取消时应无连接")
	}
	截止 := time.Now().Add(time.Second)
	for time.Now().Before(截止) && 关闭数.Load() == 0 {
		time.Sleep(10 * time.Millisecond)
	}
	if 关闭数.Load() != 1 {
		t.Fatalf("拨通后应 Close，关闭数=%d", 关闭数.Load())
	}
}

func Test取直连传输复用(t *testing.T) {
	a := 取直连传输()
	b := 取直连传输()
	if a != b {
		t.Fatal("直连传输应复用同一实例")
	}
	if a.DisableKeepAlives {
		t.Fatal("直连应启用 KeepAlive")
	}
}
