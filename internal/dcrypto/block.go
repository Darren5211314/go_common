package dcrypto

import (
	"crypto/cipher"
)

func encryptCBC(block cipher.Block, data []byte, IV []byte) ([]byte, error) {
	blockSize := block.BlockSize()
	if len(IV) != blockSize {
		return nil, errIVLen
	}
	padded := pkcs7Pad(data, blockSize)
	out := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, IV).CryptBlocks(out, padded)
	return out, nil
}

func decryptCBC(block cipher.Block, data []byte, IV []byte) ([]byte, error) {
	blockSize := block.BlockSize()
	if len(IV) != blockSize {
		return nil, errIVLen
	}
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, errCiphertext
	}
	out := make([]byte, len(data))
	cipher.NewCBCDecrypter(block, IV).CryptBlocks(out, data)
	return pkcs7Unpad(out, blockSize)
}

func encryptECB(block cipher.Block, data []byte) ([]byte, error) {
	blockSize := block.BlockSize()
	padded := pkcs7Pad(data, blockSize)
	out := make([]byte, len(padded))
	for offset := 0; offset < len(padded); offset += blockSize {
		block.Encrypt(out[offset:offset+blockSize], padded[offset:offset+blockSize])
	}
	return out, nil
}

func decryptECB(block cipher.Block, data []byte) ([]byte, error) {
	blockSize := block.BlockSize()
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, errCiphertext
	}
	out := make([]byte, len(data))
	for offset := 0; offset < len(data); offset += blockSize {
		block.Decrypt(out[offset:offset+blockSize], data[offset:offset+blockSize])
	}
	return pkcs7Unpad(out, blockSize)
}
