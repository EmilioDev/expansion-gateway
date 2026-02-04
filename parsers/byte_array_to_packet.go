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
	const filePath string = "/parsers/byte_array_to_packet.go"

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

	case enums.PING:
		return bytearraytopacket.ToPingPacket(byteArray, connectionID)

	case enums.PINGACK:
		return bytearraytopacket.ToPingAckPacket(byteArray, connectionID)

	case enums.SUBSCRIBE:
		return bytearraytopacket.ToSubscribePacket(byteArray, connectionID)

	case enums.UNSUBSCRIBE:
		return bytearraytopacket.ToUnsubscribePacket(byteArray, connectionID)

	default:
		return nil, errors.CreateInvalidPacketError(filePath, 22)
	}
}
