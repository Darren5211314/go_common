package dcrypto

import (
	"crypto/aes"
	"crypto/cipher"
)

func 解析AES密钥(密钥hex string) ([]byte, error) {
	密钥, 错误 := 解码十六进制(密钥hex)
	if 错误 != nil {
		return nil, 错误
	}
	切换 := len(密钥)
	if 切换 != 16 && 切换 != 24 && 切换 != 32 {
		return nil, 错误密钥长度
	}
	return 密钥, nil
}

// AESEncryptCBC AES-CBC 加密，返回密文十六进制。
func AESEncryptCBC(明文, 密钥hex, IVhex string) (string, error) {
	密钥, 错误 := 解析AES密钥(密钥hex)
	if 错误 != nil {
		return "", 错误
	}
	IV, 错误 := 解码十六进制(IVhex)
	if 错误 != nil {
		return "", 错误
	}
	块, 错误 := aes.NewCipher(密钥)
	if 错误 != nil {
		return "", 错误
	}
	密文, 错误 := 加密CBC(块, []byte(明文), IV)
	if 错误 != nil {
		return "", 错误
	}
	return 编码十六进制(密文), nil
}

// AESDecryptCBC AES-CBC 解密，密文为十六进制。
func AESDecryptCBC(密文hex, 密钥hex, IVhex string) (string, error) {
	密钥, 错误 := 解析AES密钥(密钥hex)
	if 错误 != nil {
		return "", 错误
	}
	IV, 错误 := 解码十六进制(IVhex)
	if 错误 != nil {
		return "", 错误
	}
	密文, 错误 := 解码十六进制(密文hex)
	if 错误 != nil {
		return "", 错误
	}
	块, 错误 := aes.NewCipher(密钥)
	if 错误 != nil {
		return "", 错误
	}
	明文, 错误 := 解密CBC(块, 密文, IV)
	if 错误 != nil {
		return "", 错误
	}
	return string(明文), nil
}

// AESEncryptECB AES-ECB 加密，返回密文十六进制。
func AESEncryptECB(明文, 密钥hex string) (string, error) {
	密钥, 错误 := 解析AES密钥(密钥hex)
	if 错误 != nil {
		return "", 错误
	}
	块, 错误 := aes.NewCipher(密钥)
	if 错误 != nil {
		return "", 错误
	}
	密文, 错误 := 加密ECB(块, []byte(明文))
	if 错误 != nil {
		return "", 错误
	}
	return 编码十六进制(密文), nil
}

// AESDecryptECB AES-ECB 解密，密文为十六进制。
func AESDecryptECB(密文hex, 密钥hex string) (string, error) {
	密钥, 错误 := 解析AES密钥(密钥hex)
	if 错误 != nil {
		return "", 错误
	}
	密文, 错误 := 解码十六进制(密文hex)
	if 错误 != nil {
		return "", 错误
	}
	块, 错误 := aes.NewCipher(密钥)
	if 错误 != nil {
		return "", 错误
	}
	明文, 错误 := 解密ECB(块, 密文)
	if 错误 != nil {
		return "", 错误
	}
	return string(明文), nil
}

// AESEncryptGCM AES-GCM 加密；nonce 单独传入；返回「密文||tag」十六进制。
func AESEncryptGCM(明文, 密钥hex, NonceHex string) (string, error) {
	密钥, 错误 := 解析AES密钥(密钥hex)
	if 错误 != nil {
		return "", 错误
	}
	nonce, 错误 := 解码十六进制(NonceHex)
	if 错误 != nil {
		return "", 错误
	}
	块, 错误 := aes.NewCipher(密钥)
	if 错误 != nil {
		return "", 错误
	}
	aead, 错误 := 新建GCM(块, len(nonce))
	if 错误 != nil {
		return "", 错误
	}
	密封 := aead.Seal(nil, nonce, []byte(明文), nil)
	return 编码十六进制(密封), nil
}

// AESDecryptGCM AES-GCM 解密；密文hex 为「密文||tag」。
func AESDecryptGCM(密文hex, 密钥hex, NonceHex string) (string, error) {
	密钥, 错误 := 解析AES密钥(密钥hex)
	if 错误 != nil {
		return "", 错误
	}
	nonce, 错误 := 解码十六进制(NonceHex)
	if 错误 != nil {
		return "", 错误
	}
	密封, 错误 := 解码十六进制(密文hex)
	if 错误 != nil {
		return "", 错误
	}
	块, 错误 := aes.NewCipher(密钥)
	if 错误 != nil {
		return "", 错误
	}
	aead, 错误 := 新建GCM(块, len(nonce))
	if 错误 != nil {
		return "", 错误
	}
	明文, 错误 := aead.Open(nil, nonce, 密封, nil)
	if 错误 != nil {
		return "", 错误
	}
	return string(明文), nil
}

func 新建GCM(块 cipher.Block, nonce长度 int) (cipher.AEAD, error) {
	if nonce长度 == 12 {
		return cipher.NewGCM(块)
	}
	if nonce长度 < 1 {
		return nil, 错误IV长度
	}
	return cipher.NewGCMWithNonceSize(块, nonce长度)
}
