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
	errEmptyURL = errors.New("URL为空")

	defaultClientOnce sync.Once
	defaultClientInst *Client
)

// 代理池 最小接口，避免强耦合具体类型。
type proxyPool interface {
	GetProxy(waitTime int) (string, error)
}

// SendOptions 可选请求参数；未填字段用零值默认。
type SendOptions struct {
	Data             any               // 发送体；nil/空串=无 body；[]byte/string 原样；map[string]string=表单。默认 nil
	Cookie           string            // 额外 Cookie；空=不加。支持 a=1; b=2 与 a=1&b=2。默认 ""
	Headers          map[string]string // 协议头；nil=无额外头。默认 nil
	NoRedirect       bool              // true=禁止跟随重定向；false=跟随。默认 false
	UseProxy         bool              // false=强制直连；true=走代理（见 Proxy）。默认 false
	Proxy            string            // 手动代理 URL；UseProxy=true 且空=粘性/池；非空=手动且重试不换。默认 ""
	MaxRetries       int               // 最大尝试次数；<=0 用内置默认 3。默认 0→3
	Timeout          int               // 超时秒；<=0 用内置默认 10。默认 0→10
	CloseMergeCookie bool              // false=自动合并 Cookie；true=关闭合并。默认 false
}

// Client HTTP 客户端。
type Client struct {
	mu sync.Mutex

	proxyPoolInst   proxyPool
	defaultTimeoutSec   int
	defaultRetries int
}

func newClient() *Client {
	return &Client{
		defaultTimeoutSec: 10,
		defaultRetries:    3,
	}
}

func defaultClient() *Client {
	defaultClientOnce.Do(func() {
		defaultClientInst = newClient()
	})
	return defaultClientInst
}

func normalizeOptions(opts []SendOptions) SendOptions {
	if len(opts) == 0 {
		return SendOptions{}
	}
	return opts[0]
}

func (c *Client) SetProxyPool(pool proxyPool) {
	c.mu.Lock()
	c.proxyPoolInst = pool
	c.mu.Unlock()
}

// NewSession 基于当前客户端配置创建会话（Cookie 持久 + 代理粘性）。
func (c *Client) NewSession() *Session {
	c.mu.Lock()
	defer c.mu.Unlock()
	return &Session{
		client: c,
		Cookie: make(map[string]string),
	}
}

// Session 保持 Cookie 与粘性代理（仅记录池代理）。
type Session struct {
	client   *Client
	mu       sync.Mutex
	Cookie   map[string]string
	stickyProxy string
}

func (s *Session) Send(method, url string, opts ...SendOptions) (*Response, error) {
	return s.client.sendInternal(s, method, url, normalizeOptions(opts))
}

// Send 单次请求（无会话粘性）。
// UseProxy=false 强制直连；=true 且 Proxy 非空 用手动代理；=true 且 Proxy 空 则粘性/池/直连。
func (c *Client) Send(method, url string, opts ...SendOptions) (*Response, error) {
	return c.sendInternal(nil, method, url, normalizeOptions(opts))
}

func (c *Client) sendInternal(session *Session, method, url string, opts SendOptions) (*Response, error) {
	url = strings.TrimSpace(url)
	if url == "" {
		return nil, errEmptyURL
	}
	method = strings.ToUpper(strings.TrimSpace(method))
	if method == "" {
		method = http.MethodGet
	}
	manualProxy := strings.TrimSpace(opts.Proxy)

	c.mu.Lock()
	defaultTimeout := c.defaultTimeoutSec
	defaultRetry := c.defaultRetries
	pool := c.proxyPoolInst
	c.mu.Unlock()

	timeoutSec := opts.Timeout
	maxAttempts := opts.MaxRetries
	if timeoutSec <= 0 {
		timeoutSec = defaultTimeout
	}
	if maxAttempts <= 0 {
		maxAttempts = defaultRetry
	}
	timeout := time.Duration(timeoutSec) * time.Second
	mergeCookie := !opts.CloseMergeCookie

	var lastErr error
	var lastResp *Response
	curProxy := ""

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		selected, fromManual, pickErr := c.pickProxy(session, opts.UseProxy, manualProxy, pool, curProxy, timeoutSec)
		if pickErr != nil {
			lastErr = pickErr
			continue
		}
		curProxy = selected

		resp, err := c.doOnce(session, method, url, opts, timeout, mergeCookie, selected)
		lastResp = resp
		lastErr = err

		if err == nil {
			// 仅池/粘性成功才写会话粘性；手动不污染粘性
			if session != nil && selected != "" && !fromManual {
				session.mu.Lock()
				session.stickyProxy = selected
				session.mu.Unlock()
			}
			return resp, nil
		}

		if !isRetryable(err, resp) || attempt >= maxAttempts {
			break
		}
		if fromManual {
			// 手动失败重试：仍用同一手动
			curProxy = manualProxy
			continue
		}
		// 池/粘性失败：清空粘性，换新 IP
		if session != nil {
			session.mu.Lock()
			if session.stickyProxy == curProxy {
				session.stickyProxy = ""
			}
			session.mu.Unlock()
		}
		curProxy = ""
	}
	if lastResp != nil {
		return lastResp, lastErr
	}
	return nil, lastErr
}

// 选择代理 返回 (代理URL, 是否来自手动, error)。
func (c *Client) pickProxy(session *Session, useProxy bool, manualProxy string, pool proxyPool, curProxy string, timeoutSec int) (string, bool, error) {
	if !useProxy {
		return "", false, nil
	}
	if manualProxy != "" {
		return manualProxy, true, nil
	}
	if curProxy != "" {
		return curProxy, false, nil
	}
	if session != nil {
		session.mu.Lock()
		sticky := session.stickyProxy
		session.mu.Unlock()
		if sticky != "" {
			return sticky, false, nil
		}
	}
	if pool == nil {
		return "", false, nil // 无池则直连
	}
	wait := timeoutSec
	if wait <= 0 {
		wait = 5
	}
	got, err := pool.GetProxy(wait)
	return got, false, err
}

func (c *Client) doOnce(session *Session, method, url string, opts SendOptions, timeout time.Duration, mergeCookie bool, proxy string) (*Response, error) {
	startTime := time.Now()

	body, contentType, err := encodeData(opts.Data)
	if err != nil {
		return nil, err
	}

	reqCookies := parseCookieString(opts.Cookie)
	sessionCookies := map[string]string{}
	if session != nil {
		session.mu.Lock()
		sessionCookies = mergeCookieMaps(session.Cookie)
		session.mu.Unlock()
	}
	merged := reqCookies
	if mergeCookie {
		merged = mergeCookieMaps(sessionCookies, reqCookies)
	}
	cookieStr := cookieMapToString(merged)

	var transport http.RoundTripper
	if strings.TrimSpace(proxy) == "" {
		transport = directTransport()
	} else {
		proxyTransport, err := newProxyTransport(proxy)
		if err != nil {
			return nil, err
		}
		transport = proxyTransport
	}
	client := &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
	if opts.NoRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(context.Background(), method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	for k, value := range opts.Headers {
		req.Header.Set(k, value)
	}
	if contentType != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentType)
	}
	if cookieStr != "" {
		req.Header.Set("Cookie", cookieStr)
	}

	httpResp, err := client.Do(req)
	elapsed := time.Since(startTime).Milliseconds()
	if elapsed <= 0 {
		elapsed = 1 // 亚毫秒本机请求仍记 1ms，避免 GetElapsedMs 为 0
	}
	if err != nil {
		return &Response{
			useProxy:   proxy,
			elapsedMs:   elapsed,
			allCookie: cookieStr,
		}, err
	}
	defer httpResp.Body.Close()
	payload, readErr := io.ReadAll(io.LimitReader(httpResp.Body, 32<<20))
	if readErr != nil {
		err = readErr
	}

	setCookieMap := cookiesFromResponse(httpResp.Header)
	setCookieStr := cookieMapToString(setCookieMap)
	all := merged
	if mergeCookie {
		all = mergeCookieMaps(merged, setCookieMap)
		if session != nil {
			session.mu.Lock()
			session.Cookie = mergeCookieMaps(session.Cookie, setCookieMap)
			session.mu.Unlock()
		}
	}

	redirect := ""
	if opts.NoRedirect {
		redirect = parseRedirect(httpResp)
	} else if httpResp.Request != nil && httpResp.Request.URL != nil {
		redirect = httpResp.Request.URL.String()
	}

	result := &Response{
		statusCode:     httpResp.StatusCode,
		body:       payload,
		redirectURL: redirect,
		setCookie: setCookieStr,
		allCookie: cookieMapToString(all),
		headers:     httpResp.Header.Clone(),
		useProxy:   proxy,
		elapsedMs:   elapsed,
	}
	if err != nil {
		return result, err
	}
	if isRetryableStatus(httpResp.StatusCode) {
		return result, errors.New("HTTP状态码可重试: " + httpResp.Status)
	}
	return result, nil
}

func encodeData(data any) ([]byte, string, error) {
	if data == nil {
		return nil, "", nil
	}
	switch v := data.(type) {
	case []byte:
		return v, "", nil
	case string:
		if v == "" {
			return nil, "", nil
		}
		return []byte(v), "", nil
	case map[string]string:
		value := url.Values{}
		for k, oneItem := range v {
			value.Set(k, oneItem)
		}
		return []byte(value.Encode()), "application/x-www-form-urlencoded", nil
	default:
		return nil, "", errors.New("不支持的data类型")
	}
}

func parseRedirect(resp *http.Response) string {
	location := resp.Header.Get("Location")
	if location == "" {
		return ""
	}
	ref, err := url.Parse(location)
	if err != nil {
		return location
	}
	if ref.IsAbs() {
		return ref.String()
	}
	if resp.Request != nil && resp.Request.URL != nil {
		return resp.Request.URL.ResolveReference(ref).String()
	}
	return location
}

func isRetryableStatus(code int) bool {
	return code == 502 || code == 503 || code == 504
}

func isRetryable(err error, resp *Response) bool {
	if err == nil {
		return false
	}
	if resp != nil && isRetryableStatus(resp.statusCode) {
		return true
	}
	msg := strings.ToLower(err.Error())
	keyword := []string{
		"timeout", "timed out", "deadline exceeded",
		"connection refused", "connection reset", "i/o timeout",
		"proxy", "eof", "broken pipe", "tls", "handshake",
		"temporary", "unavailable", "可重试",
	}
	for _, k := range keyword {
		if strings.Contains(msg, k) {
			return true
		}
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout() || netErr.Temporary()
	}
	return false
}
