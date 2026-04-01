package helpers

import "encoding/binary"

func ConvertUint16Into2Bytes(input uint16) (result [2]byte) {
	binary.BigEndian.PutUint16(result[:], input)
	return
}
