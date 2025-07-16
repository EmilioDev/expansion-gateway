package packets

import (
	"expansion-gateway/enums"
)

// Base interface of all the packets
type Packet interface {
	GetPacketType() enums.PacketType
	GetPacketHeader() PacketHeader
	GetPayload() string
}
