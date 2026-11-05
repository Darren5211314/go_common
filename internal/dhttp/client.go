package dhttp

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	错误空地址 = errors.New("URL为空")
)

// 代理池 最小接口，避免强耦合具体类型。
type 代理池 interface {
	D获取代理(等待秒 int) (string, error)
}

// D客户端 HTTP 客户端。
type D客户端 struct {
	锁 sync.Mutex

	代理池实例  代理池
	默认超时秒  int
	默认重试次数 int
}

// D新建 创建客户端。
func 新建() *D客户端 {
	return &D客户端{
		默认超时秒:  10,
		默认重试次数: 3,
	}
}

func (c *D客户端) D设置代理池(池 代理池) {
	c.锁.Lock()
	c.代理池实例 = 池
	c.锁.Unlock()
}

func (c *D客户端) D设置默认超时(秒 int) {
	if 秒 <= 0 {
		秒 = 10
	}
	c.锁.Lock()
	c.默认超时秒 = 秒
	c.锁.Unlock()
}

func (c *D客户端) D设置默认重试(次数 int) {
	if 次数 <= 0 {
		次数 = 3
	}
	c.锁.Lock()
	c.默认重试次数 = 次数
	c.锁.Unlock()
}

// D新建会话 基于当前客户端配置创建会话（Cookie 持久 + 代理粘性）。
func (c *D客户端) D新建会话() *D会话 {
	c.锁.Lock()
	defer c.锁.Unlock()
	return &D会话{
		客户端:  c,
		Cookie: make(map[string]string),
	}
}

// D会话 保持 Cookie 与粘性代理（仅记录池代理）。
type D会话 struct {
	客户端  *D客户端
	锁    sync.Mutex
	Cookie map[string]string
	粘性代理 string
}

func (s *D会话) D发送(方法, 地址 string, 数据 any, Cookie string, 协议头 map[string]string, 禁止重定向, 使用代理 bool, 代理 string, 最大重试 int, 自动合并Cookie bool, 超时秒 int) (*D响应, error) {
	return s.客户端.发送内部(s, 方法, 地址, 数据, Cookie, 协议头, 禁止重定向, 使用代理, 代理, 最大重试, 自动合并Cookie, 超时秒)
}

// D发送 单次请求（无会话粘性）。
// 使用代理=false 强制直连；=true 且 代理非空 用手动代理；=true 且 代理空 则粘性/池/直连。
func (c *D客户端) D发送(方法, 地址 string, 数据 any, Cookie string, 协议头 map[string]string, 禁止重定向, 使用代理 bool, 代理 string, 最大重试 int, 自动合并Cookie bool, 超时秒 int) (*D响应, error) {
	return c.发送内部(nil, 方法, 地址, 数据, Cookie, 协议头, 禁止重定向, 使用代理, 代理, 最大重试, 自动合并Cookie, 超时秒)
}

type 发送参数 struct {
	方法         string
	地址         string
	数据         any
	Cookie入参   string
	协议头        map[string]string
	禁止重定向      bool
	使用代理       bool
	手动代理       string
	最大尝试       int
	自动合并Cookie bool
	超时         time.Duration
}

func (c *D客户端) 发送内部(会话 *D会话, 方法, 地址 string, 数据 any, Cookie string, 协议头 map[string]string, 禁止重定向, 使用代理 bool, 代理 string, 最大重试 int, 自动合并Cookie bool, 超时秒 int) (*D响应, error) {
	地址 = strings.TrimSpace(地址)
	if 地址 == "" {
		return nil, 错误空地址
	}
	方法 = strings.ToUpper(strings.TrimSpace(方法))
	if 方法 == "" {
		方法 = http.MethodGet
	}
	手动代理 := strings.TrimSpace(代理)

	c.锁.Lock()
	默认超时 := c.默认超时秒
	默认重试 := c.默认重试次数
	池 := c.代理池实例
	c.锁.Unlock()

	if 超时秒 <= 0 {
		超时秒 = 默认超时
	}
	if 最大重试 <= 0 {
		最大重试 = 默认重试
	}

	参数 := 发送参数{
		方法:         方法,
		地址:         地址,
		数据:         数据,
		Cookie入参:   Cookie,
		协议头:        协议头,
		禁止重定向:      禁止重定向,
		使用代理:       使用代理,
		手动代理:       手动代理,
		最大尝试:       最大重试,
		自动合并Cookie: 自动合并Cookie,
		超时:         time.Duration(超时秒) * time.Second,
	}

	var 最后错误 error
	var 最后响应 *D响应
	当前代理 := ""

	for 尝试 := 1; 尝试 <= 参数.最大尝试; 尝试++ {
		选中, 来自手动, 选错 := c.选择代理(会话, 参数.使用代理, 参数.手动代理, 池, 当前代理, 超时秒)
		if 选错 != nil {
			最后错误 = 选错
			continue
		}
		当前代理 = 选中

		响应, 错误 := c.执行一次(会话, 参数, 选中)
		最后响应 = 响应
		最后错误 = 错误

		if 错误 == nil {
			// 仅池/粘性成功才写会话粘性；手动不污染粘性
			if 会话 != nil && 选中 != "" && !来自手动 {
				会话.锁.Lock()
				会话.粘性代理 = 选中
				会话.锁.Unlock()
			}
			return 响应, nil
		}

		if !可重试(错误, 响应) || 尝试 >= 参数.最大尝试 {
			break
		}
		if 来自手动 {
			// 手动失败重试：仍用同一手动
			当前代理 = 参数.手动代理
			continue
		}
		// 池/粘性失败：清空粘性，换新 IP
		if 会话 != nil {
			会话.锁.Lock()
			if 会话.粘性代理 == 当前代理 {
				会话.粘性代理 = ""
			}
			会话.锁.Unlock()
		}
		当前代理 = ""
	}
	if 最后响应 != nil {
		return 最后响应, 最后错误
	}
	return nil, 最后错误
}

// 选择代理 返回 (代理URL, 是否来自手动, error)。
func (c *D客户端) 选择代理(会话 *D会话, 使用代理 bool, 手动代理 string, 池 代理池, 当前代理 string, 超时秒 int) (string, bool, error) {
	if !使用代理 {
		return "", false, nil
	}
	if 手动代理 != "" {
		return 手动代理, true, nil
	}
	if 当前代理 != "" {
		return 当前代理, false, nil
	}
	if 会话 != nil {
		会话.锁.Lock()
		粘性 := 会话.粘性代理
		会话.锁.Unlock()
		if 粘性 != "" {
			return 粘性, false, nil
		}
	}
	if 池 == nil {
		return "", false, nil // 无池则直连
	}
	等待 := 超时秒
	if 等待 <= 0 {
		等待 = 5
	}
	取到, 错误 := 池.D获取代理(等待)
	return 取到, false, 错误
}

func (c *D客户端) 执行一次(会话 *D会话, 参数 发送参数, 代理 string) (*D响应, error) {
	开始 := time.Now()

	正文, 内容类型, 错误 := 编码数据(参数.数据)
	if 错误 != nil {
		return nil, 错误
	}

	请求Cookie映射 := 解析Cookie串(参数.Cookie入参)
	会话Cookie := map[string]string{}
	if 会话 != nil {
		会话.锁.Lock()
		会话Cookie = 合并Cookie映射(会话.Cookie)
		会话.锁.Unlock()
	}
	合并后 := 请求Cookie映射
	if 参数.自动合并Cookie {
		合并后 = 合并Cookie映射(会话Cookie, 请求Cookie映射)
	}
	Cookie串 := Cookie映射转串(合并后)

	var 传输 http.RoundTripper
	if strings.TrimSpace(代理) == "" {
		传输 = 取直连传输()
	} else {
		代理传输, 错误 := 新建代理传输(代理)
		if 错误 != nil {
			return nil, 错误
		}
		传输 = 代理传输
	}
	客户端 := &http.Client{
		Timeout:   参数.超时,
		Transport: 传输,
	}
	if 参数.禁止重定向 {
		客户端.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	var 体阅读 io.Reader
	if len(正文) > 0 {
		体阅读 = bytes.NewReader(正文)
	}
	请求, 错误 := http.NewRequestWithContext(context.Background(), 参数.方法, 参数.地址, 体阅读)
	if 错误 != nil {
		return nil, 错误
	}
	for 键, 值 := range 参数.协议头 {
		请求.Header.Set(键, 值)
	}
	if 内容类型 != "" && 请求.Header.Get("Content-Type") == "" {
		请求.Header.Set("Content-Type", 内容类型)
	}
	if Cookie串 != "" {
		请求.Header.Set("Cookie", Cookie串)
	}

	响应对象, 错误 := 客户端.Do(请求)
	耗时 := time.Since(开始).Milliseconds()
	if 耗时 <= 0 {
		耗时 = 1 // 亚毫秒本机请求仍记 1ms，避免 D取耗时 为 0
	}
	if 错误 != nil {
		return &D响应{
			使用代理:   代理,
			耗时毫秒:   耗时,
			全部Cookie: Cookie串,
		}, 错误
	}
	defer 响应对象.Body.Close()
	体, 读错 := io.ReadAll(io.LimitReader(响应对象.Body, 32<<20))
	if 读错 != nil {
		错误 = 读错
	}

	返回Cookie映射 := 从响应提取Cookie(响应对象.Header)
	返回Cookie串 := Cookie映射转串(返回Cookie映射)
	全部 := 合并后
	if 参数.自动合并Cookie {
		全部 = 合并Cookie映射(合并后, 返回Cookie映射)
		if 会话 != nil {
			会话.锁.Lock()
			会话.Cookie = 合并Cookie映射(会话.Cookie, 返回Cookie映射)
			会话.锁.Unlock()
		}
	}

	重定向 := ""
	if 参数.禁止重定向 {
		重定向 = 解析重定向(响应对象)
	} else if 响应对象.Request != nil && 响应对象.Request.URL != nil {
		重定向 = 响应对象.Request.URL.String()
	}

	结果 := &D响应{
		状态码:     响应对象.StatusCode,
		正文:       体,
		重定向地址:   重定向,
		返回Cookie: 返回Cookie串,
		全部Cookie: Cookie映射转串(全部),
		协议头:     响应对象.Header.Clone(),
		使用代理:    代理,
		耗时毫秒:    耗时,
	}
	if 错误 != nil {
		return 结果, 错误
	}
	if 可重试状态码(响应对象.StatusCode) {
		return 结果, errors.New("HTTP状态码可重试: " + 响应对象.Status)
	}
	return 结果, nil
}

func 编码数据(数据 any) ([]byte, string, error) {
	if 数据 == nil {
		return nil, "", nil
	}
	switch v := 数据.(type) {
	case []byte:
		return v, "", nil
	case string:
		if v == "" {
			return nil, "", nil
		}
		return []byte(v), "", nil
	case map[string]string:
		值 := url.Values{}
		for 键, 一项 := range v {
			值.Set(键, 一项)
		}
		return []byte(值.Encode()), "application/x-www-form-urlencoded", nil
	default:
		return nil, "", errors.New("不支持的data类型")
	}
}

func 解析重定向(响应 *http.Response) string {
	位置 := 响应.Header.Get("Location")
	if 位置 == "" {
		return ""
	}
	参考, 错误 := url.Parse(位置)
	if 错误 != nil {
		return 位置
	}
	if 参考.IsAbs() {
		return 参考.String()
	}
	if 响应.Request != nil && 响应.Request.URL != nil {
		return 响应.Request.URL.ResolveReference(参考).String()
	}
	return 位置
}

func 可重试状态码(码 int) bool {
	return 码 == 502 || 码 == 503 || 码 == 504
}

func 可重试(错误 error, 响应 *D响应) bool {
	if 错误 == nil {
		return false
	}
	if 响应 != nil && 可重试状态码(响应.状态码) {
		return true
	}
	消息 := strings.ToLower(错误.Error())
	关键字 := []string{
		"timeout", "timed out", "deadline exceeded",
		"connection refused", "connection reset", "i/o timeout",
		"proxy", "eof", "broken pipe", "tls", "handshake",
		"temporary", "unavailable", "可重试",
	}
	for _, 键 := range 关键字 {
		if strings.Contains(消息, 键) {
			return true
		}
	}
	var 网络错误 net.Error
	if errors.As(错误, &网络错误) {
		return 网络错误.Timeout() || 网络错误.Temporary()
	}
	return false
}
