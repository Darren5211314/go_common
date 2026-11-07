# darren

面向 Windows exe 与 Ubuntu 等 Linux 项目的个人 Go 工具库。调用方只导入根包，通过命名空间使用方法。

## 环境

- Go 1.26（本机 SDK：`D:\DevAPP\Go\SDK\1.26.10`）
- Module：`github.com/Darren5211314/go_common`（根包名仍为 `darren`，调用 `darren.Info` / `darren.IP` 等）。
- 即使业务只用 goja，`go.mod` 仍会拉取 `rogchap.com/v8go`（含预编译库，体积大）。Windows 上勿启用 V8；后续可再拆成 opt-in build tag / 子模块。

## 调用示例

```go
package main

import "github.com/Darren5211314/go_common"

func main() {
	左边 := darren.Text.LeftOf("abc_123", "_")
	右边 := darren.Text.RightOf("abc_123", "_")
	前两字 := darren.Text.Left("你好世界", 2)
	位置 := darren.Text.Find("你好世界", "世界") // 从 1 起算，结果为 3
	验证码 := darren.Text.Random(6, false, false, false) // 默认小写+数字
	戳 := darren.Time.Now(false) // false=13位毫秒，true=10位秒
	文本 := darren.Time.FormatUnix(戳, "yyyy-mm-dd hh:mm:ss")
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

- **导出**：标识符用英文 PascalCase（如 `LeftOf`、`Create`）。注释可用中文；参数名可用中文（包内局部）。
- 命名空间：`Text`、`Time`、`File`、`AES`、`DES`、`DES3`、`SM4`、`Hash`（含 `SM3`）、`HMAC`、`CRC`、`RSA`、`JSON`、`QR`、`IP`、`Proxy`、`Thread`、`HTTP`、`JS`、`Sync`、`Codec`。勿用小写 `string`/`time`/`file`。日志用根包 `darren.Debug` / `Info` / `Warn` / `Error` / `DPanic` / `Panic` / `Fatal` 及对应 `*f`。
- 字符串按 UTF-8 字符处理，不按字节切。
- `Find`、`Mid` 的位置从 **1** 起算；找不到返回 `0` / 空串。
- `Random(位数, 小写字母, 大写字母, 数字)`：三者都为 `false` 时默认小写字母 + 数字。
- `Text` 找不到内容或空输入时返回空字符串，不返回 `error`。加解密、HMAC、RSA 返回 `(结果, error)`；Hash/CRC 只返回值。
- `Now(取十位)`：`true`=秒（约 10 位），`false`=毫秒（约 13 位）。
- 时间输出格式只用通配符：`yyyy-mm-dd hh:mm:ss`。
- `FormatUnix` 自动识别 10/13 位；`ParseUnix(文本, 取十位)` / `Format(文本, 目标格式)` **自动识别输入**，不必传原格式；空文本格式化=现行时间。
- `Sleep(毫秒)`：阻塞当前协程；`<=0` 立即返回。
- **File**：`Size(路径, 单位)`（B/KB/MB/GB/TB，1024 进制）；`Delete` / `Trash`（Windows + Linux Trash）；`Move` / `Copy` / `Rename`；`Ext`（含点）；`Exists` / `DirExists`；`Mkdir(路径)`；`ExeDir()`；`Name(路径, 带扩展名)`；`SaveText(路径, 分隔符, 覆盖, ...内容)`（空分隔符默认 `----`，自动建父目录；`覆盖=false` 追加，`true` 覆盖）。
- **加解密**：模式写进方法名（`EncryptCBC` / `EncryptECB` / `EncryptGCM`）；密钥、IV/nonce、data、摘要均为十六进制字符串；CBC/ECB 用 PKCS7；GCM 返回「data||tag」hex，nonce 单独传。国密：`Hash.SM3`、`SM4`（128-bit 密钥，CBC/ECB/GCM）；依赖 `github.com/emmansun/gmsm`。
- **JSON**：路径用点号，数组下标为数字段（`data.user.name`、`data.tags.0`）。取值找不到返回空/`0`/`false`；`SetString`/`SetInt`/… 返回新字符串并自动创建中间节点；`Compact` / `Pretty` 非法返回空串。
- **QR**：`Generate` 返回 PNG Base64；`GenerateFile` 写 PNG；`ParseFile` / `ParseBase64` 读回文本。边长 `<=0` 默认 256；纠错固定 Medium。依赖 `skip2/go-qrcode`、`makiuchi-d/gozxing`。
- **IP**：`Random` / `Valid` / `Version` / `IsPrivate` / `IsLoopback` / `Local`（网卡局域网）/ `Public(接口) (string, error)`（按编号单源；1=3322 2=gdt 3=qq-checktime 4=pconline 5=ipinfo 6=icanhazip-v4 7=126 8=httpbin 9=ip-api 10=ipip 11=qq-ip2city 12=icanhazip 13=api64-ipify）；地区查询 `LoadDefaultDB()`（`data/ip2region_v4.xdb`）或 `LoadDB(路径)`，再 `Query(ip)`。依赖 `lionsoul2014/ip2region`。
- **Proxy**：`池 := Proxy.Create()`；`SetProxy` / `SetPoolNum` / `SetVerify` / `SetVerifyTimeout` / `SetTTL` / `SetVerifyWorkers` / `SetGetInterval` / `SetVerifyURL` / `SetProxyType` / `SetLog`；`Start` 后 `GetProxy(waitTime)`（`-1` 无限等，`>=0` 超时秒）；`Stats()`（`LinkCount`/`PoolCount`/`GetCount`/`VerifyOK`/`VerifyFail`/`Verifying`/`Expired`）；`Close`。默认打运行日志，`SetLog(false)` 关闭。
- **Thread**：`Thread.Create(最大线程数)` → `Start`；`Post`（无返回，爬虫常用）/ `PostTask`（可 `Wait` 结果）；`Pause`/`Resume`；`SetMaxWorkers` 动态扩缩；`Stop` 排空 / `ForceStop` 丢弃排队。
- **HTTP**：`HTTP.Send` / `HTTP.NewSession` / `HTTP.SetProxyPool`（共用默认客户端）。`SendOptions` 字段：`Data`/`Cookie`/`Headers`/`NoRedirect`/`UseProxy`/`Proxy`/`MaxRetries`/`Timeout`/`CloseMergeCookie`。`MaxRetries`/`Timeout` 未填或 `<=0` 用内置默认（3 / 10s）。`UseProxy=false` 直连；`true`+`Proxy`非空用手动；`true`+空则粘性/池。跟随重定向时跨跳合并 Cookie（抗跨域剥头）；空值 Cookie（如 `P_ck_ctl=`）不合并、不写出。`Session.SetCookie(串)` 可注入整串。响应：`GetStatusCode`/`GetBody`/`GetBytes`/`GetLocation`/`GetSetCookie`/`GetCookie`/`GetHeaders`/`GetHeader`/`GetUseProxy`/`GetElapsedMs`。直连复用 KeepAlive；代理禁用 KeepAlive。
- **JS**：`引擎, _ := JS.Create("goja")`（或 `"v8"`）；`SetVar` / `Eval` / `EvalFile` / `GetString`/`GetInt`/`GetFloat`/`GetBool` / `Close`。单实例内部加锁可并发调用（串行）；高并发请每协程一个引擎。默认 goja（免 CGO，Win/Linux 均可）。V8 依赖 `rogchap.com/v8go`，**仅 Linux/macOS + CGO**；Windows 上请用 goja。注意：当前 `go.mod` 硬依赖 v8go，即使用 goja 也会下载该依赖。
- **Sync**：`Sync.NewMutex`（`Lock`/`Unlock`/`TryLock`）；`NewRWMutex`（`RLock`/`RUnlock`/`Lock`/`Unlock`/`TryRLock`/`TryLock`）；`NewAtomicInt` / `NewAtomicIntValue`（`Add`/`Sub`/`Get`/`Set`/`CompareAndSwap`）。
- **Codec**：`URLEncode`/`URLDecode`（query，空格 `%20`）；`QueryParam(url, key)` 取查询参数；`Base64Encode`/`Base64Decode`；`Base64URLEncode`/`Base64URLDecode`（RawURL）；`HexEncode`/`HexDecode`（小写）；`IsUTF8`；`ToGBK`/`FromGBK`（返回 `error`）。解码/取值失败返回空串。
- 命名空间对象一般不保存可变全局状态（`IP.LoadDB` / `LoadDefaultDB`、`Proxy`/`Thread`/`HTTP`/`JS`/`Sync` 实例除外）。
- **打印输出**：`Debug`/`Info`/`Warn`/`Error`/`DPanic`/`Panic`/`Fatal` 及 `Debugf`/`Infof`/…（zap 彩色，自动初始化，带行号，时间 `yyyy-mm-dd hh:mm:ss.SSS`）。不要用 `fmt.Print`。

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

`examples/` 仅本地联调、已 `.gitignore`，不入仓。Go 的 import 路径只能用 ASCII，所以 `internal` 目录用英文；对外调用是 `darren.Text` / `darren.AES` / `darren.Info` 等。

新增能力固定步骤：

1. 在 `internal/` 下用英文目录实现并写测试；在 `ns.go` 定义 `XxxNS` 与方法。
2. 根包 `darren.go` 增加未导出别名（如 `type httpNS = dhttp.HTTPNS`）与 `var Xxx httpNS`（命名空间变量用英文 PascalCase，如 `Text`/`HTTP`）。
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
