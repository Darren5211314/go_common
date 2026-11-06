package dcrypto

import (
	"encoding/hex"
	"errors"
	"strings"
)

var (
	errKeyLen = errors.New("密钥长度不正确")
	errIVLen  = errors.New("IV/Nonce长度不正确")
	errPadding    = errors.New("PKCS7填充无效")
	errCiphertext    = errors.New("密文无效")
	errData    = errors.New("数据无效")
)

func decodeHex(text string) ([]byte, error) {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, " ", "")
	data, err := hex.DecodeString(text)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func encodeHex(data []byte) string {
	return hex.EncodeToString(data)
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padLen := blockSize - len(data)%blockSize
	pad := make([]byte, padLen)
	for i := range pad {
		pad[i] = byte(padLen)
	}
	return append(data, pad...)
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	n := len(data)
	if n == 0 || n%blockSize != 0 {
		return nil, errPadding
	}
	padLen := int(data[n-1])
	if padLen == 0 || padLen > blockSize || padLen > n {
		return nil, errPadding
	}
	for i := 0; i < padLen; i++ {
		if data[n-1-i] != byte(padLen) {
			return nil, errPadding
		}
	}
	return data[:n-padLen], nil
}
