package dhttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Darren5211314/go_common/internal/dlog"
)

func Test发送GET与Cookie(t *testing.T) {
	源 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dlog.Info("收到请求", "cookie", r.Header.Get("Cookie"), "方法", r.Method)
		http.SetCookie(w, &http.Cookie{Name: "sid", Value: "xyz"})
		w.Header().Set("X-Test", "1")
		fmt.Fprint(w, "hello")
	}))
	defer 源.Close()

	客户端 := 新建()
	响应, err := 客户端.D发送("GET", 源.URL, nil, "a=1", map[string]string{"User-Agent": "github.com/Darren5211314/go_common"}, false, false, "", 1, true, 5)
	if err != nil {
		t.Fatalf("发送失败: %v", err)
	}
	dlog.Info("响应", "正文", 响应.D取正文(), "耗时", 响应.D取耗时(), "返回Cookie", 响应.D取返回Cookie(), "全部", 响应.D取Cookie())
	if 响应.D取正文() != "hello" || 响应.D取状态码() != 200 {
		t.Errorf("正文/状态不符")
	}
	if 响应.D取耗时() <= 0 {
		t.Errorf("耗时异常")
	}
	if !strings.Contains(响应.D取返回Cookie(), "sid=xyz") {
		t.Errorf("返回Cookie不符")
	}
	if !strings.Contains(响应.D取Cookie(), "a=1") || !strings.Contains(响应.D取Cookie(), "sid=xyz") {
		t.Errorf("合并Cookie不符: %s", 响应.D取Cookie())
	}
	if 响应.D取协议头值("X-Test") != "1" {
		t.Errorf("协议头值不符")
	}
}

func Test禁止重定向(t *testing.T) {
	源 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/from" {
			http.Redirect(w, r, "/to", http.StatusFound)
			return
		}
		fmt.Fprint(w, "target")
	}))
	defer 源.Close()

	客户端 := 新建()
	响应, err := 客户端.D发送("GET", 源.URL+"/from", nil, "", nil, true, false, "", 1, true, 5)
	if err != nil {
		t.Fatalf("发送失败: %v", err)
	}
	地址 := 响应.D取重定向地址()
	dlog.Info("重定向", "地址", 地址, "状态", 响应.D取状态码())
	if 响应.D取状态码() != http.StatusFound {
		t.Errorf("应返回302")
	}
	if !strings.Contains(地址, "/to") {
		t.Errorf("重定向地址不符: %s", 地址)
	}
}

func TestPOST表单(t *testing.T) {
	源 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		fmt.Fprint(w, r.Form.Get("name"))
	}))
	defer 源.Close()

	客户端 := 新建()
	响应, err := 客户端.D发送("POST", 源.URL, map[string]string{"name": "张三"}, "", nil, false, false, "", 1, true, 5)
	if err != nil {
		t.Fatalf("发送失败: %v", err)
	}
	if 响应.D取正文() != "张三" {
		t.Errorf("表单不符: %s", 响应.D取正文())
	}
}

func Test会话Cookie持久(t *testing.T) {
	源 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/set" {
			http.SetCookie(w, &http.Cookie{Name: "k", Value: "v"})
			fmt.Fprint(w, "set")
			return
		}
		fmt.Fprint(w, r.Header.Get("Cookie"))
	}))
	defer 源.Close()

	客户端 := 新建()
	会话 := 客户端.D新建会话()
	_, err := 会话.D发送("GET", 源.URL+"/set", nil, "", nil, false, false, "", 1, true, 5)
	if err != nil {
		t.Fatal(err)
	}
	响应, err := 会话.D发送("GET", 源.URL+"/echo", nil, "", nil, false, false, "", 1, true, 5)
	if err != nil {
		t.Fatal(err)
	}
	dlog.Info("会话cookie", "正文", 响应.D取正文())
	if !strings.Contains(响应.D取正文(), "k=v") {
		t.Errorf("会话未保持cookie")
	}
}

type 假代理池 struct {
	次数 atomic.Int64
	列表 []string
}

func (p *假代理池) D获取代理(等待秒 int) (string, error) {
	i := int(p.次数.Add(1)-1) % len(p.列表)
	return p.列表[i], nil
}

func Test使用代理false不取池(t *testing.T) {
	源 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	}))
	defer 源.Close()

	客户端 := 新建()
	池 := &假代理池{列表: []string{"http://127.0.0.1:1"}}
	客户端.D设置代理池(池)
	响应, err := 客户端.D发送("GET", 源.URL, nil, "", nil, false, false, "", 1, true, 5)
	if err != nil {
		t.Fatal(err)
	}
	if 响应.D取使用代理() != "" {
		t.Errorf("直连不应有代理")
	}
	if 池.次数.Load() != 0 {
		t.Errorf("未使用代理时不应取池")
	}
}

func Test手动代理优先于池(t *testing.T) {
	客户端 := 新建()
	池 := &假代理池{列表: []string{"http://127.0.0.1:9"}}
	客户端.D设置代理池(池)
	手动 := "http://127.0.0.1:8"
	响应, err := 客户端.D发送("GET", "http://127.0.0.1:1/", nil, "", nil, false, true, 手动, 1, true, 1)
	if 响应 == nil {
		t.Fatalf("应有响应对象: %v", err)
	}
	if 响应.D取使用代理() != 手动 {
		t.Errorf("应使用手动代理, 得到 %q", 响应.D取使用代理())
	}
	if 池.次数.Load() != 0 {
		t.Errorf("有手动代理时不应取池")
	}
}

func Test会话粘性与手动不污染(t *testing.T) {
	客户端 := 新建()
	池 := &假代理池{列表: []string{"http://127.0.0.1:7"}}
	客户端.D设置代理池(池)
	会话 := 客户端.D新建会话()

	粘性 := "http://127.0.0.1:6"
	会话.锁.Lock()
	会话.粘性代理 = 粘性
	会话.锁.Unlock()

	响应1, _ := 会话.D发送("GET", "http://127.0.0.1:1/", nil, "", nil, false, true, "", 1, true, 1)
	if 响应1 == nil || 响应1.D取使用代理() != 粘性 {
		t.Fatalf("应复用粘性代理, 得到 %#v", 响应1)
	}
	if 池.次数.Load() != 0 {
		t.Errorf("有粘性时不应取池")
	}

	手动 := "http://127.0.0.1:5"
	响应2, _ := 会话.D发送("GET", "http://127.0.0.1:1/", nil, "", nil, false, true, 手动, 1, true, 1)
	if 响应2 == nil || 响应2.D取使用代理() != 手动 {
		t.Fatalf("应使用手动代理")
	}
	会话.锁.Lock()
	粘性仍在 := 会话.粘性代理
	会话.锁.Unlock()
	if 粘性仍在 != 粘性 {
		t.Errorf("手动成功/失败都不应改写粘性, 得到 %q", 粘性仍在)
	}
	if 池.次数.Load() != 0 {
		t.Errorf("手动路径不应取池")
	}

	响应3, _ := 会话.D发送("GET", "http://127.0.0.1:1/", nil, "", nil, false, true, "", 1, true, 1)
	if 响应3 == nil || 响应3.D取使用代理() != 粘性 {
		t.Fatalf("手动后空串仍应走原粘性, 得到 %#v", 响应3)
	}
}

func Test重试次数(t *testing.T) {
	var 次数 atomic.Int64
	源 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := 次数.Add(1)
		if n < 3 {
			w.WriteHeader(503)
			fmt.Fprint(w, "busy")
			return
		}
		fmt.Fprint(w, "ok")
	}))
	defer 源.Close()

	客户端 := 新建()
	开始 := time.Now()
	响应, err := 客户端.D发送("GET", 源.URL, nil, "", nil, false, false, "", 3, true, 5)
	dlog.Info("重试", "错误", err, "正文", 响应.D取正文(), "次数", 次数.Load(), "耗时", time.Since(开始))
	if err != nil {
		t.Fatalf("应最终成功: %v", err)
	}
	if 响应.D取正文() != "ok" || 次数.Load() != 3 {
		t.Errorf("重试行为不符")
	}
}
