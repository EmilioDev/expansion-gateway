package engines

import (
	"encoding/binary"
	"expansion-gateway/helpers"
)

// generates a new nonce of a specified size from the connection id and the counter
func buildNonce(connectionID int64, counter uint64, size int) []byte {
	nonce := make([]byte, size)

	binary.BigEndian.PutUint64(nonce[0:8], uint64(connectionID))
	binary.BigEndian.PutUint64(nonce[8:16], counter)

	return nonce[:size]
}

func buildNonceAesGcm(connectionID int64, counter uint64) []byte {
	nonce := make([]byte, 12)

	binary.BigEndian.PutUint64(nonce[0:8], uint64(connectionID))
	binary.BigEndian.PutUint32(nonce[8:12], helpers.ConvertUint64IntoUint32(counter))

	return nonce
}
