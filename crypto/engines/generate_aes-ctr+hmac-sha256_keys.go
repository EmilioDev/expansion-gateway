package engines

import (
	"crypto/sha256"
	"expansion-gateway/errors/cryptoerror"
	"expansion-gateway/interfaces/errorinfo"
	"io"

	"golang.org/x/crypto/hkdf"
)

// derives both keys from a shared secret for AES-CTR + HMAC-SHA256
func deriveAesCtrHmacSha256Keys(sharedSecret [32]byte, info []byte) ([]byte, []byte, errorinfo.GatewayError) {
	hkdfReader := hkdf.New(sha256.New, sharedSecret[:], nil, info) // hkdf reader to derive 64 bytes keys (32+32)
	keys := make([]byte, 64)                                       // 32 bytes encryption + 32 bytes HMAC

	if _, err := io.ReadFull(hkdfReader, keys); err != nil { // generate keys
		return nil, nil, cryptoerror.CreateAesCtrHmacSha256KeysNotGenerated(
			"/crypto/engines/generate_aes-ctr+hmac-sha256_keys.go",
			16,
		)
	}

	return keys[:32], keys[32:], nil
}
