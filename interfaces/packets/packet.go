package packets

import (
	"expansion-gateway/enums"
	errors "expansion-gateway/interfaces/errorinfo"
)

// Base interface of all the packets
type Packet interface {
	GetPacketType() enums.PacketType
	GetPacketHeader() PacketHeader
	GetPayload() string
	Marshal() ([]byte, errors.GatewayError)
}
