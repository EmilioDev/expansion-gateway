package layers

import "expansion-gateway/interfaces/packets"

type Layer1 interface {
	DumbLayer[packets.Packet]
}
