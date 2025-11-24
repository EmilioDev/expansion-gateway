package helpers

import "encoding/base64"

// ===== Convenience helpers for logging / wire encoding (base64) =====

// encodes a byte array into a base-64 string
func EncodeToBase64String(rawInput []byte) string {
	return base64.RawStdEncoding.EncodeToString(rawInput)
}

// encodes a public 32-bytes key
func PubKeyToBase64(pub [32]byte) string {
	return base64.RawStdEncoding.EncodeToString(pub[:])
}

// encodes a private 32-bytes key
func PrivKeyToBase64(priv [32]byte) string {
	return base64.RawStdEncoding.EncodeToString(priv[:])
}
