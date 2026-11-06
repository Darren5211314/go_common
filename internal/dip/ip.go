package dip

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

var (
	errNotLoaded = errors.New("ip2region数据未加载，请先调用D加载数据文件")
	searchMu   sync.RWMutex
	searcher   *xdb.Searcher
)

func parseIP(text string) net.IP {
	return net.ParseIP(text)
}

// D是否有效 判断是否为合法 IPv4/IPv6。
func isValidIP(text string) bool {
	return parseIP(text) != nil
}

// D取版本 返回 "IPv4" / "IPv6"；非法返回空串。
func ipVersion(text string) string {
	ip := parseIP(text)
	if ip == nil {
		return ""
	}
	if ip.To4() != nil {
		return "IPv4"
	}
	return "IPv6"
}

// D是否私有 判断是否为私有地址；非法返回 false。
func isPrivate(text string) bool {
	ip := parseIP(text)
	if ip == nil {
		return false
	}
	return ip.IsPrivate()
}

// D是否回环 判断是否为回环地址；非法返回 false。
func isLoopback(text string) bool {
	ip := parseIP(text)
	if ip == nil {
		return false
	}
	return ip.IsLoopback()
}

func isUsablePublicIPv4(ip net.IP) bool {
	v4 := ip.To4()
	if v4 == nil {
		return false
	}
	if ip.IsUnspecified() || ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsMulticast() {
		return false
	}
	// 100.64.0.0/10 CGNAT
	if v4[0] == 100 && v4[1] >= 64 && v4[1] <= 127 {
		return false
	}
	// 198.18.0.0/15 benchmarking
	if v4[0] == 198 && (v4[1] == 18 || v4[1] == 19) {
		return false
	}
	return true
}

// D取随机IP 生成随机公网风格 IPv4 字符串。
func randomIP() string {
	buf := make([]byte, 4)
	for {
		if _, err := rand.Read(buf); err != nil {
			return ""
		}
		ip := net.IPv4(buf[0], buf[1], buf[2], buf[3])
		if isUsablePublicIPv4(ip) {
			return ip.String()
		}
	}
}

// D取本机IP 返回第一个非回环、非 link-local 的 IPv4；没有则 ""。
func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		ip := ipNet.IP
		v4 := ip.To4()
		if v4 == nil {
			continue
		}
		if ip.IsLoopback() || ip.IsLinkLocalUnicast() {
			continue
		}
		return v4.String()
	}
	return ""
}

const defaultXDBRelPath = "data/ip2region_v4.xdb"

// 加载数据文件 加载 ip2region IPv4 xdb 到内存（可重复调用覆盖）。
func loadXDB(path string) error {
	content, err := xdb.LoadContentFromFile(path)
	if err != nil {
		return err
	}
	newSearcher, err := xdb.NewWithBuffer(xdb.IPv4, content)
	if err != nil {
		return err
	}
	searchMu.Lock()
	searcher = newSearcher
	searchMu.Unlock()
	return nil
}

func defaultXDBCandidates() []string {
	list := []string{defaultXDBRelPath}
	if runnable, err := os.Executable(); err == nil {
		dir := filepath.Dir(runnable)
		list = append(list, filepath.Join(dir, defaultXDBRelPath))
	}
	return list
}

// 加载默认数据文件 按约定路径加载 data/ip2region_v4.xdb（工作目录，其次可执行文件旁）。
func loadDefaultXDB() error {
	candidate := defaultXDBCandidates()
	var try []string
	for _, path := range candidate {
		try = append(try, path)
		if _, err := os.Stat(path); err != nil {
			continue
		}
		return loadXDB(path)
	}
	return fmt.Errorf("未找到 ip2region 数据文件，已尝试: %s", strings.Join(try, ", "))
}

// D查询 按已加载的 xdb 查询地区字符串；未加载或失败返回 error；未命中返回 ""。
func query(ipText string) (string, error) {
	searchMu.RLock()
	cur := searcher
	searchMu.RUnlock()
	if cur == nil {
		return "", errNotLoaded
	}
	region, err := cur.Search(ipText)
	if err != nil {
		return "", err
	}
	return region, nil
}
