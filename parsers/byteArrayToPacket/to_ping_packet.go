package bytearraytopacket

import (
	"expansion-gateway/dto"
	"expansion-gateway/enums"
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

func ToPingPacket(byteArray *[]byte, connectionID int64) (packets.Packet, errorinfo.GatewayError) {
	const filePath string = "/parsers/byteArrayToPacket/to_ping_packet.go"
	byteArraySize := len(*byteArray)

	if byteArraySize != 1 {
		return nil, errors.CreateInvalidPacketSizeError(filePath, 15, enums.PING, byteArraySize)
	}

	return dto.CreatePingPacket(connectionID), nil
}
