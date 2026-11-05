package dproxy

import "testing"

func Test解析代理行(t *testing.T) {
	得到, err := 解析代理行("1.2.3.4:8080", "http")
	if err != nil || 得到 != "http://1.2.3.4:8080" {
		t.Errorf("默认协议解析不符: %v %v", 得到, err)
	}
	得到, err = 解析代理行("socks5://u:p@1.2.3.4:1080", "http")
	if err != nil || 得到 != "socks5://u:p@1.2.3.4:1080" {
		t.Errorf("socks5 解析不符: %v %v", 得到, err)
	}
	得到, err = 解析代理行("https://1.2.3.4:8443", "http")
	if err != nil || 得到 != "https://1.2.3.4:8443" {
		t.Errorf("https 解析不符: %v %v", 得到, err)
	}
}
