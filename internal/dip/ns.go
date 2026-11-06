package dip

// IP类 对外命名空间类型（根包通过类型别名导出）。
type IP类 struct{}

func (IP类) D是否有效(text string) bool {
	return isValidIP(text)
}

func (IP类) D取版本(text string) string {
	return ipVersion(text)
}

func (IP类) D是否私有(text string) bool {
	return isPrivate(text)
}

func (IP类) D是否回环(text string) bool {
	return isLoopback(text)
}

func (IP类) D取随机IP() string {
	return randomIP()
}

func (IP类) D取本机IP() string {
	return localIP()
}

// D取外网IP 按接口编号单源取公网 IP；未知编号或失败返回 error。
func (IP类) D取外网IP(iface int) (string, error) {
	return publicIP(iface)
}

func (IP类) D加载数据文件(path string) error {
	return loadXDB(path)
}

// D加载默认数据文件 加载 data/ip2region_v4.xdb（工作目录或可执行文件旁）。
func (IP类) D加载默认数据文件() error {
	return loadDefaultXDB()
}

func (IP类) D查询(ipText string) (string, error) {
	return query(ipText)
}
