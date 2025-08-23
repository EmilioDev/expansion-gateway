// file: /interfaces/packets/encryptable_packet.go
package packets

import "expansion-gateway/internal/crypto"

type EncryptablePacket interface {
	Packet
	GetCryptoAppendix() *crypto.CryptoAppendix
}
