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
	错误未加载 = errors.New("ip2region数据未加载，请先调用D加载数据文件")
	查询锁   sync.RWMutex
	查询器   *xdb.Searcher
)

func 解析IP(文本 string) net.IP {
	return net.ParseIP(文本)
}

// D是否有效 判断是否为合法 IPv4/IPv6。
func 是否有效(文本 string) bool {
	return 解析IP(文本) != nil
}

// D取版本 返回 "IPv4" / "IPv6"；非法返回空串。
func 取版本(文本 string) string {
	ip := 解析IP(文本)
	if ip == nil {
		return ""
	}
	if ip.To4() != nil {
		return "IPv4"
	}
	return "IPv6"
}

// D是否私有 判断是否为私有地址；非法返回 false。
func 是否私有(文本 string) bool {
	ip := 解析IP(文本)
	if ip == nil {
		return false
	}
	return ip.IsPrivate()
}

// D是否回环 判断是否为回环地址；非法返回 false。
func 是否回环(文本 string) bool {
	ip := 解析IP(文本)
	if ip == nil {
		return false
	}
	return ip.IsLoopback()
}

func 是可用公网IPv4(ip net.IP) bool {
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
func 取随机IP() string {
	缓冲 := make([]byte, 4)
	for {
		if _, 错误 := rand.Read(缓冲); 错误 != nil {
			return ""
		}
		ip := net.IPv4(缓冲[0], 缓冲[1], 缓冲[2], 缓冲[3])
		if 是可用公网IPv4(ip) {
			return ip.String()
		}
	}
}

// D取本机IP 返回第一个非回环、非 link-local 的 IPv4；没有则 ""。
func 取本机IP() string {
	地址列表, 错误 := net.InterfaceAddrs()
	if 错误 != nil {
		return ""
	}
	for _, 地址 := range 地址列表 {
		ip网络, 成功 := 地址.(*net.IPNet)
		if !成功 {
			continue
		}
		ip := ip网络.IP
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

const 默认xdb相对路径 = "data/ip2region_v4.xdb"

// 加载数据文件 加载 ip2region IPv4 xdb 到内存（可重复调用覆盖）。
func 加载数据文件(路径 string) error {
	内容, 错误 := xdb.LoadContentFromFile(路径)
	if 错误 != nil {
		return 错误
	}
	新查询器, 错误 := xdb.NewWithBuffer(xdb.IPv4, 内容)
	if 错误 != nil {
		return 错误
	}
	查询锁.Lock()
	查询器 = 新查询器
	查询锁.Unlock()
	return nil
}

func 默认xdb候选路径() []string {
	列表 := []string{默认xdb相对路径}
	if 可执行, 错误 := os.Executable(); 错误 == nil {
		目录 := filepath.Dir(可执行)
		列表 = append(列表, filepath.Join(目录, 默认xdb相对路径))
	}
	return 列表
}

// 加载默认数据文件 按约定路径加载 data/ip2region_v4.xdb（工作目录，其次可执行文件旁）。
func 加载默认数据文件() error {
	候选 := 默认xdb候选路径()
	var 尝试 []string
	for _, 路径 := range 候选 {
		尝试 = append(尝试, 路径)
		if _, 错误 := os.Stat(路径); 错误 != nil {
			continue
		}
		return 加载数据文件(路径)
	}
	return fmt.Errorf("未找到 ip2region 数据文件，已尝试: %s", strings.Join(尝试, ", "))
}

// D查询 按已加载的 xdb 查询地区字符串；未加载或失败返回 error；未命中返回 ""。
func 查询(IP文本 string) (string, error) {
	查询锁.RLock()
	当前 := 查询器
	查询锁.RUnlock()
	if 当前 == nil {
		return "", 错误未加载
	}
	地区, 错误 := 当前.Search(IP文本)
	if 错误 != nil {
		return "", 错误
	}
	return 地区, nil
}
