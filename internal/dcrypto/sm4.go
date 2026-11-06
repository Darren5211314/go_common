package dcrypto

import (
	"github.com/emmansun/gmsm/sm4"
)

func parseSM4Key(keyHex string) ([]byte, error) {
	key, err := decodeHex(keyHex)
	if err != nil {
		return nil, err
	}
	if len(key) != 16 {
		return nil, errKeyLen
	}
	return key, nil
}

// SM4EncryptCBC SM4-CBC 加密，返回密文十六进制。
func SM4EncryptCBC(data, keyHex, IVhex string) (string, error) {
	key, err := parseSM4Key(keyHex)
	if err != nil {
		return "", err
	}
	IV, err := decodeHex(IVhex)
	if err != nil {
		return "", err
	}
	block, err := sm4.NewCipher(key)
	if err != nil {
		return "", err
	}
	out, err := encryptCBC(block, []byte(data), IV)
	if err != nil {
		return "", err
	}
	return encodeHex(out), nil
}

// SM4DecryptCBC SM4-CBC 解密，密文为十六进制。
func SM4DecryptCBC(data, keyHex, IVhex string) (string, error) {
	key, err := parseSM4Key(keyHex)
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
	block, err := sm4.NewCipher(key)
	if err != nil {
		return "", err
	}
	out, err := decryptCBC(block, raw, IV)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// SM4EncryptECB SM4-ECB 加密，返回密文十六进制。
func SM4EncryptECB(data, keyHex string) (string, error) {
	key, err := parseSM4Key(keyHex)
	if err != nil {
		return "", err
	}
	block, err := sm4.NewCipher(key)
	if err != nil {
		return "", err
	}
	out, err := encryptECB(block, []byte(data))
	if err != nil {
		return "", err
	}
	return encodeHex(out), nil
}

// SM4DecryptECB SM4-ECB 解密，密文为十六进制。
func SM4DecryptECB(data, keyHex string) (string, error) {
	key, err := parseSM4Key(keyHex)
	if err != nil {
		return "", err
	}
	raw, err := decodeHex(data)
	if err != nil {
		return "", err
	}
	block, err := sm4.NewCipher(key)
	if err != nil {
		return "", err
	}
	out, err := decryptECB(block, raw)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// SM4EncryptGCM SM4-GCM 加密；nonce 单独传入；返回「data||tag」十六进制。
func SM4EncryptGCM(data, keyHex, NonceHex string) (string, error) {
	key, err := parseSM4Key(keyHex)
	if err != nil {
		return "", err
	}
	nonce, err := decodeHex(NonceHex)
	if err != nil {
		return "", err
	}
	block, err := sm4.NewCipher(key)
	if err != nil {
		return "", err
	}
	aead, err := newGCM(block, len(nonce))
	if err != nil {
		return "", err
	}
	sealed := aead.Seal(nil, nonce, []byte(data), nil)
	return encodeHex(sealed), nil
}

// SM4DecryptGCM SM4-GCM 解密；data 为「data||tag」。
func SM4DecryptGCM(data, keyHex, NonceHex string) (string, error) {
	key, err := parseSM4Key(keyHex)
	if err != nil {
		return "", err
	}
	nonce, err := decodeHex(NonceHex)
	if err != nil {
		return "", err
	}
	sealed, err := decodeHex(data)
	if err != nil {
		return "", err
	}
	block, err := sm4.NewCipher(key)
	if err != nil {
		return "", err
	}
	aead, err := newGCM(block, len(nonce))
	if err != nil {
		return "", err
	}
	out, err := aead.Open(nil, nonce, sealed, nil)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
