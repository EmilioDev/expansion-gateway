// file: /interfaces/packets/packet.go
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
	GetRawPayload() []byte                  // the payload of the packet as a ray byte array
	GetIdentifier() string                  // the subscription key identifier of the packet
	SetNewOwner(int64)                      // changes the owner of this packet
}
