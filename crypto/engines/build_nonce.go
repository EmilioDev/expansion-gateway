package engines

import "encoding/binary"

// generates a new nonce of a specified size from the connection id and the counter
func buildNonce(connectionID int64, counter uint64, size int) []byte {
	nonce := make([]byte, size)

	binary.BigEndian.PutUint64(nonce[0:8], uint64(connectionID))
	binary.BigEndian.PutUint64(nonce[8:16], counter)

	return nonce[:size]
}
