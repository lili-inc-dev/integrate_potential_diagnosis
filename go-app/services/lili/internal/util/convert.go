package util

import (
	"encoding/binary"
)

func Uint64ToBytes(n uint64) []byte {
	bytes := make([]byte, binary.MaxVarintLen64)
	binary.BigEndian.PutUint64(bytes, n)
	return bytes
}
