package layers

import (
	"expansion-gateway/interfaces/dispatchers"
	"expansion-gateway/interfaces/errorinfo"
)

type DumbLayer interface {
	Layer
	// these are the channels/dispatchers that will be used between this layer
	// and the bussines logic layer to comunicate
	ConfigureDumbLayer(outputChannel dispatchers.Dispatcher, inputChannel dispatchers.Reciver) errorinfo.GatewayError
	CloseSession(sessionId int64) errorinfo.GatewayError
	MoveClientTo(origin, destiny int64)
}
