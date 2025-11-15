package helpers

import "encoding/binary"

func ConvertInt32Into4Bytes(input int32) [4]byte {
	var result [4]byte

	binary.BigEndian.PutUint32(result[:], uint32(input))

	return result
}
