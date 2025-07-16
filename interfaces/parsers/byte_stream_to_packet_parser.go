package parsers

import (
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/packets"
)

type ByteStreamToPacketParser interface {
	// converts a byte array into a packet
	ParseByteArrayToPacket(byteArray *[]byte, connectionID int64) (*packets.Packet, errors.GatewayError)
}
