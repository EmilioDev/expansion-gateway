package layers

import (
	"expansion-gateway/interfaces/commands"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

type DumbLayer[O packets.BaseMessage] interface {
	Layer
	// these are the channel that will be used between this layer
	// and the bussines logic layer to comunicate
	ConfigureDumbLayer(outputChannel chan<- O, inputChannel <-chan commands.Command) errorinfo.GatewayError
}
