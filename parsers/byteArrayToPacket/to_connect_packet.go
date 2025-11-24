// file: /parsers/byteArrayToPacket/to_connect_packet.go
package bytearraytopacket

import (
	"expansion-gateway/dto"
	"expansion-gateway/enums"
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

func ToConnectPacket(byteArray *[]byte, connectionID int64) (packets.Packet, errorinfo.GatewayError) {
	packetSize := len(*byteArray)
	const filePath string = "/parsers/byteArrayToPacket/to_connect_packet.go"

	var signature [64]byte

	switch packetSize {
	case 65:
		copy(signature[:], (*byteArray)[1:])

		return dto.CreateConnectPacketWithoutEphemeralKey(connectionID, signature), nil

	case 97:
		copy(signature[:], (*byteArray)[1:65])

		var ephemeralKey [32]byte
		copy(ephemeralKey[:], (*byteArray)[65:])

		return dto.CreateConnectPacketWithEphemeralKey(connectionID, signature, ephemeralKey), nil

	default:
		return nil, errors.CreateInvalidPacketSizeError(filePath, 32, enums.CONNECT, packetSize)
	}
}
