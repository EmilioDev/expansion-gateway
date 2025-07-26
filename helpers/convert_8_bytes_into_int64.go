package helpers

import "encoding/binary"

func Convert8BytesIntoSingleInt64(byte0, byte1, byte2, byte3, byte4, byte5, byte6, byte7 byte) int64 {
	array := [8]byte{byte0, byte1, byte2, byte3, byte4, byte5, byte6, byte7}
	return int64(binary.BigEndian.Uint64(array[:]))
}

func ConvertBytesArrayIntoSingleInt64(byteArray []byte) int64 {
	if len(byteArray) < 8 {
		return 0
	}

	return int64(binary.BigEndian.Uint64(byteArray[:8]))
}
