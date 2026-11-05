package dip

// IP类 对外命名空间类型（根包通过类型别名导出）。
type IP类 struct{}

func (IP类) D是否有效(文本 string) bool {
	return 是否有效(文本)
}

func (IP类) D取版本(文本 string) string {
	return 取版本(文本)
}

func (IP类) D是否私有(文本 string) bool {
	return 是否私有(文本)
}

func (IP类) D是否回环(文本 string) bool {
	return 是否回环(文本)
}

func (IP类) D取随机IP() string {
	return 取随机IP()
}

func (IP类) D取本机IP() string {
	return 取本机IP()
}

// D取外网IP 按接口编号单源取公网 IP；未知编号或失败返回 error。
func (IP类) D取外网IP(接口 int) (string, error) {
	return 取外网IP(接口)
}

func (IP类) D加载数据文件(路径 string) error {
	return 加载数据文件(路径)
}

// D加载默认数据文件 加载 data/ip2region_v4.xdb（工作目录或可执行文件旁）。
func (IP类) D加载默认数据文件() error {
	return 加载默认数据文件()
}

func (IP类) D查询(IP文本 string) (string, error) {
	return 查询(IP文本)
}
