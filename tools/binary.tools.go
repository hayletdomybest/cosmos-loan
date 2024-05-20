package tools

import "encoding/binary"

func Uint64ToBinary(u uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, u)
	return bz
}
