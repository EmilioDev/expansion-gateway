package layers

import (
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

type Layer1 interface {
	DumbLayer[packets.Packet]
	SendPacket(packets.Packet) errorinfo.GatewayError
	DisableSession(int64)
}
