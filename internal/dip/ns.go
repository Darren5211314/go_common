package dip

// IPNS 对外命名空间类型（根包通过类型别名导出）。
type IPNS struct{}

func (IPNS) Valid(text string) bool {
	return isValidIP(text)
}

func (IPNS) Version(text string) string {
	return ipVersion(text)
}

func (IPNS) IsPrivate(text string) bool {
	return isPrivate(text)
}

func (IPNS) IsLoopback(text string) bool {
	return isLoopback(text)
}

func (IPNS) Random() string {
	return randomIP()
}

func (IPNS) Local() string {
	return localIP()
}

// Public 按接口编号单源取公网 IP；未知编号或失败返回 error。
func (IPNS) Public(iface int) (string, error) {
	return publicIP(iface)
}

func (IPNS) LoadDB(path string) error {
	return loadXDB(path)
}

// LoadDefaultDB 加载 data/ip2region_v4.xdb（工作目录或可执行文件旁）。
func (IPNS) LoadDefaultDB() error {
	return loadDefaultXDB()
}

func (IPNS) Query(ipText string) (string, error) {
	return query(ipText)
}
