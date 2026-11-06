package dcrypto

// AESNS AES 加解密命名空间。
type AESNS struct{}

func (AESNS) EncryptCBC(data, keyHex, IVhex string) (string, error) {
	return AESEncryptCBC(data, keyHex, IVhex)
}
func (AESNS) DecryptCBC(data, keyHex, IVhex string) (string, error) {
	return AESDecryptCBC(data, keyHex, IVhex)
}
func (AESNS) EncryptECB(data, keyHex string) (string, error) {
	return AESEncryptECB(data, keyHex)
}
func (AESNS) DecryptECB(data, keyHex string) (string, error) {
	return AESDecryptECB(data, keyHex)
}
func (AESNS) EncryptGCM(data, keyHex, NonceHex string) (string, error) {
	return AESEncryptGCM(data, keyHex, NonceHex)
}
func (AESNS) DecryptGCM(data, keyHex, NonceHex string) (string, error) {
	return AESDecryptGCM(data, keyHex, NonceHex)
}

// DESNS DES 加解密命名空间。
type DESNS struct{}

func (DESNS) EncryptCBC(data, keyHex, IVhex string) (string, error) {
	return DESEncryptCBC(data, keyHex, IVhex)
}
func (DESNS) DecryptCBC(data, keyHex, IVhex string) (string, error) {
	return DESDecryptCBC(data, keyHex, IVhex)
}
func (DESNS) EncryptECB(data, keyHex string) (string, error) {
	return DESEncryptECB(data, keyHex)
}
func (DESNS) DecryptECB(data, keyHex string) (string, error) {
	return DESDecryptECB(data, keyHex)
}

// DES3NS 3DES 加解密命名空间。
type DES3NS struct{}

func (DES3NS) EncryptCBC(data, keyHex, IVhex string) (string, error) {
	return DES3EncryptCBC(data, keyHex, IVhex)
}
func (DES3NS) DecryptCBC(data, keyHex, IVhex string) (string, error) {
	return DES3DecryptCBC(data, keyHex, IVhex)
}
func (DES3NS) EncryptECB(data, keyHex string) (string, error) {
	return DES3EncryptECB(data, keyHex)
}
func (DES3NS) DecryptECB(data, keyHex string) (string, error) {
	return DES3DecryptECB(data, keyHex)
}

// HashNS 摘要算法命名空间。
type HashNS struct{}

func (HashNS) MD5(data string) string    { return MD5(data) }
func (HashNS) SHA1(data string) string   { return SHA1(data) }
func (HashNS) SHA256(data string) string { return SHA256(data) }
func (HashNS) SHA512(data string) string { return SHA512(data) }
func (HashNS) SM3(data string) string    { return SM3(data) }

// HMACNS HMAC 命名空间。
type HMACNS struct{}

func (HMACNS) SHA1(data, key string) (string, error) {
	return HMACSHA1(data, key)
}
func (HMACNS) SHA256(data, key string) (string, error) {
	return HMACSHA256(data, key)
}
func (HMACNS) SHA512(data, key string) (string, error) {
	return HMACSHA512(data, key)
}

// CRCNS CRC 校验命名空间。
type CRCNS struct{}

func (CRCNS) CRC32(data string) uint32 { return crc32Of(data) }
func (CRCNS) CRC64(data string) uint64 { return crc64Of(data) }

// RSANS RSA 加解密与签名命名空间。
type RSANS struct{}

func (RSANS) GenerateKeyPair(bits int) (publicKeyPEM, privateKeyPEM string, err error) {
	return RSAGenerateKeyPair(bits)
}
func (RSANS) Encrypt(data, publicKeyPEM string) (string, error) {
	return RSAEncrypt(data, publicKeyPEM)
}
func (RSANS) Decrypt(data, privateKeyPEM string) (string, error) {
	return RSADecrypt(data, privateKeyPEM)
}
func (RSANS) Sign(data, privateKeyPEM string) (string, error) {
	return RSASign(data, privateKeyPEM)
}
func (RSANS) Verify(data, signatureHex, publicKeyPEM string) error {
	return RSAVerify(data, signatureHex, publicKeyPEM)
}

// SM4NS SM4 加解密命名空间。
type SM4NS struct{}

func (SM4NS) EncryptCBC(data, keyHex, IVhex string) (string, error) {
	return SM4EncryptCBC(data, keyHex, IVhex)
}
func (SM4NS) DecryptCBC(data, keyHex, IVhex string) (string, error) {
	return SM4DecryptCBC(data, keyHex, IVhex)
}
func (SM4NS) EncryptECB(data, keyHex string) (string, error) {
	return SM4EncryptECB(data, keyHex)
}
func (SM4NS) DecryptECB(data, keyHex string) (string, error) {
	return SM4DecryptECB(data, keyHex)
}
func (SM4NS) EncryptGCM(data, keyHex, NonceHex string) (string, error) {
	return SM4EncryptGCM(data, keyHex, NonceHex)
}
func (SM4NS) DecryptGCM(data, keyHex, NonceHex string) (string, error) {
	return SM4DecryptGCM(data, keyHex, NonceHex)
}
