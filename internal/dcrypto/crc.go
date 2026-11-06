package dcrypto

import (
	"hash/crc32"
	"hash/crc64"
)

func crc32Of(data string) uint32 {
	return crc32.ChecksumIEEE([]byte(data))
}

func crc64Of(data string) uint64 {
	table := crc64.MakeTable(crc64.ECMA)
	return crc64.Checksum([]byte(data), table)
}
