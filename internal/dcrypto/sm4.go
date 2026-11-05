package dcrypto

import (
	"github.com/emmansun/gmsm/sm4"
)

func 解析SM4密钥(密钥hex string) ([]byte, error) {
	密钥, 错误 := 解码十六进制(密钥hex)
	if 错误 != nil {
		return nil, 错误
	}
	if len(密钥) != 16 {
		return nil, 错误密钥长度
	}
	return 密钥, nil
}

// SM4EncryptCBC SM4-CBC 加密，返回密文十六进制。
func SM4EncryptCBC(明文, 密钥hex, IVhex string) (string, error) {
	密钥, 错误 := 解析SM4密钥(密钥hex)
	if 错误 != nil {
		return "", 错误
	}
	IV, 错误 := 解码十六进制(IVhex)
	if 错误 != nil {
		return "", 错误
	}
	块, 错误 := sm4.NewCipher(密钥)
	if 错误 != nil {
		return "", 错误
	}
	密文, 错误 := 加密CBC(块, []byte(明文), IV)
	if 错误 != nil {
		return "", 错误
	}
	return 编码十六进制(密文), nil
}

// SM4DecryptCBC SM4-CBC 解密，密文为十六进制。
func SM4DecryptCBC(密文hex, 密钥hex, IVhex string) (string, error) {
	密钥, 错误 := 解析SM4密钥(密钥hex)
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
	块, 错误 := sm4.NewCipher(密钥)
	if 错误 != nil {
		return "", 错误
	}
	明文, 错误 := 解密CBC(块, 密文, IV)
	if 错误 != nil {
		return "", 错误
	}
	return string(明文), nil
}

// SM4EncryptECB SM4-ECB 加密，返回密文十六进制。
func SM4EncryptECB(明文, 密钥hex string) (string, error) {
	密钥, 错误 := 解析SM4密钥(密钥hex)
	if 错误 != nil {
		return "", 错误
	}
	块, 错误 := sm4.NewCipher(密钥)
	if 错误 != nil {
		return "", 错误
	}
	密文, 错误 := 加密ECB(块, []byte(明文))
	if 错误 != nil {
		return "", 错误
	}
	return 编码十六进制(密文), nil
}

// SM4DecryptECB SM4-ECB 解密，密文为十六进制。
func SM4DecryptECB(密文hex, 密钥hex string) (string, error) {
	密钥, 错误 := 解析SM4密钥(密钥hex)
	if 错误 != nil {
		return "", 错误
	}
	密文, 错误 := 解码十六进制(密文hex)
	if 错误 != nil {
		return "", 错误
	}
	块, 错误 := sm4.NewCipher(密钥)
	if 错误 != nil {
		return "", 错误
	}
	明文, 错误 := 解密ECB(块, 密文)
	if 错误 != nil {
		return "", 错误
	}
	return string(明文), nil
}

// SM4EncryptGCM SM4-GCM 加密；nonce 单独传入；返回「密文||tag」十六进制。
func SM4EncryptGCM(明文, 密钥hex, NonceHex string) (string, error) {
	密钥, 错误 := 解析SM4密钥(密钥hex)
	if 错误 != nil {
		return "", 错误
	}
	nonce, 错误 := 解码十六进制(NonceHex)
	if 错误 != nil {
		return "", 错误
	}
	块, 错误 := sm4.NewCipher(密钥)
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

// SM4DecryptGCM SM4-GCM 解密；密文hex 为「密文||tag」。
func SM4DecryptGCM(密文hex, 密钥hex, NonceHex string) (string, error) {
	密钥, 错误 := 解析SM4密钥(密钥hex)
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
	块, 错误 := sm4.NewCipher(密钥)
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
