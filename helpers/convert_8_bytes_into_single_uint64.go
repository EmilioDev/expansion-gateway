package helpers

import "encoding/binary"

func Convert8BytesIntoSingleUInt64(byte0, byte1, byte2, byte3, byte4, byte5, byte6, byte7 byte) uint64 {
	array := [8]byte{byte0, byte1, byte2, byte3, byte4, byte5, byte6, byte7}
	return binary.BigEndian.Uint64(array[:])
}

func ConvertBytesArrayIntoSingleUInt64(byteArray [8]byte) uint64 {
	return binary.BigEndian.Uint64(byteArray[:])
}
