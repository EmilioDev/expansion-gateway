package helpers

import "encoding/binary"

func Convert8BytesIntoSingleInt64(byte0, byte1, byte2, byte3, byte4, byte5, byte6, byte7 byte) int64 {
	return int64(binary.BigEndian.Uint64([]byte{byte0, byte1, byte2, byte3, byte4, byte5, byte6, byte7}))
}

func ConvertBytesArrayIntoSingleInt64(byteArray []byte) int64 {
	if len(byteArray) < 8 {
		return 0
	}

	return Convert8BytesIntoSingleInt64(
		byteArray[0],
		byteArray[1],
		byteArray[2],
		byteArray[3],
		byteArray[4],
		byteArray[5],
		byteArray[6],
		byteArray[7])
}
