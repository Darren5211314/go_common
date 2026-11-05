package darren

import (
	"github.com/Darren5211314/go_common/internal/dcodec"
	"github.com/Darren5211314/go_common/internal/dcrypto"
	"github.com/Darren5211314/go_common/internal/dfile"
	"github.com/Darren5211314/go_common/internal/dhttp"
	"github.com/Darren5211314/go_common/internal/dip"
	"github.com/Darren5211314/go_common/internal/djs"
	"github.com/Darren5211314/go_common/internal/djson"
	"github.com/Darren5211314/go_common/internal/dproxy"
	"github.com/Darren5211314/go_common/internal/dqr"
	"github.com/Darren5211314/go_common/internal/dsync"
	"github.com/Darren5211314/go_common/internal/dthread"
	"github.com/Darren5211314/go_common/internal/dtime"
	"github.com/Darren5211314/go_common/internal/text"
)

// 命名空间类型别名：实现与方法在 internal，根包只导出入口。

type D文本类 = text.D文本类
type D时间类 = dtime.D时间类
type D文件类 = dfile.D文件类

type AES类 = dcrypto.AES类
type DES类 = dcrypto.DES类
type DES3类 = dcrypto.DES3类
type Hash类 = dcrypto.Hash类
type HMAC类 = dcrypto.HMAC类
type CRC类 = dcrypto.CRC类
type RSA类 = dcrypto.RSA类
type SM4类 = dcrypto.SM4类

type JSON类 = djson.JSON类
type QR类 = dqr.QR类
type IP类 = dip.IP类

type Proxy类 = dproxy.Proxy类
type Thread类 = dthread.Thread类
type HTTP类 = dhttp.HTTP类
type JS类 = djs.JS类
type Sync类 = dsync.Sync类
type Codec类 = dcodec.Codec类

// 有状态类型别名。

type D代理池 = dproxy.D代理池
type D代理统计 = dproxy.D代理统计

type D线程池 = dthread.D线程池
type D线程池统计 = dthread.D线程池统计
type D任务 = dthread.D任务

type D客户端 = dhttp.D客户端
type D会话 = dhttp.D会话
type D响应 = dhttp.D响应

type D引擎 = djs.D引擎

type D互斥锁 = dsync.D互斥锁
type D读写锁 = dsync.D读写锁
type D原子计数 = dsync.D原子计数

// 命名空间入口。

var (
	D文本 D文本类
	D时间 D时间类
	D文件 D文件类

	AES    AES类
	DES    DES类
	DES3   DES3类
	Hash   Hash类
	HMAC   HMAC类
	CRC    CRC类
	RSA    RSA类
	SM4    SM4类
	JSON   JSON类
	QR     QR类
	IP     IP类
	Proxy  Proxy类
	Thread Thread类
	HTTP   HTTP类
	JS     JS类
	Sync   Sync类
	Codec  Codec类
)
