package dcrypto

// AES类 AES 加解密命名空间。
type AES类 struct{}

func (AES类) EncryptCBC(明文, 密钥hex, IVhex string) (string, error) {
	return AESEncryptCBC(明文, 密钥hex, IVhex)
}
func (AES类) DecryptCBC(密文hex, 密钥hex, IVhex string) (string, error) {
	return AESDecryptCBC(密文hex, 密钥hex, IVhex)
}
func (AES类) EncryptECB(明文, 密钥hex string) (string, error) {
	return AESEncryptECB(明文, 密钥hex)
}
func (AES类) DecryptECB(密文hex, 密钥hex string) (string, error) {
	return AESDecryptECB(密文hex, 密钥hex)
}
func (AES类) EncryptGCM(明文, 密钥hex, NonceHex string) (string, error) {
	return AESEncryptGCM(明文, 密钥hex, NonceHex)
}
func (AES类) DecryptGCM(密文hex, 密钥hex, NonceHex string) (string, error) {
	return AESDecryptGCM(密文hex, 密钥hex, NonceHex)
}

// DES类 DES 加解密命名空间。
type DES类 struct{}

func (DES类) EncryptCBC(明文, 密钥hex, IVhex string) (string, error) {
	return DESEncryptCBC(明文, 密钥hex, IVhex)
}
func (DES类) DecryptCBC(密文hex, 密钥hex, IVhex string) (string, error) {
	return DESDecryptCBC(密文hex, 密钥hex, IVhex)
}
func (DES类) EncryptECB(明文, 密钥hex string) (string, error) {
	return DESEncryptECB(明文, 密钥hex)
}
func (DES类) DecryptECB(密文hex, 密钥hex string) (string, error) {
	return DESDecryptECB(密文hex, 密钥hex)
}

// DES3类 3DES 加解密命名空间。
type DES3类 struct{}

func (DES3类) EncryptCBC(明文, 密钥hex, IVhex string) (string, error) {
	return DES3EncryptCBC(明文, 密钥hex, IVhex)
}
func (DES3类) DecryptCBC(密文hex, 密钥hex, IVhex string) (string, error) {
	return DES3DecryptCBC(密文hex, 密钥hex, IVhex)
}
func (DES3类) EncryptECB(明文, 密钥hex string) (string, error) {
	return DES3EncryptECB(明文, 密钥hex)
}
func (DES3类) DecryptECB(密文hex, 密钥hex string) (string, error) {
	return DES3DecryptECB(密文hex, 密钥hex)
}

// Hash类 摘要算法命名空间。
type Hash类 struct{}

func (Hash类) MD5(数据 string) string    { return 摘要MD5(数据) }
func (Hash类) SHA1(数据 string) string   { return 摘要SHA1(数据) }
func (Hash类) SHA256(数据 string) string { return 摘要SHA256(数据) }
func (Hash类) SHA512(数据 string) string { return 摘要SHA512(数据) }
func (Hash类) SM3(数据 string) string    { return 摘要SM3(数据) }

// HMAC类 HMAC 命名空间。
type HMAC类 struct{}

func (HMAC类) SHA1(数据, 密钥 string) (string, error) {
	return HMACSHA1(数据, 密钥)
}
func (HMAC类) SHA256(数据, 密钥 string) (string, error) {
	return HMACSHA256(数据, 密钥)
}
func (HMAC类) SHA512(数据, 密钥 string) (string, error) {
	return HMACSHA512(数据, 密钥)
}

// CRC类 CRC 校验命名空间。
type CRC类 struct{}

func (CRC类) CRC32(数据 string) uint32 { return 计算CRC32(数据) }
func (CRC类) CRC64(数据 string) uint64 { return 计算CRC64(数据) }

// RSA类 RSA 加解密与签名命名空间。
type RSA类 struct{}

func (RSA类) GenerateKeyPair(位数 int) (公钥PEM, 私钥PEM string, err error) {
	return RSAGenerateKeyPair(位数)
}
func (RSA类) Encrypt(明文, 公钥PEM string) (string, error) {
	return RSAEncrypt(明文, 公钥PEM)
}
func (RSA类) Decrypt(密文hex, 私钥PEM string) (string, error) {
	return RSADecrypt(密文hex, 私钥PEM)
}
func (RSA类) Sign(数据, 私钥PEM string) (string, error) {
	return RSASign(数据, 私钥PEM)
}
func (RSA类) Verify(数据, 签名hex, 公钥PEM string) error {
	return RSAVerify(数据, 签名hex, 公钥PEM)
}

// SM4类 SM4 加解密命名空间。
type SM4类 struct{}

func (SM4类) EncryptCBC(明文, 密钥hex, IVhex string) (string, error) {
	return SM4EncryptCBC(明文, 密钥hex, IVhex)
}
func (SM4类) DecryptCBC(密文hex, 密钥hex, IVhex string) (string, error) {
	return SM4DecryptCBC(密文hex, 密钥hex, IVhex)
}
func (SM4类) EncryptECB(明文, 密钥hex string) (string, error) {
	return SM4EncryptECB(明文, 密钥hex)
}
func (SM4类) DecryptECB(密文hex, 密钥hex string) (string, error) {
	return SM4DecryptECB(密文hex, 密钥hex)
}
func (SM4类) EncryptGCM(明文, 密钥hex, NonceHex string) (string, error) {
	return SM4EncryptGCM(明文, 密钥hex, NonceHex)
}
func (SM4类) DecryptGCM(密文hex, 密钥hex, NonceHex string) (string, error) {
	return SM4DecryptGCM(密文hex, 密钥hex, NonceHex)
}
