package engines

import "encoding/binary"

func buildNonce(connectionID int64, counter uint64, size int) []byte {
	nonce := make([]byte, size)

	binary.LittleEndian.PutUint64(nonce[0:8], uint64(connectionID))
	binary.LittleEndian.PutUint64(nonce[8:16], counter)

	return nonce[:size]
}
