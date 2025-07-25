package helpers

import "encoding/binary"

func ConvertInt64Into8Bytes(input int64) [8]byte {
	var b [8]byte

	binary.BigEndian.PutUint64(b[:], uint64(input))

	return b
}
