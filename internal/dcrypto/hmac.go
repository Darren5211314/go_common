package dcrypto

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

func hmacHex(data, key string, newFn func() hash.Hash) (string, error) {
	mac := hmac.New(newFn, []byte(key))
	if _, err := mac.Write([]byte(data)); err != nil {
		return "", err
	}
	return encodeHex(mac.Sum(nil)), nil
}

// HMACSHA1 计算 HMAC-SHA1，返回十六进制。
func HMACSHA1(data, key string) (string, error) {
	return hmacHex(data, key, sha1.New)
}

// HMACSHA256 计算 HMAC-SHA256，返回十六进制。
func HMACSHA256(data, key string) (string, error) {
	return hmacHex(data, key, sha256.New)
}

// HMACSHA512 计算 HMAC-SHA512，返回十六进制。
func HMACSHA512(data, key string) (string, error) {
	return hmacHex(data, key, sha512.New)
}
