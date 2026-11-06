// Package darren 是面向 Windows exe 与 Ubuntu 等 Linux 项目的个人工具库。
//
// 调用方只导入根包，通过命名空间对象使用方法，例如：
//
//	darren.Text.D取文本左边("abc_123", "_")
//	darren.Time.D取现行时间戳(false)
//
// # 命名约定
//
// Go 要求跨包导出的标识符必须以 Unicode 大写字母开头；纯中文属于“其他字母”，
// 无法导出。命名空间用英文入口（Text / Time / File / AES / HTTP 等）；
// 中文方法仍以大写 D 开头，例如 D取文本左边。注释仍使用中文。
// 不要使用小写 string、time、file 作为命名空间名（无法导出，且与内置类型/标准库冲突）。
//
// # 错误约定
//
// Text（TextNS）方法在找不到内容、空输入等情况下返回空字符串，不返回 error。
// 加解密、文件、系统调用等后续能力返回 (结果, error)。
//
// # 文本约定
//
// 字符串按 UTF-8 字符（rune）处理，不按字节切片，避免中文被截断。
// D寻找文本、D取文本中间 的位置从 1 起算；找不到或非法参数返回 0 / 空串。
// D取随机字符 在小写/大写/数字开关全为 false 时默认使用「小写字母 + 数字」。
//
// # 时间约定
//
// D取现行时间戳(取十位)：true 为秒（约 10 位），false 为毫秒（约 13 位，默认写法传 false）。
// D时间戳转时间 自动识别 10/13 位时间戳。
// 输出格式只使用通配符（如 yyyy-mm-dd hh:mm:ss）。
// D时间转时间戳 / D时间格式化 自动识别输入时间文本，不必传原格式；解析失败返回 0 / 空串。
// D时间格式化 时间文本为空时使用现行时间。D延时(毫秒) 阻塞当前协程，<=0 立即返回。
// 命名空间对象不保存可变全局状态。
//
// # 文件约定
//
// File：D取文件大小(路径, 单位) 单位为 B/KB/MB/GB/TB（1024 进制）。
// D删除文件 永久删除；D删除到回收站 支持 Windows 与 Linux XDG Trash。
// D移动文件 / D复制文件 / D重命名；目标父目录自动创建。
// D取扩展名 含点；D取文件名(路径, 带扩展名)；D文件是否存在 / D目录是否存在。
// D保存文本(路径, 分隔符, ...内容)：分隔符空则用 "----"；自动建父目录并覆盖写入 UTF-8。
//
// # 日志约定
//
//	darren.Info("消息", "键", 值)
//	darren.Error("失败", "错误", err)
//	darren.Infof("码=%d", 500)
//
// 首次调用自动初始化（zap 彩色控制台）。
// 带源码文件名与行号，时间为 yyyy-mm-dd hh:mm:ss.SSS。不要使用 fmt.Print。
//
// # 加解密约定
//
// 英文命名空间不加 D：AES、DES、DES3、SM4、Hash、HMAC、CRC、RSA、JSON、QR、IP、Proxy、Thread、HTTP、JS、Sync、Codec。
// JSON / QR / IP / Proxy / Thread / HTTP / JS / Sync / Codec 的方法仍为中文加 D。
// 模式写进方法名，例如 AES.EncryptCBC / EncryptECB / EncryptGCM。
// 密钥、IV/nonce、data、摘要、HMAC 均为十六进制字符串；明文为普通 string。
// CBC/ECB 使用 PKCS7；GCM 返回「data||tag」hex，nonce 单独参数传入。
// 国密：Hash.SM3；SM4 固定 128-bit 密钥，支持 CBC/ECB/GCM（依赖 emmansun/gmsm）。
// Hash/CRC 只返回值；加解密、HMAC、RSA、JSON.D设置*、QR、IP、Proxy、Thread 投递、HTTP.Send 等返回 (结果, error)。
//
//	摘要 := darren.Hash.MD5("abc")
//	本机 := darren.IP.D取本机IP()
//	代理池 := darren.Proxy.D新建()
//	线程池 := darren.Thread.D新建(32)
//	引擎, _ := darren.JS.D新建("goja")
//
// # JS 约定
//
// JS.D新建("goja"|"v8"|"" ) 得到有状态引擎；空=goja。
// V8 需 CGO，且 rogchap/v8go 仅提供 Linux/macOS 预编译库；Windows 请用 goja。
// 即使只用 goja，go.mod 仍拉取 v8go（体积大）；后续可改为 opt-in tag / 子模块。
// D设置变量 / D执行 / D执行文件 / D取文本|整数|小数|逻辑 / D关闭。
// 单实例内部互斥锁，多 goroutine 共用安全但串行；高并发请多实例。
//
// # Sync 约定
//
// Sync.D新建互斥锁 / D新建读写锁 / D新建原子计数 / D新建原子计数初值。
// 互斥：D加锁 D解锁 D尝试加锁；读写：D读/写加解锁与尝试；原子：D加 D减 D取 D设置 D比较交换。
//
// # Codec 约定
//
// Codec.D网址编码 / D网址解码（query 组件，空格为 %20）；失败解码返回空串。
// DBase64编码 / DBase64解码（StdEncoding）；DBase64网址编码 / DBase64网址解码（RawURLEncoding）。
// D十六进制编码 / D十六进制解码（小写）；D是否UTF8；D到GBK / D从GBK 返回 (结果, error)。
//
// # JSON 约定
//
// 路径用点号分隔，数组下标为数字段（data.tags.0）。
// 取值找不到或非法返回空字符串 / 0 / false；D设置* 写回新 JSON 并创建中间节点。
// D压缩 / D格式化 非法输入返回空串。
//
// # QR 约定
//
// D生成 返回 PNG Base64（无 data: 前缀）；D生成到文件 写 PNG。
// D解析文件 / D解析Base64 解析出文本；边长 <=0 默认 256；纠错固定 Medium。
//
// # IP 约定
//
// D取随机IP 生成公网风格 IPv4；D取本机IP 取首个非回环非 link-local 的 IPv4。
// D取外网IP(接口) (string, error)：按编号单源取公网 IP（约 3s）；未知编号/网络/解析失败均返回 error。
// 1=3322 2=gdt 3=qq-checktime 4=pconline 5=ipinfo 6=icanhazip-v4 7=126-locate
// 8=httpbin 9=ip-api-demo 10=ipip 11=qq-ip2city 12=icanhazip 13=api64-ipify。
// 地区查询：D加载默认数据文件()（data/ip2region_v4.xdb）或 D加载数据文件(路径)，再 D查询(ip)；xdb 不嵌入仓库。
// HTTP 直连复用 KeepAlive Transport；代理请求禁用 KeepAlive。
//
// # Proxy 约定
//
// Proxy.D新建 得到有状态池；D设置代理链接/容量/是否验证/验证超时/存活时间/验证线程数/获取间隔/验证地址/默认协议/输出日志；
// D启动 后 GetProxy(waitTime)（-1 无限等，>=0 超时秒）；D取统计；D关闭。默认输出拉取/验证/取出日志，D设置输出日志(false) 关闭。
// 过期仅在获取时剔除；仅在可用数低于容量时按获取间隔补充。支持 http/https/socks5。
//
// # Thread 约定
//
// Thread.D新建(最大线程数) 后 D启动；D投递(无返回) / D投递任务(可 D等待)；D暂停/D继续；
// D设置最大线程数 动态扩缩；D停止 排空队列；D强制停止 丢弃排队（不强杀执行中任务）。
// panic 计入失败数，不拖垮池子。D任务上下文 在强制停止时会被 cancel。
//
// # HTTP 约定
//
// HTTP.Send / HTTP.NewSession / HTTP.SetProxyPool（共用默认客户端）。
// SendOptions：Data、Cookie、Headers、NoRedirect、UseProxy、Proxy、MaxRetries、Timeout、CloseMergeCookie。
// UseProxy=false 直连；=true 且 Proxy 非空用手动；=true 且空则粘性/池/直连。
// Timeout/MaxRetries 未填或 <=0 用内置默认（10s / 3 次）；CloseMergeCookie 默认 false（合并）。
// Data：nil/空串无 body；[]byte/string 原样；map[string]string 为表单。
// Cookie 支持 a=1; b=2 与 a=1&b=2；值中勿含未编码的 &。
// 直连复用 KeepAlive；代理禁用 KeepAlive。仅超时/连接/代理/TLS 及 502/503/504 可重试。
// D响应：GetStatusCode、GetBody、GetBytes、GetLocation、GetSetCookie、GetCookie、GetHeaders、GetHeader、GetUseProxy、GetElapsedMs。
//
// # 扩展方式
//
// 新增能力步骤：internal 实现 + ns.go 定义 XxxNS/方法 → 测试 → darren.go 未导出别名与 var → 更新本文件。
// import 路径必须是 ASCII（例如 internal/text、internal/dcrypto）。
// 中文导出加 D 前缀；英文标识本身已大写开头则不加 D。
// Windows / Linux 专有 API 使用文件名 build tag（xxx_windows.go / xxx_linux.go）拆分。
package darren
