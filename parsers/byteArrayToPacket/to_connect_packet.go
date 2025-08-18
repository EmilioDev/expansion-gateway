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

	if packetSize != 65 {
		return nil, errors.CreateInvalidPacketSizeError(filePath, 16, enums.CONNECT, packetSize)
	}

	var signature [64]byte
	copy(signature[:], (*byteArray)[1:65])

	return &dto.ConnectPacket{
		OwnerID:   connectionID,
		Signature: signature,
	}, nil
}
