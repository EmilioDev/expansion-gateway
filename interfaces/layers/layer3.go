package layers

import (
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

type Layer3 interface {
	DumbLayer[packets.OutputPacket]
	Publish(data packets.OutputPacket) errorinfo.GatewayError
	SubscribeTo(topic string) errorinfo.GatewayError
	UnsubscribeTo(topic string) errorinfo.GatewayError
	RenameGateway(newName string)
}
