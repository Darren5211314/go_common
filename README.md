# darren

面向 Windows exe 与 Ubuntu 等 Linux 项目的个人 Go 工具库。调用方只导入根包，通过中文命名空间使用方法。

## 环境

- Go 1.26（本机 SDK：`D:\DevAPP\Go\SDK\1.26.10`）
- Module：`github.com/Darren5211314/go_common`（根包名仍为 `darren`，调用 `darren.Info` / `darren.IP` 等）。
- 即使业务只用 goja，`go.mod` 仍会拉取 `rogchap.com/v8go`（含预编译库，体积大）。Windows 上勿启用 V8；后续可再拆成 opt-in build tag / 子模块。

## 调用示例

```go
package main

import "github.com/Darren5211314/go_common"

func main() {
	左边 := darren.D文本.D取文本左边("abc_123", "_")
	右边 := darren.D文本.D取文本右边("abc_123", "_")
	前两字 := darren.D文本.D取左边("你好世界", 2)
	位置 := darren.D文本.D寻找文本("你好世界", "世界") // 从 1 起算，结果为 3
	验证码 := darren.D文本.D取随机字符(6, false, false, false) // 默认小写+数字
	戳 := darren.D时间.D取现行时间戳(false) // false=13位毫秒，true=10位秒
	文本 := darren.D时间.D时间戳转时间(戳, "yyyy-mm-dd hh:mm:ss")
	摘要 := darren.Hash.MD5("abc")
	密钥 := "00112233445566778899aabbccddeeff"
	IV := "0102030405060708090a0b0c0d0e0f10"
	data, err := darren.AES.EncryptCBC("hello", 密钥, IV)

	darren.Info("示例结果",
		"左边", 左边,
		"右边", 右边,
		"前两字", 前两字,
		"位置", 位置,
		"验证码", 验证码,
		"戳", 戳,
		"文本", 文本,
		"摘要", 摘要,
		"data", data,
		"错误", err,
	)
	darren.Errorf("printf风格示例 码=%d", 500)
}
```

## 约定

- **导出名前缀 `D`**：Go 规定跨包可见的标识符必须首字母为大写拉丁字母；纯中文无法导出，因此统一写成 `D文本`、`D取文本左边`。
- 命名空间：中文 `D文本`、`D时间`、`D文件`；英文不加 `D`：`AES`、`DES`、`DES3`、`SM4`、`Hash`（含 `SM3`）、`HMAC`、`CRC`、`RSA`、`JSON`、`QR`、`IP`、`Proxy`、`Thread`、`HTTP`、`JS`、`Sync`、`Codec`（后九者方法仍中文加 `D`）。日志用根包 `darren.Info` / `Warn` / `Error` / `Infof` 等。以后可加 `D系统`。
- 注释使用中文；参数名可用中文（包内局部，不必加 D）。
- 字符串按 UTF-8 字符处理，不按字节切。
- `D寻找文本`、`D取文本中间` 的位置从 **1** 起算；找不到返回 `0` / 空串。
- `D取随机字符(位数, 小写字母, 大写字母, 数字)`：三者都为 `false` 时默认小写字母 + 数字。
- `D文本` 找不到内容或空输入时返回空字符串，不返回 `error`。加解密、HMAC、RSA 返回 `(结果, error)`；Hash/CRC 只返回值。
- `D取现行时间戳(取十位)`：`true`=秒（约 10 位），`false`=毫秒（约 13 位）。
- 时间输出格式只用通配符：`yyyy-mm-dd hh:mm:ss`。
- `D时间戳转时间` 自动识别 10/13 位；`D时间转时间戳(文本, 取十位)` / `D时间格式化(文本, 目标格式)` **自动识别输入**，不必传原格式；空文本格式化=现行时间。
- `D延时(毫秒)`：阻塞当前协程；`<=0` 立即返回。
- **D文件**：`D取文件大小(路径, 单位)`（B/KB/MB/GB/TB，1024 进制）；`D删除文件` / `D删除到回收站`（Windows + Linux Trash）；`D移动文件` / `D复制文件` / `D重命名`；`D取扩展名`（含点）；`D文件是否存在` / `D目录是否存在`；`D取文件名(路径, 带扩展名)`；`D保存文本(路径, 分隔符, ...内容)`（空分隔符默认 `----`，自动建父目录并覆盖）。
- **加解密**：模式写进方法名（`EncryptCBC` / `EncryptECB` / `EncryptGCM`）；密钥、IV/nonce、data、摘要均为十六进制字符串；CBC/ECB 用 PKCS7；GCM 返回「data||tag」hex，nonce 单独传。国密：`Hash.SM3`、`SM4`（128-bit 密钥，CBC/ECB/GCM）；依赖 `github.com/emmansun/gmsm`。
- **JSON**：路径用点号，数组下标为数字段（`data.user.name`、`data.tags.0`）。取值找不到返回空/`0`/`false`；`D设置*` 返回新字符串并自动创建中间节点；`D压缩` / `D格式化` 非法返回空串。
- **QR**：`D生成` 返回 PNG Base64；`D生成到文件` 写 PNG；`D解析文件` / `D解析Base64` 读回文本。边长 `<=0` 默认 256；纠错固定 Medium。依赖 `skip2/go-qrcode`、`makiuchi-d/gozxing`。
- **IP**：`D取随机IP` / `D是否有效` / `D取版本` / `D是否私有` / `D是否回环` / `D取本机IP`（网卡局域网）/ `D取外网IP(接口) (string, error)`（按编号单源；1=3322 2=gdt 3=qq-checktime 4=pconline 5=ipinfo 6=icanhazip-v4 7=126 8=httpbin 9=ip-api 10=ipip 11=qq-ip2city 12=icanhazip 13=api64-ipify）；地区查询 `D加载默认数据文件()`（`data/ip2region_v4.xdb`）或 `D加载数据文件(路径)`，再 `D查询(ip)`。依赖 `lionsoul2014/ip2region`。
- **Proxy**：`池 := Proxy.D新建()` 有状态实例；配置链接/容量/验证/存活等后 `D启动`；`GetProxy(waitTime)`（`-1` 无限等，`>=0` 超时秒）时剔除过期并按缺口异步补充；支持 http/https/socks5；`D取统计` 看计数；`D关闭` 收尾。
- **Thread**：`Thread.D新建(最大线程数)` → `D启动`；`D投递`（无返回，爬虫常用）/ `D投递任务`（可 `D等待` 结果）；`D暂停`/`D继续`；`D设置最大线程数` 动态扩缩；`D停止` 排空 / `D强制停止` 丢弃排队。
- **HTTP**：日常 `HTTP.Send("GET", url)` 或 `HTTP.Send("GET", url, D发送选项{...})`，无需 `D新建`。命名空间可 `D设置代理池` / `D设置默认超时|重试` / `D新建会话`。独立客户端用 `HTTP.D新建()`。选项字段：`Data`/`Cookie`/`Headers`/`NoRedirect`/`UseProxy`/`Proxy`/`MaxRetries`/`Timeout`/`CloseMergeCookie`。`UseProxy=false` 直连；`true`+`Proxy`非空用手动；`true`+空则粘性/池。`CloseMergeCookie` 默认 false（合并）。`MaxRetries`/`Timeout` `<=0` 用默认（3 / 10s）。`Data`：`nil`/空串无 body，`[]byte`/`string` 原样，`map[string]string` 表单。响应：`D取状态码`/`D取正文`/`D取字节`/`D取重定向地址`/`D取返回Cookie`/`D取Cookie`/`D取返回协议头`/`D取协议头值`/`D取使用代理`/`D取耗时`。直连复用 KeepAlive；代理禁用 KeepAlive。
- **JS**：`引擎, _ := JS.D新建("goja")`（或 `"v8"`）；`D设置变量` / `D执行` / `D执行文件` / `D取文本|整数|小数|逻辑` / `D关闭`。单实例内部加锁可并发调用（串行）；高并发请每协程一个引擎。默认 goja（免 CGO，Win/Linux 均可）。V8 依赖 `rogchap.com/v8go`，**仅 Linux/macOS + CGO**；Windows 上请用 goja。注意：当前 `go.mod` 硬依赖 v8go，即使用 goja 也会下载该依赖。
- **Sync**：`Sync.D新建互斥锁`（`D加锁`/`D解锁`/`D尝试加锁`）；`D新建读写锁`（读/写加解锁与尝试）；`D新建原子计数` / `D新建原子计数初值`（`D加`/`D减`/`D取`/`D设置`/`D比较交换`）。
- **Codec**：`D网址编码`/`D网址解码`（query，空格 `%20`）；`DBase64编码`/`DBase64解码`；`DBase64网址编码`/`DBase64网址解码`（RawURL）；`D十六进制编码`/`D十六进制解码`（小写）；`D是否UTF8`；`D到GBK`/`D从GBK`（返回 `error`）。解码失败返回空串。
- 命名空间对象一般不保存可变全局状态（`IP.D加载数据文件` / `D加载默认数据文件`、`Proxy`/`Thread`/`HTTP`/`JS`/`Sync` 实例除外）。
- **打印输出**：`darren.Info` / `darren.Error` / `darren.Infof`（zap 彩色，自动初始化，带行号，时间 `yyyy-mm-dd hh:mm:ss`）。不要用 `fmt.Print`。

## 目录与维护

```
darren.go                         类型别名 + 命名空间变量
log_api.go                        darren.Info / Error 等
doc.go                            包说明
internal/*/ns.go                  命名空间类型与导出方法
internal/text                     文本实现与测试
internal/dtime                    时间实现与测试
internal/dfile                    文件操作 / 回收站实现与测试
internal/dlog                     zap 实现（供根包日志转发）
internal/dcrypto                  加解密 / 摘要实现与测试
internal/djson                    JSON 取值 / 设置实现与测试
internal/dqr                      二维码生成 / 解析实现与测试
internal/dip                      IP 工具 / ip2region 实现与测试
internal/dproxy                   代理池实现与测试
internal/dthread                  线程池实现与测试
internal/dhttp                    HTTP 客户端 / 会话实现与测试
internal/djs                      JS 引擎（goja / 可选 V8）
internal/dsync                    互斥锁 / 读写锁 / 原子计数
internal/dcodec                   URL / Base64 / Hex / UTF-8 / GBK
```

`examples/` 仅本地联调、已 `.gitignore`，不入仓。Go 的 import 路径只能用 ASCII，所以 `internal` 目录用英文；对外调用是 `darren.D文本` / `darren.AES` / `darren.Info` 等。

新增能力固定步骤：

1. 在 `internal/` 下用英文目录实现并写测试；在 `ns.go` 定义 `Xxx类` 与方法。
2. 根包 `darren.go` 增加 `type Xxx类 = internal.Xxx类` 与 `var Xxx Xxx类`（中文能力加 `D` 前缀，英文不加）。
3. 更新 `doc.go` 与本 README。

密钥长度、IV、填充必须写进方法名（例如 `AES.EncryptCBC(data, 密钥hex, IVhex)`），不要一个 `Encrypt()` 暗猜模式。

Windows / Linux 专有 API 用 `xxx_windows.go`、`xxx_linux.go` 拆分；文本、时间、加解密保持纯标准库，两平台同一套代码。

## 在其他项目中引用

```
go get github.com/Darren5211314/go_common@v0.1.0
```

```go
import "github.com/Darren5211314/go_common"

func main() {
	darren.Info("ok")
}
```

本地联调可在业务 `go.mod` 加：

```
replace github.com/Darren5211314/go_common => F:/GoProject/go_commons
```

若仓库为 Private：`go env -w GOPRIVATE=github.com/Darren5211314`。

删除或改名已有方法会破坏兼容；只能新增方法，旧方法保留。
