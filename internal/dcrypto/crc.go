package dcrypto

import (
	"hash/crc32"
	"hash/crc64"
)

func 计算CRC32(数据 string) uint32 {
	return crc32.ChecksumIEEE([]byte(数据))
}

func 计算CRC64(数据 string) uint64 {
	表 := crc64.MakeTable(crc64.ECMA)
	return crc64.Checksum([]byte(数据), 表)
}
