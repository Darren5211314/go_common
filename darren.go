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

// 命名空间类型：根包未导出，避免补全出现 AES / AESNS 成对；调用方只用变量入口。

type textNS = text.TextNS
type timeNS = dtime.TimeNS
type fileNS = dfile.FileNS

type aesNS = dcrypto.AESNS
type desNS = dcrypto.DESNS
type des3NS = dcrypto.DES3NS
type hashNS = dcrypto.HashNS
type hmacNS = dcrypto.HMACNS
type crcNS = dcrypto.CRCNS
type rsaNS = dcrypto.RSANS
type sm4NS = dcrypto.SM4NS

type jsonNS = djson.JSONNS
type qrNS = dqr.QRNS
type ipNS = dip.IPNS

type proxyNS = dproxy.ProxyNS
type threadNS = dthread.ThreadNS
type httpNS = dhttp.HTTPNS
type jsNS = djs.JSNS
type syncNS = dsync.SyncNS
type codecNS = dcodec.CodecNS

// 有状态类型别名（调用方需要写类型名）。

type ProxyPool = dproxy.ProxyPool
type ProxyStats = dproxy.ProxyStats

type ThreadPool = dthread.ThreadPool
type ThreadStats = dthread.ThreadStats
type Task = dthread.Task

type Session = dhttp.Session
type Response = dhttp.Response
type SendOptions = dhttp.SendOptions

type Engine = djs.Engine

type Mutex = dsync.Mutex
type RWMutex = dsync.RWMutex
type AtomicInt = dsync.AtomicInt

// 命名空间入口。

var (
	Text textNS
	Time timeNS
	File fileNS

	AES    aesNS
	DES    desNS
	DES3   des3NS
	Hash   hashNS
	HMAC   hmacNS
	CRC    crcNS
	RSA    rsaNS
	SM4    sm4NS
	JSON   jsonNS
	QR     qrNS
	IP     ipNS
	Proxy  proxyNS
	Thread threadNS
	HTTP   httpNS
	JS     jsNS
	Sync   syncNS
	Codec  codecNS
)
