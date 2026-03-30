package crypto

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
)

type CryptoEngine interface {
	Encrypt(counter uint64, data []byte) ([]byte, errorinfo.GatewayError) // encrypts data, and returns that data encrypted, or an error
	Decrypt(counter uint64, data []byte) ([]byte, errorinfo.GatewayError) // decrypts data, and returns that data decrypted, or an error
	EncryptionSupported() enums.EncryptionAlgorithm                       // says the kind of cryptographic algorithm supported by this engine
}
