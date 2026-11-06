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

var errRSAKey = errors.New("RSA密钥无效")

// RSAGenerateKeyPair 生成 PEM 格式公钥与私钥；位数小于等于 0 时使用 2048。
func RSAGenerateKeyPair(bits int) (publicKeyPEM, privateKeyPEM string, err error) {
	if bits <= 0 {
		bits = 2048
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}
	privBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubBlock := &pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes}
	return string(pem.EncodeToMemory(pubBlock)), string(pem.EncodeToMemory(privBlock)), nil
}

func parsePublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(strings.TrimSpace(publicKeyPEM)))
	if block == nil {
		return nil, errRSAKey
	}
	pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		// 兼容 PKCS1 公钥
		publicKey, err2 := x509.ParsePKCS1PublicKey(block.Bytes)
		if err2 != nil {
			return nil, err
		}
		return publicKey, nil
	}
	publicKey, ok := pubAny.(*rsa.PublicKey)
	if !ok {
		return nil, errRSAKey
	}
	return publicKey, nil
}

func parsePrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(strings.TrimSpace(privateKeyPEM)))
	if block == nil {
		return nil, errRSAKey
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		key, err2 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err2 != nil {
			return nil, err
		}
		rsaPrivateKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errRSAKey
		}
		return rsaPrivateKey, nil
	}
	return privateKey, nil
}

// RSAEncrypt 公钥加密（PKCS#1 v1.5），返回密文十六进制。
func RSAEncrypt(data, publicKeyPEM string) (string, error) {
	publicKey, err := parsePublicKey(publicKeyPEM)
	if err != nil {
		return "", err
	}
	out, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(data))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(out), nil
}

// RSADecrypt 私钥解密（PKCS#1 v1.5），密文为十六进制。
func RSADecrypt(data, privateKeyPEM string) (string, error) {
	privateKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return "", err
	}
	raw, err := hex.DecodeString(strings.TrimSpace(data))
	if err != nil {
		return "", err
	}
	out, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, raw)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// RSASign 私钥签名（SHA256 + PKCS1v15），返回签名十六进制。
func RSASign(data, privateKeyPEM string) (string, error) {
	privateKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(data))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, sum[:])
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(signature), nil
}

// RSAVerify 公钥验签（SHA256 + PKCS1v15）。
func RSAVerify(data, signatureHex, publicKeyPEM string) error {
	publicKey, err := parsePublicKey(publicKeyPEM)
	if err != nil {
		return err
	}
	signature, err := hex.DecodeString(strings.TrimSpace(signatureHex))
	if err != nil {
		return err
	}
	sum := sha256.Sum256([]byte(data))
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, sum[:], signature)
}
