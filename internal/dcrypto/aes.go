package dcrypto

import (
	"crypto/aes"
	"crypto/cipher"
)

func parseAESKey(keyHex string) ([]byte, error) {
	key, err := decodeHex(keyHex)
	if err != nil {
		return nil, err
	}
	keyLen := len(key)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return nil, errKeyLen
	}
	return key, nil
}

// AESEncryptCBC AES-CBC 加密，返回密文十六进制。
func AESEncryptCBC(data, keyHex, IVhex string) (string, error) {
	key, err := parseAESKey(keyHex)
	if err != nil {
		return "", err
	}
	IV, err := decodeHex(IVhex)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	out, err := encryptCBC(block, []byte(data), IV)
	if err != nil {
		return "", err
	}
	return encodeHex(out), nil
}

// AESDecryptCBC AES-CBC 解密，密文为十六进制。
func AESDecryptCBC(data, keyHex, IVhex string) (string, error) {
	key, err := parseAESKey(keyHex)
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
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	out, err := decryptCBC(block, raw, IV)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// AESEncryptECB AES-ECB 加密，返回密文十六进制。
func AESEncryptECB(data, keyHex string) (string, error) {
	key, err := parseAESKey(keyHex)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	out, err := encryptECB(block, []byte(data))
	if err != nil {
		return "", err
	}
	return encodeHex(out), nil
}

// AESDecryptECB AES-ECB 解密，密文为十六进制。
func AESDecryptECB(data, keyHex string) (string, error) {
	key, err := parseAESKey(keyHex)
	if err != nil {
		return "", err
	}
	raw, err := decodeHex(data)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	out, err := decryptECB(block, raw)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// AESEncryptGCM AES-GCM 加密；nonce 单独传入；返回「data||tag」十六进制。
func AESEncryptGCM(data, keyHex, NonceHex string) (string, error) {
	key, err := parseAESKey(keyHex)
	if err != nil {
		return "", err
	}
	nonce, err := decodeHex(NonceHex)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
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

// AESDecryptGCM AES-GCM 解密；data 为「data||tag」。
func AESDecryptGCM(data, keyHex, NonceHex string) (string, error) {
	key, err := parseAESKey(keyHex)
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
	block, err := aes.NewCipher(key)
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

func newGCM(block cipher.Block, nonceLen int) (cipher.AEAD, error) {
	if nonceLen == 12 {
		return cipher.NewGCM(block)
	}
	if nonceLen < 1 {
		return nil, errIVLen
	}
	return cipher.NewGCMWithNonceSize(block, nonceLen)
}
