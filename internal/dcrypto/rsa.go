package dcrypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"strings"
)

var 错误RSA密钥 = errors.New("RSA密钥无效")

// RSAGenerateKeyPair 生成 PEM 格式公钥与私钥；位数小于等于 0 时使用 2048。
func RSAGenerateKeyPair(位数 int) (公钥PEM, 私钥PEM string, err error) {
	if 位数 <= 0 {
		位数 = 2048
	}
	私钥, 错误 := rsa.GenerateKey(rand.Reader, 位数)
	if 错误 != nil {
		return "", "", 错误
	}
	私钥块 := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(私钥)}
	公钥字节, 错误 := x509.MarshalPKIXPublicKey(&私钥.PublicKey)
	if 错误 != nil {
		return "", "", 错误
	}
	公钥块 := &pem.Block{Type: "PUBLIC KEY", Bytes: 公钥字节}
	return string(pem.EncodeToMemory(公钥块)), string(pem.EncodeToMemory(私钥块)), nil
}

func 解析公钥(公钥PEM string) (*rsa.PublicKey, error) {
	块, _ := pem.Decode([]byte(strings.TrimSpace(公钥PEM)))
	if 块 == nil {
		return nil, 错误RSA密钥
	}
	公钥任意, 错误 := x509.ParsePKIXPublicKey(块.Bytes)
	if 错误 != nil {
		// 兼容 PKCS1 公钥
		公钥, 错误2 := x509.ParsePKCS1PublicKey(块.Bytes)
		if 错误2 != nil {
			return nil, 错误
		}
		return 公钥, nil
	}
	公钥, 成功 := 公钥任意.(*rsa.PublicKey)
	if !成功 {
		return nil, 错误RSA密钥
	}
	return 公钥, nil
}

func 解析私钥(私钥PEM string) (*rsa.PrivateKey, error) {
	块, _ := pem.Decode([]byte(strings.TrimSpace(私钥PEM)))
	if 块 == nil {
		return nil, 错误RSA密钥
	}
	私钥, 错误 := x509.ParsePKCS1PrivateKey(块.Bytes)
	if 错误 != nil {
		密钥, 错误2 := x509.ParsePKCS8PrivateKey(块.Bytes)
		if 错误2 != nil {
			return nil, 错误
		}
		rsa私钥, 成功 := 密钥.(*rsa.PrivateKey)
		if !成功 {
			return nil, 错误RSA密钥
		}
		return rsa私钥, nil
	}
	return 私钥, nil
}

// RSAEncrypt 公钥加密（PKCS#1 v1.5），返回密文十六进制。
func RSAEncrypt(明文, 公钥PEM string) (string, error) {
	公钥, 错误 := 解析公钥(公钥PEM)
	if 错误 != nil {
		return "", 错误
	}
	密文, 错误 := rsa.EncryptPKCS1v15(rand.Reader, 公钥, []byte(明文))
	if 错误 != nil {
		return "", 错误
	}
	return hex.EncodeToString(密文), nil
}

// RSADecrypt 私钥解密（PKCS#1 v1.5），密文为十六进制。
func RSADecrypt(密文hex, 私钥PEM string) (string, error) {
	私钥, 错误 := 解析私钥(私钥PEM)
	if 错误 != nil {
		return "", 错误
	}
	密文, 错误 := hex.DecodeString(strings.TrimSpace(密文hex))
	if 错误 != nil {
		return "", 错误
	}
	明文, 错误 := rsa.DecryptPKCS1v15(rand.Reader, 私钥, 密文)
	if 错误 != nil {
		return "", 错误
	}
	return string(明文), nil
}

// RSASign 私钥签名（SHA256 + PKCS1v15），返回签名十六进制。
func RSASign(数据, 私钥PEM string) (string, error) {
	私钥, 错误 := 解析私钥(私钥PEM)
	if 错误 != nil {
		return "", 错误
	}
	和 := sha256.Sum256([]byte(数据))
	签名, 错误 := rsa.SignPKCS1v15(rand.Reader, 私钥, crypto.SHA256, 和[:])
	if 错误 != nil {
		return "", 错误
	}
	return hex.EncodeToString(签名), nil
}

// RSAVerify 公钥验签（SHA256 + PKCS1v15）。
func RSAVerify(数据, 签名hex, 公钥PEM string) error {
	公钥, 错误 := 解析公钥(公钥PEM)
	if 错误 != nil {
		return 错误
	}
	签名, 错误 := hex.DecodeString(strings.TrimSpace(签名hex))
	if 错误 != nil {
		return 错误
	}
	和 := sha256.Sum256([]byte(数据))
	return rsa.VerifyPKCS1v15(公钥, crypto.SHA256, 和[:], 签名)
}
