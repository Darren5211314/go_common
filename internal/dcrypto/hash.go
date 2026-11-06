package dcrypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"

	"github.com/emmansun/gmsm/sm3"
)

func MD5(data string) string {
	sum := md5.Sum([]byte(data))
	return encodeHex(sum[:])
}

func SHA1(data string) string {
	sum := sha1.Sum([]byte(data))
	return encodeHex(sum[:])
}

func SHA256(data string) string {
	sum := sha256.Sum256([]byte(data))
	return encodeHex(sum[:])
}

func SHA512(data string) string {
	sum := sha512.Sum512([]byte(data))
	return encodeHex(sum[:])
}

func SM3(data string) string {
	sum := sm3.Sum([]byte(data))
	return encodeHex(sum[:])
}
