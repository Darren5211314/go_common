package dcrypto

import (
	"crypto/des"
)

func 解析DES密钥(密钥hex string) ([]byte, error) {
	密钥, 错误 := 解码十六进制(密钥hex)
	if 错误 != nil {
		return nil, 错误
	}
	if len(密钥) != 8 {
		return nil, 错误密钥长度
	}
	return 密钥, nil
}

func 解析DES3密钥(密钥hex string) ([]byte, error) {
	密钥, 错误 := 解码十六进制(密钥hex)
	if 错误 != nil {
		return nil, 错误
	}
	切换 := len(密钥)
	switch 切换 {
	case 24:
		return 密钥, nil
	case 16:
		// EDE2：K1|K2 → K1|K2|K1
		扩展 := make([]byte, 24)
		copy(扩展, 密钥)
		copy(扩展[16:], 密钥[:8])
		return 扩展, nil
	default:
		return nil, 错误密钥长度
	}
}

// DESEncryptCBC DES-CBC 加密。
func DESEncryptCBC(明文, 密钥hex, IVhex string) (string, error) {
	密钥, 错误 := 解析DES密钥(密钥hex)
	if 错误 != nil {
		return "", 错误
	}
	IV, 错误 := 解码十六进制(IVhex)
	if 错误 != nil {
		return "", 错误
	}
	块, 错误 := des.NewCipher(密钥)
	if 错误 != nil {
		return "", 错误
	}
	密文, 错误 := 加密CBC(块, []byte(明文), IV)
	if 错误 != nil {
		return "", 错误
	}
	return 编码十六进制(密文), nil
}

// DESDecryptCBC DES-CBC 解密。
func DESDecryptCBC(密文hex, 密钥hex, IVhex string) (string, error) {
	密钥, 错误 := 解析DES密钥(密钥hex)
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
	块, 错误 := des.NewCipher(密钥)
	if 错误 != nil {
		return "", 错误
	}
	明文, 错误 := 解密CBC(块, 密文, IV)
	if 错误 != nil {
		return "", 错误
	}
	return string(明文), nil
}

// DESEncryptECB DES-ECB 加密。
func DESEncryptECB(明文, 密钥hex string) (string, error) {
	密钥, 错误 := 解析DES密钥(密钥hex)
	if 错误 != nil {
		return "", 错误
	}
	块, 错误 := des.NewCipher(密钥)
	if 错误 != nil {
		return "", 错误
	}
	密文, 错误 := 加密ECB(块, []byte(明文))
	if 错误 != nil {
		return "", 错误
	}
	return 编码十六进制(密文), nil
}

// DESDecryptECB DES-ECB 解密。
func DESDecryptECB(密文hex, 密钥hex string) (string, error) {
	密钥, 错误 := 解析DES密钥(密钥hex)
	if 错误 != nil {
		return "", 错误
	}
	密文, 错误 := 解码十六进制(密文hex)
	if 错误 != nil {
		return "", 错误
	}
	块, 错误 := des.NewCipher(密钥)
	if 错误 != nil {
		return "", 错误
	}
	明文, 错误 := 解密ECB(块, 密文)
	if 错误 != nil {
		return "", 错误
	}
	return string(明文), nil
}

// DES3EncryptCBC 3DES-CBC 加密。
func DES3EncryptCBC(明文, 密钥hex, IVhex string) (string, error) {
	密钥, 错误 := 解析DES3密钥(密钥hex)
	if 错误 != nil {
		return "", 错误
	}
	IV, 错误 := 解码十六进制(IVhex)
	if 错误 != nil {
		return "", 错误
	}
	块, 错误 := des.NewTripleDESCipher(密钥)
	if 错误 != nil {
		return "", 错误
	}
	密文, 错误 := 加密CBC(块, []byte(明文), IV)
	if 错误 != nil {
		return "", 错误
	}
	return 编码十六进制(密文), nil
}

// DES3DecryptCBC 3DES-CBC 解密。
func DES3DecryptCBC(密文hex, 密钥hex, IVhex string) (string, error) {
	密钥, 错误 := 解析DES3密钥(密钥hex)
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
	块, 错误 := des.NewTripleDESCipher(密钥)
	if 错误 != nil {
		return "", 错误
	}
	明文, 错误 := 解密CBC(块, 密文, IV)
	if 错误 != nil {
		return "", 错误
	}
	return string(明文), nil
}

// DES3EncryptECB 3DES-ECB 加密。
func DES3EncryptECB(明文, 密钥hex string) (string, error) {
	密钥, 错误 := 解析DES3密钥(密钥hex)
	if 错误 != nil {
		return "", 错误
	}
	块, 错误 := des.NewTripleDESCipher(密钥)
	if 错误 != nil {
		return "", 错误
	}
	密文, 错误 := 加密ECB(块, []byte(明文))
	if 错误 != nil {
		return "", 错误
	}
	return 编码十六进制(密文), nil
}

// DES3DecryptECB 3DES-ECB 解密。
func DES3DecryptECB(密文hex, 密钥hex string) (string, error) {
	密钥, 错误 := 解析DES3密钥(密钥hex)
	if 错误 != nil {
		return "", 错误
	}
	密文, 错误 := 解码十六进制(密文hex)
	if 错误 != nil {
		return "", 错误
	}
	块, 错误 := des.NewTripleDESCipher(密钥)
	if 错误 != nil {
		return "", 错误
	}
	明文, 错误 := 解密ECB(块, 密文)
	if 错误 != nil {
		return "", 错误
	}
	return string(明文), nil
}
