package bytearraytopacket

import (
	"expansion-gateway/dto"
	"expansion-gateway/enums"
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

func ToDisconnectPacket(byteArray *[]byte, connectionID int64) (packets.Packet, errorinfo.GatewayError) {
	byteArraySize := len(*byteArray)
	const filePath string = "/parsers/byteArrayToPacket/to_disconnect_packet.go"

	if byteArraySize != 2 {
		return nil, errors.CreateInvalidPacketSizeError(filePath, 15, enums.NONE, byteArraySize)
	}

	if !enums.IsValidDisconnectReason((*byteArray)[1]) {
		return nil, errors.CreateInvalidDisconnectReasonError(filePath, 19, (*byteArray)[1])
	}

	return dto.CreateDisconnectPacket(connectionID, enums.DisconnectReason((*byteArray)[1])), nil
}
