package dcrypto

import (
	"testing"

	"github.com/Darren5211314/go_common/internal/dlog"
)

func TestHash(t *testing.T) {
	// "abc" 的标准向量
	md5 := 摘要MD5("abc")
	dlog.Info("MD5", "结果", md5, "期望", "900150983cd24fb0d6963f7d28e17f72")
	if md5 != "900150983cd24fb0d6963f7d28e17f72" {
		t.Errorf("MD5 不符")
	}
	sha256 := 摘要SHA256("abc")
	dlog.Info("SHA256", "结果", sha256)
	if sha256 != "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad" {
		t.Errorf("SHA256 不符")
	}
	sm3期望 := "66c7f0f462eeedd9d1f2d46bdc10e4e24167c4875cf2f7a2297da02b8f4ba8e0"
	sm3 := 摘要SM3("abc")
	dlog.Info("SM3", "结果", sm3, "期望", sm3期望)
	if sm3 != sm3期望 {
		t.Errorf("SM3 不符")
	}
}

func TestCRC(t *testing.T) {
	v := 计算CRC32("123456789")
	dlog.Info("CRC32", "结果", v, "期望", uint32(0xcbf43926))
	if v != 0xcbf43926 {
		t.Errorf("CRC32 不符")
	}
}

func TestHMAC(t *testing.T) {
	mac, err := HMACSHA256("abc", "key")
	dlog.Info("HMAC.SHA256", "结果", mac, "错误", err)
	if err != nil || mac == "" {
		t.Errorf("HMAC.SHA256 失败")
	}
}

func TestAES往返(t *testing.T) {
	密钥 := "00112233445566778899aabbccddeeff"       // 16 bytes
	IV := "0102030405060708090a0b0c0d0e0f10"         // 16 bytes
	nonce := "0102030405060708090a0b0c"               // 12 bytes
	明文 := "hello-aes"

	cbc, err := AESEncryptCBC(明文, 密钥, IV)
	dlog.Info("AES.CBC加密", "密文", cbc, "错误", err)
	if err != nil {
		t.Fatalf("CBC加密失败: %v", err)
	}
	解, err := AESDecryptCBC(cbc, 密钥, IV)
	dlog.Info("AES.CBC解密", "明文", 解, "错误", err)
	if err != nil || 解 != 明文 {
		t.Errorf("CBC往返不符")
	}

	ecb, err := AESEncryptECB(明文, 密钥)
	解, err2 := AESDecryptECB(ecb, 密钥)
	dlog.Info("AES.ECB往返", "密文", ecb, "明文", 解, "错误1", err, "错误2", err2)
	if err != nil || err2 != nil || 解 != 明文 {
		t.Errorf("ECB往返不符")
	}

	gcm, err := AESEncryptGCM(明文, 密钥, nonce)
	解, err2 = AESDecryptGCM(gcm, 密钥, nonce)
	dlog.Info("AES.GCM往返", "密文", gcm, "明文", 解, "错误1", err, "错误2", err2)
	if err != nil || err2 != nil || 解 != 明文 {
		t.Errorf("GCM往返不符")
	}
}

func TestDES往返(t *testing.T) {
	密钥 := "0011223344556677"
	IV := "0102030405060708"
	明文 := "hello-des"
	cbc, err := DESEncryptCBC(明文, 密钥, IV)
	解, err2 := DESDecryptCBC(cbc, 密钥, IV)
	dlog.Info("DES.CBC往返", "密文", cbc, "明文", 解, "错误1", err, "错误2", err2)
	if err != nil || err2 != nil || 解 != 明文 {
		t.Errorf("DES CBC往返不符")
	}
	ecb, err := DESEncryptECB(明文, 密钥)
	解, err2 = DESDecryptECB(ecb, 密钥)
	dlog.Info("DES.ECB往返", "密文", ecb, "明文", 解)
	if err != nil || err2 != nil || 解 != 明文 {
		t.Errorf("DES ECB往返不符")
	}
}

func TestDES3往返(t *testing.T) {
	密钥 := "00112233445566778899aabbccddeeff0011223344556677" // 24 bytes
	IV := "0102030405060708"
	明文 := "hello-3des"
	cbc, err := DES3EncryptCBC(明文, 密钥, IV)
	解, err2 := DES3DecryptCBC(cbc, 密钥, IV)
	dlog.Info("DES3.CBC往返", "密文", cbc, "明文", 解, "错误1", err, "错误2", err2)
	if err != nil || err2 != nil || 解 != 明文 {
		t.Errorf("DES3 CBC往返不符")
	}
}

func TestSM4往返(t *testing.T) {
	密钥 := "0123456789abcdeffedcba9876543210" // 16 bytes
	IV := "000102030405060708090a0b0c0d0e0f"
	nonce := "0102030405060708090a0b0c"
	明文 := "hello-sm4"

	cbc, err := SM4EncryptCBC(明文, 密钥, IV)
	解, err2 := SM4DecryptCBC(cbc, 密钥, IV)
	dlog.Info("SM4.CBC往返", "密文", cbc, "明文", 解, "错误1", err, "错误2", err2)
	if err != nil || err2 != nil || 解 != 明文 {
		t.Errorf("SM4 CBC往返不符")
	}

	ecb, err := SM4EncryptECB(明文, 密钥)
	解, err2 = SM4DecryptECB(ecb, 密钥)
	dlog.Info("SM4.ECB往返", "密文", ecb, "明文", 解, "错误1", err, "错误2", err2)
	if err != nil || err2 != nil || 解 != 明文 {
		t.Errorf("SM4 ECB往返不符")
	}

	gcm, err := SM4EncryptGCM(明文, 密钥, nonce)
	解, err2 = SM4DecryptGCM(gcm, 密钥, nonce)
	dlog.Info("SM4.GCM往返", "密文", gcm, "明文", 解, "错误1", err, "错误2", err2)
	if err != nil || err2 != nil || 解 != 明文 {
		t.Errorf("SM4 GCM往返不符")
	}
}

func TestRSA往返(t *testing.T) {
	pub, priv, err := RSAGenerateKeyPair(1024)
	dlog.Info("RSA生成", "错误", err, "公钥长度", len(pub), "私钥长度", len(priv))
	if err != nil {
		t.Fatalf("生成失败: %v", err)
	}
	密文, err := RSAEncrypt("hello-rsa", pub)
	明文, err2 := RSADecrypt(密文, priv)
	dlog.Info("RSA加解密", "密文", 密文, "明文", 明文, "错误1", err, "错误2", err2)
	if err != nil || err2 != nil || 明文 != "hello-rsa" {
		t.Errorf("RSA加解密不符")
	}
	签名, err := RSASign("hello-rsa", priv)
	验签 := RSAVerify("hello-rsa", 签名, pub)
	dlog.Info("RSA签名", "签名", 签名, "验签", 验签, "错误", err)
	if err != nil || 验签 != nil {
		t.Errorf("RSA签名验签不符")
	}
}
