package helpers

import (
	"crypto/rand"
	"encoding/binary"
)

// generates a random int64. if the generation fails, it returns 0 (zero)
func GenerateRandomInt64() int64 {
	var buffer [8]byte
	var result int64 = 0

	if _, err := rand.Read(buffer[:]); err == nil {
		raw := binary.LittleEndian.Uint64(buffer[:])
		result = int64(raw)
	}

	return result
}
