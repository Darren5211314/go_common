package dcrypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"

	"github.com/emmansun/gmsm/sm3"
)

func 摘要MD5(数据 string) string {
	和 := md5.Sum([]byte(数据))
	return 编码十六进制(和[:])
}

func 摘要SHA1(数据 string) string {
	和 := sha1.Sum([]byte(数据))
	return 编码十六进制(和[:])
}

func 摘要SHA256(数据 string) string {
	和 := sha256.Sum256([]byte(数据))
	return 编码十六进制(和[:])
}

func 摘要SHA512(数据 string) string {
	和 := sha512.Sum512([]byte(数据))
	return 编码十六进制(和[:])
}

func 摘要SM3(数据 string) string {
	和 := sm3.Sum([]byte(数据))
	return 编码十六进制(和[:])
}
