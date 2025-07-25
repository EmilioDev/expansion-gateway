package interfaces

import "expansion-gateway/interfaces/packets"

type OneWayPipe interface {
	// Initialices this layer with the cannel
	InitOutputChannel(channel packets.Packet)
}
