package crypto

import (
	"crypto/rand"
	"expansion-gateway/errors"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/errorinfo"

	"golang.org/x/crypto/curve25519"
)

// GenerateX25519Keypair generates a new ephemeral X25519 keypair.
// Returns private and public keys as [32]byte. The private key is already clamped.
func GenerateX25519KeyPair() (priv, pub [32]byte, err errorinfo.GatewayError) {
	const filePath string = "/crypto/generate_x25519_key_pair.go"

	// Fill private with secure random bytes
	if _, err := rand.Read(priv[:]); err != nil {
		return priv, pub, errors.CreateErrorWrapper(filePath, 15, err)
	}

	// RFC7748 clamping (defensive; X25519 does its own clamping, but being explicit is ok)
	priv[0] &= 248
	priv[31] &= 127
	priv[31] |= 64

	// standard basepoint: first byte 9, rest 0
	var basepoint [32]byte
	basepoint[0] = 9

	// compute public key
	pubSlice, err2 := curve25519.X25519(priv[:], basepoint[:])

	if err2 != nil {
		// wipe private before returning on error
		helpers.ZeroKey(&priv)
		return priv, pub, errors.CreateErrorWrapper(filePath, 33, err2)
	}

	copy(pub[:], pubSlice)

	return priv, pub, nil
}
