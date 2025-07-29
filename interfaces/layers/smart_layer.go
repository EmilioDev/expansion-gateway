package layers

import "expansion-gateway/interfaces/packets"

type SmartLayer interface {
	Layer
	// method for configuring layer 1
	ConfigureFirstLayer(layer DumbLayer[packets.Packet])

	// ...later we add for layer 3
}
