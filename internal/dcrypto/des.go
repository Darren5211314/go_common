package dcrypto

import (
	"crypto/des"
)

func parseDESKey(keyHex string) ([]byte, error) {
	key, err := decodeHex(keyHex)
	if err != nil {
		return nil, err
	}
	if len(key) != 8 {
		return nil, errKeyLen
	}
	return key, nil
}

func parseDES3Key(keyHex string) ([]byte, error) {
	key, err := decodeHex(keyHex)
	if err != nil {
		return nil, err
	}
	keyLen := len(key)
	switch keyLen {
	case 24:
		return key, nil
	case 16:
		// EDE2：K1|K2 → K1|K2|K1
		expanded := make([]byte, 24)
		copy(expanded, key)
		copy(expanded[16:], key[:8])
		return expanded, nil
	default:
		return nil, errKeyLen
	}
}

// DESEncryptCBC DES-CBC 加密。
func DESEncryptCBC(data, keyHex, IVhex string) (string, error) {
	key, err := parseDESKey(keyHex)
	if err != nil {
		return "", err
	}
	IV, err := decodeHex(IVhex)
	if err != nil {
		return "", err
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	out, err := encryptCBC(block, []byte(data), IV)
	if err != nil {
		return "", err
	}
	return encodeHex(out), nil
}

// DESDecryptCBC DES-CBC 解密。
func DESDecryptCBC(data, keyHex, IVhex string) (string, error) {
	key, err := parseDESKey(keyHex)
	if err != nil {
		return "", err
	}
	IV, err := decodeHex(IVhex)
	if err != nil {
		return "", err
	}
	raw, err := decodeHex(data)
	if err != nil {
		return "", err
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	out, err := decryptCBC(block, raw, IV)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// DESEncryptECB DES-ECB 加密。
func DESEncryptECB(data, keyHex string) (string, error) {
	key, err := parseDESKey(keyHex)
	if err != nil {
		return "", err
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	out, err := encryptECB(block, []byte(data))
	if err != nil {
		return "", err
	}
	return encodeHex(out), nil
}

// DESDecryptECB DES-ECB 解密。
func DESDecryptECB(data, keyHex string) (string, error) {
	key, err := parseDESKey(keyHex)
	if err != nil {
		return "", err
	}
	raw, err := decodeHex(data)
	if err != nil {
		return "", err
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	out, err := decryptECB(block, raw)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// DES3EncryptCBC 3DES-CBC 加密。
func DES3EncryptCBC(data, keyHex, IVhex string) (string, error) {
	key, err := parseDES3Key(keyHex)
	if err != nil {
		return "", err
	}
	IV, err := decodeHex(IVhex)
	if err != nil {
		return "", err
	}
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	out, err := encryptCBC(block, []byte(data), IV)
	if err != nil {
		return "", err
	}
	return encodeHex(out), nil
}

// DES3DecryptCBC 3DES-CBC 解密。
func DES3DecryptCBC(data, keyHex, IVhex string) (string, error) {
	key, err := parseDES3Key(keyHex)
	if err != nil {
		return "", err
	}
	IV, err := decodeHex(IVhex)
	if err != nil {
		return "", err
	}
	raw, err := decodeHex(data)
	if err != nil {
		return "", err
	}
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	out, err := decryptCBC(block, raw, IV)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// DES3EncryptECB 3DES-ECB 加密。
func DES3EncryptECB(data, keyHex string) (string, error) {
	key, err := parseDES3Key(keyHex)
	if err != nil {
		return "", err
	}
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	out, err := encryptECB(block, []byte(data))
	if err != nil {
		return "", err
	}
	return encodeHex(out), nil
}

// DES3DecryptECB 3DES-ECB 解密。
func DES3DecryptECB(data, keyHex string) (string, error) {
	key, err := parseDES3Key(keyHex)
	if err != nil {
		return "", err
	}
	raw, err := decodeHex(data)
	if err != nil {
		return "", err
	}
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	out, err := decryptECB(block, raw)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
