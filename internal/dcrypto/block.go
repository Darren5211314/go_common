package dcrypto

import (
	"crypto/cipher"
)

func 加密CBC(块 cipher.Block, 明文 []byte, IV []byte) ([]byte, error) {
	块大小 := 块.BlockSize()
	if len(IV) != 块大小 {
		return nil, 错误IV长度
	}
	填充后 := PKCS7填充(明文, 块大小)
	密文 := make([]byte, len(填充后))
	cipher.NewCBCEncrypter(块, IV).CryptBlocks(密文, 填充后)
	return 密文, nil
}

func 解密CBC(块 cipher.Block, 密文 []byte, IV []byte) ([]byte, error) {
	块大小 := 块.BlockSize()
	if len(IV) != 块大小 {
		return nil, 错误IV长度
	}
	if len(密文) == 0 || len(密文)%块大小 != 0 {
		return nil, 错误密文
	}
	明文 := make([]byte, len(密文))
	cipher.NewCBCDecrypter(块, IV).CryptBlocks(明文, 密文)
	return PKCS7去填充(明文, 块大小)
}

func 加密ECB(块 cipher.Block, 明文 []byte) ([]byte, error) {
	块大小 := 块.BlockSize()
	填充后 := PKCS7填充(明文, 块大小)
	密文 := make([]byte, len(填充后))
	for 偏移 := 0; 偏移 < len(填充后); 偏移 += 块大小 {
		块.Encrypt(密文[偏移:偏移+块大小], 填充后[偏移:偏移+块大小])
	}
	return 密文, nil
}

func 解密ECB(块 cipher.Block, 密文 []byte) ([]byte, error) {
	块大小 := 块.BlockSize()
	if len(密文) == 0 || len(密文)%块大小 != 0 {
		return nil, 错误密文
	}
	明文 := make([]byte, len(密文))
	for 偏移 := 0; 偏移 < len(密文); 偏移 += 块大小 {
		块.Decrypt(明文[偏移:偏移+块大小], 密文[偏移:偏移+块大小])
	}
	return PKCS7去填充(明文, 块大小)
}
