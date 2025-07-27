package packets

import (
	"expansion-gateway/enums"
	errors "expansion-gateway/interfaces/errorinfo"
)

// Base interface of all the packets
type Packet interface {
	GetPacketType() enums.PacketType        // type of packet
	GetPacketHeader() PacketHeader          // packet header
	GetPayload() string                     // packet payload
	Marshal() ([]byte, errors.GatewayError) // convert the packet into a byte array
	GetSender() int64                       // the sender of the packet
}
