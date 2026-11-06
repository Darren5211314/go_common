package dcrypto

// AES类 AES 加解密命名空间。
type AES类 struct{}

func (AES类) EncryptCBC(data, keyHex, IVhex string) (string, error) {
	return AESEncryptCBC(data, keyHex, IVhex)
}
func (AES类) DecryptCBC(data, keyHex, IVhex string) (string, error) {
	return AESDecryptCBC(data, keyHex, IVhex)
}
func (AES类) EncryptECB(data, keyHex string) (string, error) {
	return AESEncryptECB(data, keyHex)
}
func (AES类) DecryptECB(data, keyHex string) (string, error) {
	return AESDecryptECB(data, keyHex)
}
func (AES类) EncryptGCM(data, keyHex, NonceHex string) (string, error) {
	return AESEncryptGCM(data, keyHex, NonceHex)
}
func (AES类) DecryptGCM(data, keyHex, NonceHex string) (string, error) {
	return AESDecryptGCM(data, keyHex, NonceHex)
}

// DES类 DES 加解密命名空间。
type DES类 struct{}

func (DES类) EncryptCBC(data, keyHex, IVhex string) (string, error) {
	return DESEncryptCBC(data, keyHex, IVhex)
}
func (DES类) DecryptCBC(data, keyHex, IVhex string) (string, error) {
	return DESDecryptCBC(data, keyHex, IVhex)
}
func (DES类) EncryptECB(data, keyHex string) (string, error) {
	return DESEncryptECB(data, keyHex)
}
func (DES类) DecryptECB(data, keyHex string) (string, error) {
	return DESDecryptECB(data, keyHex)
}

// DES3类 3DES 加解密命名空间。
type DES3类 struct{}

func (DES3类) EncryptCBC(data, keyHex, IVhex string) (string, error) {
	return DES3EncryptCBC(data, keyHex, IVhex)
}
func (DES3类) DecryptCBC(data, keyHex, IVhex string) (string, error) {
	return DES3DecryptCBC(data, keyHex, IVhex)
}
func (DES3类) EncryptECB(data, keyHex string) (string, error) {
	return DES3EncryptECB(data, keyHex)
}
func (DES3类) DecryptECB(data, keyHex string) (string, error) {
	return DES3DecryptECB(data, keyHex)
}

// Hash类 摘要算法命名空间。
type Hash类 struct{}

func (Hash类) MD5(data string) string    { return MD5(data) }
func (Hash类) SHA1(data string) string   { return SHA1(data) }
func (Hash类) SHA256(data string) string { return SHA256(data) }
func (Hash类) SHA512(data string) string { return SHA512(data) }
func (Hash类) SM3(data string) string    { return SM3(data) }

// HMAC类 HMAC 命名空间。
type HMAC类 struct{}

func (HMAC类) SHA1(data, key string) (string, error) {
	return HMACSHA1(data, key)
}
func (HMAC类) SHA256(data, key string) (string, error) {
	return HMACSHA256(data, key)
}
func (HMAC类) SHA512(data, key string) (string, error) {
	return HMACSHA512(data, key)
}

// CRC类 CRC 校验命名空间。
type CRC类 struct{}

func (CRC类) CRC32(data string) uint32 { return crc32Of(data) }
func (CRC类) CRC64(data string) uint64 { return crc64Of(data) }

// RSA类 RSA 加解密与签名命名空间。
type RSA类 struct{}

func (RSA类) GenerateKeyPair(bits int) (publicKeyPEM, privateKeyPEM string, err error) {
	return RSAGenerateKeyPair(bits)
}
func (RSA类) Encrypt(data, publicKeyPEM string) (string, error) {
	return RSAEncrypt(data, publicKeyPEM)
}
func (RSA类) Decrypt(data, privateKeyPEM string) (string, error) {
	return RSADecrypt(data, privateKeyPEM)
}
func (RSA类) Sign(data, privateKeyPEM string) (string, error) {
	return RSASign(data, privateKeyPEM)
}
func (RSA类) Verify(data, signatureHex, publicKeyPEM string) error {
	return RSAVerify(data, signatureHex, publicKeyPEM)
}

// SM4类 SM4 加解密命名空间。
type SM4类 struct{}

func (SM4类) EncryptCBC(data, keyHex, IVhex string) (string, error) {
	return SM4EncryptCBC(data, keyHex, IVhex)
}
func (SM4类) DecryptCBC(data, keyHex, IVhex string) (string, error) {
	return SM4DecryptCBC(data, keyHex, IVhex)
}
func (SM4类) EncryptECB(data, keyHex string) (string, error) {
	return SM4EncryptECB(data, keyHex)
}
func (SM4类) DecryptECB(data, keyHex string) (string, error) {
	return SM4DecryptECB(data, keyHex)
}
func (SM4类) EncryptGCM(data, keyHex, NonceHex string) (string, error) {
	return SM4EncryptGCM(data, keyHex, NonceHex)
}
func (SM4类) DecryptGCM(data, keyHex, NonceHex string) (string, error) {
	return SM4DecryptGCM(data, keyHex, NonceHex)
}
