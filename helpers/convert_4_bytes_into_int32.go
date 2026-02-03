package helpers

import "encoding/binary"

// converts 4 bytes into an int32
func Convert4bytesIntoInt32(input [4]byte) int32 {
	return int32(binary.BigEndian.Uint32(input[:]))
}

// converts 4 bytes into an uint32
func Convert4bytesIntoUint32(input [4]byte) uint32 {
	return binary.BigEndian.Uint32(input[:])
}
