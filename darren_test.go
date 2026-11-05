package darren_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Darren5211314/go_common"
)

func TestD文本门面转发(t *testing.T) {
	darren.Error("啊是哒是哒大")
	得到 := darren.D文本.D取文本左边("你好_世界", "_")
	darren.Info("取文本左边", "输入", "你好_世界", "标记", "_", "结果", 得到, "期望", "你好")
	if 得到 != "你好" {
		t.Errorf("取文本左边 不符")
	}

	得到 = darren.D文本.D取文本右边("你好_世界", "_")
	darren.Info("取文本右边", "输入", "你好_世界", "标记", "_", "结果", 得到, "期望", "世界")
	if 得到 != "世界" {
		t.Errorf("取文本右边 不符")
	}

	得到 = darren.D文本.D取左边("你好世界", 2)
	darren.Info("取左边", "输入", "你好世界", "长度", 2, "结果", 得到, "期望", "你好")
	if 得到 != "你好" {
		t.Errorf("取左边 不符")
	}

	得到 = darren.D文本.D取右边("你好世界", 2)
	darren.Info("取右边", "输入", "你好世界", "长度", 2, "结果", 得到, "期望", "世界")
	if 得到 != "世界" {
		t.Errorf("取右边 不符")
	}

	长度 := darren.D文本.D取文本长度("你好")
	darren.Info("取文本长度", "输入", "你好", "结果", 长度, "期望", 2)
	if 长度 != 2 {
		t.Errorf("取文本长度 不符")
	}

	位置 := darren.D文本.D寻找文本("你好世界", "世界")
	darren.Info("寻找文本", "输入", "你好世界", "子文本", "世界", "结果", 位置, "期望", 3)
	if 位置 != 3 {
		t.Errorf("寻找文本 不符")
	}

	得到 = darren.D文本.D删首空("  a")
	darren.Info("删首空", "输入", "  a", "结果", 得到, "期望", "a")
	if 得到 != "a" {
		t.Errorf("删首空 不符")
	}

	得到 = darren.D文本.D删尾空("a  ")
	darren.Info("删尾空", "输入", "a  ", "结果", 得到, "期望", "a")
	if 得到 != "a" {
		t.Errorf("删尾空 不符")
	}

	得到 = darren.D文本.D删首尾空(" a ")
	darren.Info("删首尾空", "输入", " a ", "结果", 得到, "期望", "a")
	if 得到 != "a" {
		t.Errorf("删首尾空 不符")
	}

	得到 = darren.D文本.D删全部空("a b")
	darren.Info("删全部空", "输入", "a b", "结果", 得到, "期望", "ab")
	if 得到 != "ab" {
		t.Errorf("删全部空 不符")
	}

	得到 = darren.D文本.D子文本替换("a_b", "_", "-")
	darren.Info("子文本替换", "输入", "a_b", "原文本", "_", "新文本", "-", "结果", 得到, "期望", "a-b")
	if 得到 != "a-b" {
		t.Errorf("子文本替换 不符")
	}

	得到 = darren.D文本.D取文本中间("你好世界", 2, 2)
	darren.Info("取文本中间", "输入", "你好世界", "起始", 2, "长度", 2, "结果", 得到, "期望", "好世")
	if 得到 != "好世" {
		t.Errorf("取文本中间 不符")
	}

	得到 = darren.D文本.D取文本中间文本("a[123]b", "[", "]")
	darren.Info("取文本中间文本", "输入", "a[123]b", "左边", "[", "右边", "]", "结果", 得到, "期望", "123")
	if 得到 != "123" {
		t.Errorf("取文本中间文本 不符")
	}

	是字母 := darren.D文本.D是否为字母("Ab")
	是数字 := darren.D文本.D是否为数字("12")
	是汉字 := darren.D文本.D是否为汉字("你好")
	darren.Info("类别判断", "是否为字母", 是字母, "是否为数字", 是数字, "是否为汉字", 是汉字)
	if !是字母 || !是数字 || !是汉字 {
		t.Errorf("类别判断 不符")
	}

	得到 = darren.D文本.D只取汉字("a你b好")
	darren.Info("只取汉字", "输入", "a你b好", "结果", 得到, "期望", "你好")
	if 得到 != "你好" {
		t.Errorf("只取汉字 不符")
	}

	随机 := darren.D文本.D取随机字符(8, false, false, false)
	darren.Info("取随机字符", "位数", 8, "结果", 随机, "期望长度", 8)
	if len(随机) != 8 {
		t.Errorf("取随机字符 长度不符")
	}
}

func TestD时间门面转发(t *testing.T) {
	前秒 := time.Now().Unix()
	十位 := darren.D时间.D取现行时间戳(true)
	后秒 := time.Now().Unix()
	darren.Info("取现行时间戳", "取十位", true, "结果", 十位, "范围起", 前秒, "范围止", 后秒)
	if 十位 < 前秒 || 十位 > 后秒 {
		t.Errorf("取现行时间戳(true) 不符")
	}

	前毫秒 := time.Now().UnixMilli()
	十三位 := darren.D时间.D取现行时间戳(false)
	后毫秒 := time.Now().UnixMilli()
	darren.Info("取现行时间戳", "取十位", false, "结果", 十三位, "范围起", 前毫秒, "范围止", 后毫秒)
	if 十三位 < 前毫秒 || 十三位 > 后毫秒 {
		t.Errorf("取现行时间戳(false) 不符")
	}

	固定 := time.Date(2026, 11, 5, 20, 53, 0, 0, time.Local)
	转时间 := darren.D时间.D时间戳转时间(固定.UnixMilli(), "yyyy-mm-dd hh:mm:ss")
	darren.Info("时间戳转时间", "戳", 固定.UnixMilli(), "格式", "yyyy-mm-dd hh:mm:ss", "结果", 转时间, "期望", "2026-11-05 20:53:00")
	if 转时间 != "2026-11-05 20:53:00" {
		t.Errorf("时间戳转时间 不符")
	}

	秒转 := darren.D时间.D时间戳转时间(固定.Unix(), "yyyy-mm-dd")
	darren.Info("时间戳转时间", "戳", 固定.Unix(), "格式", "yyyy-mm-dd", "结果", 秒转, "期望", "2026-11-05")
	if 秒转 != "2026-11-05" {
		t.Errorf("时间戳转时间(秒) 不符")
	}

	转戳 := darren.D时间.D时间转时间戳("2026-11-05 20:53:00", true)
	darren.Info("时间转时间戳", "输入", "2026-11-05 20:53:00", "取十位", true, "结果", 转戳, "期望", 固定.Unix())
	if 转戳 != 固定.Unix() {
		t.Errorf("时间转时间戳 不符")
	}

	格式化 := darren.D时间.D时间格式化("2026/11/05", "yyyy-mm-dd")
	darren.Info("时间格式化", "输入", "2026/11/05", "目标格式", "yyyy-mm-dd", "结果", 格式化, "期望", "2026-11-05")
	if 格式化 != "2026-11-05" {
		t.Errorf("时间格式化 不符")
	}

	开始 := time.Now()
	darren.D时间.D延时(30)
	耗时 := time.Since(开始)
	darren.Info("延时", "毫秒", 30, "耗时", 耗时)
	if 耗时 < 20*time.Millisecond {
		t.Errorf("D延时 耗时过短: %v", 耗时)
	}
}

func Test日志包直接调用(t *testing.T) {
	darren.Info("直接log.Info", "说明", "无需先初始化")
	darren.Warn("直接log.Warn", "级别", "警告")
	darren.Error("直接log.Error", "级别", "错误")
	darren.Infof("Infof示例 值=%d", 42)
	darren.Errorf("Errorf示例 码=%d", 500)
}

func TestD文件门面转发(t *testing.T) {
	目录 := t.TempDir()
	路径 := filepath.Join(目录, "子", "demo.txt")
	if err := darren.D文件.D保存文本(路径, "", "a", "b"); err != nil {
		t.Fatalf("保存失败: %v", err)
	}
	if !darren.D文件.D文件是否存在(路径) {
		t.Fatal("文件应存在")
	}
	if darren.D文件.D取文件名(路径, false) != "demo" {
		t.Errorf("文件名不符")
	}
	大小, err := darren.D文件.D取文件大小(路径, "B")
	darren.Info("D文件", "大小", 大小, "扩展名", darren.D文件.D取扩展名(路径))
	if err != nil || 大小 <= 0 {
		t.Errorf("大小不符: %v %v", 大小, err)
	}
	目标 := filepath.Join(目录, "copy.txt")
	if err := darren.D文件.D复制文件(路径, 目标); err != nil {
		t.Fatalf("复制失败: %v", err)
	}
	if err := darren.D文件.D删除到回收站(目标); err != nil {
		t.Fatalf("回收站失败: %v", err)
	}
	if darren.D文件.D文件是否存在(目标) {
		t.Fatal("回收站后不应存在")
	}
}

func Test加密门面转发(t *testing.T) {
	摘要 := darren.Hash.MD5("abc")
	darren.Info("Hash.MD5", "结果", 摘要, "期望", "900150983cd24fb0d6963f7d28e17f72")
	if 摘要 != "900150983cd24fb0d6963f7d28e17f72" {
		t.Errorf("Hash.MD5 不符")
	}

	sm3期望 := "66c7f0f462eeedd9d1f2d46bdc10e4e24167c4875cf2f7a2297da02b8f4ba8e0"
	sm3摘要 := darren.Hash.SM3("abc")
	darren.Info("Hash.SM3", "结果", sm3摘要, "期望", sm3期望)
	if sm3摘要 != sm3期望 {
		t.Errorf("Hash.SM3 不符")
	}

	mac, err := darren.HMAC.SHA256("abc", "key")
	darren.Info("HMAC.SHA256", "结果", mac, "错误", err)
	if err != nil || mac == "" {
		t.Errorf("HMAC.SHA256 失败")
	}

	crc := darren.CRC.CRC32("123456789")
	darren.Info("CRC.CRC32", "结果", crc, "期望", uint32(0xcbf43926))
	if crc != 0xcbf43926 {
		t.Errorf("CRC.CRC32 不符")
	}

	密钥 := "00112233445566778899aabbccddeeff"
	IV := "0102030405060708090a0b0c0d0e0f10"
	nonce := "0102030405060708090a0b0c"
	明文 := "hello-aes"

	密文, err := darren.AES.EncryptCBC(明文, 密钥, IV)
	解, err2 := darren.AES.DecryptCBC(密文, 密钥, IV)
	darren.Info("AES.CBC往返", "密文", 密文, "明文", 解, "错误1", err, "错误2", err2)
	if err != nil || err2 != nil || 解 != 明文 {
		t.Errorf("AES.CBC 往返不符")
	}

	密文, err = darren.AES.EncryptECB(明文, 密钥)
	解, err2 = darren.AES.DecryptECB(密文, 密钥)
	darren.Info("AES.ECB往返", "密文", 密文, "明文", 解)
	if err != nil || err2 != nil || 解 != 明文 {
		t.Errorf("AES.ECB 往返不符")
	}

	密文, err = darren.AES.EncryptGCM(明文, 密钥, nonce)
	解, err2 = darren.AES.DecryptGCM(密文, 密钥, nonce)
	darren.Info("AES.GCM往返", "密文", 密文, "明文", 解)
	if err != nil || err2 != nil || 解 != 明文 {
		t.Errorf("AES.GCM 往返不符")
	}

	des密钥 := "0011223344556677"
	desIV := "0102030405060708"
	密文, err = darren.DES.EncryptCBC("hello-des", des密钥, desIV)
	解, err2 = darren.DES.DecryptCBC(密文, des密钥, desIV)
	darren.Info("DES.CBC往返", "密文", 密文, "明文", 解)
	if err != nil || err2 != nil || 解 != "hello-des" {
		t.Errorf("DES.CBC 往返不符")
	}

	des3密钥 := "00112233445566778899aabbccddeeff0011223344556677"
	密文, err = darren.DES3.EncryptCBC("hello-3des", des3密钥, desIV)
	解, err2 = darren.DES3.DecryptCBC(密文, des3密钥, desIV)
	darren.Info("DES3.CBC往返", "密文", 密文, "明文", 解)
	if err != nil || err2 != nil || 解 != "hello-3des" {
		t.Errorf("DES3.CBC 往返不符")
	}

	pub, priv, err := darren.RSA.GenerateKeyPair(1024)
	darren.Info("RSA.GenerateKeyPair", "错误", err, "公钥长度", len(pub), "私钥长度", len(priv))
	if err != nil {
		t.Fatalf("RSA 生成失败: %v", err)
	}
	密文, err = darren.RSA.Encrypt("hello-rsa", pub)
	解, err2 = darren.RSA.Decrypt(密文, priv)
	darren.Info("RSA加解密", "密文", 密文, "明文", 解, "错误1", err, "错误2", err2)
	if err != nil || err2 != nil || 解 != "hello-rsa" {
		t.Errorf("RSA 加解密不符")
	}
	签名, err := darren.RSA.Sign("hello-rsa", priv)
	验签 := darren.RSA.Verify("hello-rsa", 签名, pub)
	darren.Info("RSA签名", "签名", 签名, "验签", 验签, "错误", err)
	if err != nil || 验签 != nil {
		t.Errorf("RSA 签名验签不符")
	}

	sm4密钥 := "0123456789abcdeffedcba9876543210"
	sm4IV := "000102030405060708090a0b0c0d0e0f"
	sm4明文 := "hello-sm4"
	密文, err = darren.SM4.EncryptCBC(sm4明文, sm4密钥, sm4IV)
	解, err2 = darren.SM4.DecryptCBC(密文, sm4密钥, sm4IV)
	darren.Info("SM4.CBC往返", "密文", 密文, "明文", 解, "错误1", err, "错误2", err2)
	if err != nil || err2 != nil || 解 != sm4明文 {
		t.Errorf("SM4.CBC 往返不符")
	}
	密文, err = darren.SM4.EncryptECB(sm4明文, sm4密钥)
	解, err2 = darren.SM4.DecryptECB(密文, sm4密钥)
	darren.Info("SM4.ECB往返", "密文", 密文, "明文", 解)
	if err != nil || err2 != nil || 解 != sm4明文 {
		t.Errorf("SM4.ECB 往返不符")
	}
	密文, err = darren.SM4.EncryptGCM(sm4明文, sm4密钥, nonce)
	解, err2 = darren.SM4.DecryptGCM(密文, sm4密钥, nonce)
	darren.Info("SM4.GCM往返", "密文", 密文, "明文", 解)
	if err != nil || err2 != nil || 解 != sm4明文 {
		t.Errorf("SM4.GCM 往返不符")
	}
}

func TestJSON门面转发(t *testing.T) {
	样本 := `{
  "code": 0,
  "msg": "ok",
  "data": {
    "user": {"id": 1001, "name": "张三", "vip": true},
    "tags": ["go", "json"],
    "score": 98.5
  }
}`
	if !darren.JSON.D是否有效(样本) {
		t.Errorf("D是否有效 不符")
	}
	名 := darren.JSON.D取文本(样本, "data.user.name")
	id := darren.JSON.D取整数(样本, "data.user.id")
	分 := darren.JSON.D取小数(样本, "data.score")
	vip := darren.JSON.D取逻辑(样本, "data.user.vip")
	darren.Info("JSON取值", "name", 名, "id", id, "score", 分, "vip", vip)
	if 名 != "张三" || id != 1001 || 分 != 98.5 || !vip {
		t.Errorf("JSON 取值不符")
	}
	if darren.JSON.D取成员数(样本, "data.tags") != 2 {
		t.Errorf("D取成员数 不符")
	}
	新, err := darren.JSON.D设置文本(样本, "data.user.name", "李四")
	darren.Info("JSON设置", "结果", 新, "错误", err)
	if err != nil || darren.JSON.D取文本(新, "data.user.name") != "李四" {
		t.Errorf("D设置文本 不符")
	}
	紧凑 := darren.JSON.D压缩(样本)
	if 紧凑 == "" || !darren.JSON.D是否有效(紧凑) {
		t.Errorf("D压缩 不符")
	}
}

func TestQR门面转发(t *testing.T) {
	原文 := "你好二维码-门面"
	b64, err := darren.QR.D生成(原文, 200)
	darren.Info("QR.D生成", "长度", len(b64), "错误", err)
	if err != nil {
		t.Fatalf("D生成失败: %v", err)
	}
	解出, err := darren.QR.D解析Base64(b64)
	darren.Info("QR.D解析Base64", "结果", 解出, "错误", err)
	if err != nil || 解出 != 原文 {
		t.Errorf("QR Base64往返不符")
	}

	路径 := filepath.Join(t.TempDir(), "qr-facade.png")
	if err := darren.QR.D生成到文件(原文, 路径, 0); err != nil {
		t.Fatalf("D生成到文件失败: %v", err)
	}
	解出, err = darren.QR.D解析文件(路径)
	darren.Info("QR.D解析文件", "结果", 解出, "错误", err)
	if err != nil || 解出 != 原文 {
		t.Errorf("QR 文件往返不符")
	}
}

func TestIP门面转发(t *testing.T) {
	随机 := darren.IP.D取随机IP()
	本机 := darren.IP.D取本机IP()
	darren.Info("IP门面", "随机", 随机, "本机", 本机, "版本", darren.IP.D取版本(随机))
	if !darren.IP.D是否有效(随机) || darren.IP.D取版本(随机) != "IPv4" {
		t.Errorf("D取随机IP 不符")
	}
	if !darren.IP.D是否私有("10.0.0.1") || !darren.IP.D是否回环("127.0.0.1") {
		t.Errorf("私有/回环判断不符")
	}
	if _, err := darren.IP.D查询("8.8.8.8"); err == nil {
		t.Errorf("未加载应报错")
	}
	xdb := os.Getenv("IP2REGION_XDB")
	if xdb == "" {
		t.Log("未设置 IP2REGION_XDB，跳过 D查询门面断言")
		return
	}
	if err := darren.IP.D加载数据文件(xdb); err != nil {
		t.Fatalf("D加载数据文件失败: %v", err)
	}
	地区, err := darren.IP.D查询("8.8.8.8")
	darren.Info("IP.D查询", "地区", 地区, "错误", err)
	if err != nil || 地区 == "" {
		t.Errorf("D查询 不符")
	}
}

func TestProxy门面转发(t *testing.T) {
	源 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "10.9.9.1:8000\n10.9.9.2:8000\n")
	}))
	defer 源.Close()

	池 := darren.Proxy.D新建()
	池.D设置代理链接(源.URL)
	池.D设置代理池容量(2)
	池.D设置是否验证(false)
	池.D设置获取间隔(1)
	池.D设置默认协议("http")
	if err := 池.D启动(); err != nil {
		t.Fatalf("D启动失败: %v", err)
	}
	defer 池.D关闭()

	代理, err := 池.D获取代理(5)
	darren.Info("Proxy.D获取代理", "代理", 代理, "错误", err)
	if err != nil || 代理 == "" {
		t.Fatalf("D获取代理失败: %v", err)
	}
	统计 := 池.D取统计()
	darren.Info("Proxy.D取统计", "从链接", 统计.D从链接获取数, "已取出", 统计.D已取出数)
	if 统计.D已取出数 < 1 || 统计.D从链接获取数 < 1 {
		t.Errorf("统计不符: %+v", 统计)
	}
}

func TestThread门面转发(t *testing.T) {
	池 := darren.Thread.D新建(4)
	池.D设置队列容量(32)
	if err := 池.D启动(); err != nil {
		t.Fatalf("D启动失败: %v", err)
	}
	defer 池.D停止()

	任务, err := 池.D投递任务(func() (any, error) {
		return 42, nil
	})
	if err != nil {
		t.Fatalf("D投递任务失败: %v", err)
	}
	结果, err := 任务.D等待()
	darren.Info("Thread任务", "结果", 结果, "错误", err)
	if err != nil || 结果 != 42 {
		t.Errorf("任务结果不符")
	}

	var 完成 sync.WaitGroup
	完成.Add(5)
	for i := 0; i < 5; i++ {
		if err := 池.D投递(func() { 完成.Done() }); err != nil {
			t.Fatalf("D投递失败: %v", err)
		}
	}
	完成.Wait()
	统计 := 池.D取统计()
	darren.Info("Thread统计", "投递", 统计.D已投递数, "完成", 统计.D已完成数)
	if 统计.D已投递数 < 6 {
		t.Errorf("投递统计不符")
	}
}

func TestHTTP门面转发(t *testing.T) {
	服务 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Set-Cookie", "sid=1")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer 服务.Close()

	客户端 := darren.HTTP.D新建()
	响应, err := 客户端.D发送("GET", 服务.URL, nil, "", nil, false, false, "", 3, true, 5)
	if err != nil {
		t.Fatalf("D发送失败: %v", err)
	}
	darren.Info("HTTP响应", "状态", 响应.D取状态码(), "正文", 响应.D取正文(), "耗时", 响应.D取耗时())
	if 响应.D取状态码() != 200 || !strings.Contains(响应.D取正文(), "ok") {
		t.Errorf("门面响应不符")
	}
	if 响应.D取耗时() <= 0 {
		t.Errorf("耗时应 > 0")
	}

	会话 := 客户端.D新建会话()
	响应2, err := 会话.D发送("POST", 服务.URL, map[string]string{"a": "1"}, "token=abc", map[string]string{"User-Agent": "github.com/Darren5211314/go_common"}, false, false, "", 3, true, 5)
	if err != nil {
		t.Fatalf("会话D发送失败: %v", err)
	}
	if 响应2.D取状态码() != 200 {
		t.Errorf("会话状态码不符")
	}
}
