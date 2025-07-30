package layers

import (
	"expansion-gateway/interfaces/commands"
	"expansion-gateway/interfaces/dispatchers"
	"expansion-gateway/interfaces/errorinfo"
)

type DumbLayer interface {
	Layer
	// these are the channels/dispatchers that will be used between this layer
	// and the bussines logic layer to comunicate
	ConfigureDumbLayer(outputChannel dispatchers.Dispatcher, inputChannel <-chan commands.Command) errorinfo.GatewayError
}
