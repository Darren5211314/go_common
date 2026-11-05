package dcrypto

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

func hmacHex(数据, 密钥 string, 新建 func() hash.Hash) (string, error) {
	mac := hmac.New(新建, []byte(密钥))
	if _, 错误 := mac.Write([]byte(数据)); 错误 != nil {
		return "", 错误
	}
	return 编码十六进制(mac.Sum(nil)), nil
}

// HMACSHA1 计算 HMAC-SHA1，返回十六进制。
func HMACSHA1(数据, 密钥 string) (string, error) {
	return hmacHex(数据, 密钥, sha1.New)
}

// HMACSHA256 计算 HMAC-SHA256，返回十六进制。
func HMACSHA256(数据, 密钥 string) (string, error) {
	return hmacHex(数据, 密钥, sha256.New)
}

// HMACSHA512 计算 HMAC-SHA512，返回十六进制。
func HMACSHA512(数据, 密钥 string) (string, error) {
	return hmacHex(数据, 密钥, sha512.New)
}
