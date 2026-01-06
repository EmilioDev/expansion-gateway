package parsers

import (
	"expansion-gateway/enums"
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
	bytearraytopacket "expansion-gateway/parsers/byteArrayToPacket"
)

type BasicByteArrayToPacketParser struct{}

func (parser *BasicByteArrayToPacketParser) ParseByteArrayToPacket(byteArray *[]byte, connectionID int64) (packets.Packet, errorinfo.GatewayError) {
	byteArraySize := len(*byteArray)
	const filePath string = "/parsers/byte_array_to_packet.go"

	if byteArraySize == 0 {
		return nil, errors.CreateInvalidPacketSizeError(filePath, 11, enums.NONE, byteArraySize)
	}

	indexByte := (*byteArray)[0]

	switch enums.PacketType(indexByte) {
	case enums.HELLO:
		return bytearraytopacket.ToHelloPacket(byteArray, connectionID)

	case enums.CONNECT:
		return bytearraytopacket.ToConnectPacket(byteArray, connectionID)

	case enums.REDIRECTED:
		return bytearraytopacket.ToRedirectedPacket(byteArray, connectionID)

	case enums.DISCONNECT:
		return bytearraytopacket.ToDisconnectPacket(byteArray, connectionID)

	default:
		return nil, errors.CreateInvalidPacketError(filePath, 22)
	}
}
