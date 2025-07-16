package parsers

import (
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/packets"
)

type BasicByteArrayToPacketParser struct{}

func (parser *BasicByteArrayToPacketParser) ParseByteArrayToPacket(byteArray *[]byte, connectionID int64) (*packets.Packet, errors.GatewayError) {
	// implement

	return nil, nil
}
