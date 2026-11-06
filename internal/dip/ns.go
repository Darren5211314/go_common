package dip

// IPNS 对外命名空间类型（根包通过类型别名导出）。
type IPNS struct{}

func (IPNS) D是否有效(text string) bool {
	return isValidIP(text)
}

func (IPNS) D取版本(text string) string {
	return ipVersion(text)
}

func (IPNS) D是否私有(text string) bool {
	return isPrivate(text)
}

func (IPNS) D是否回环(text string) bool {
	return isLoopback(text)
}

func (IPNS) D取随机IP() string {
	return randomIP()
}

func (IPNS) D取本机IP() string {
	return localIP()
}

// D取外网IP 按接口编号单源取公网 IP；未知编号或失败返回 error。
func (IPNS) D取外网IP(iface int) (string, error) {
	return publicIP(iface)
}

func (IPNS) D加载数据文件(path string) error {
	return loadXDB(path)
}

// D加载默认数据文件 加载 data/ip2region_v4.xdb（工作目录或可执行文件旁）。
func (IPNS) D加载默认数据文件() error {
	return loadDefaultXDB()
}

func (IPNS) D查询(ipText string) (string, error) {
	return query(ipText)
}
