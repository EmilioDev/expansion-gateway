package layers

import (
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

type Layer1 interface {
	DumbLayer
	SendPacket(packets.Packet) errorinfo.GatewayError
}
