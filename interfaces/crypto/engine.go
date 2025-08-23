// file: /interfaces/crypto/engine.go
package crypto

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

type CryptoEngine interface {
	Algorithm() enums.EncryptionAlgorithm
	Decrypt(packet packets.EncryptablePacket) (string, errorinfo.GatewayError)
	Encrypt(textToEncrypt string) ([]byte, errorinfo.GatewayError)
	HandshakeParams() []byte
}
