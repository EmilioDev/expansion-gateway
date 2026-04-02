package layers

import (
	"expansion-gateway/interfaces/dispatchers"
	"expansion-gateway/interfaces/errorinfo"
)

type DumbLayer[T any] interface {
	Layer
	// these are the channels/dispatchers that will be used between this layer
	// and the bussines logic layer to comunicate
	CloseSession(sessionId int64) errorinfo.GatewayError
	MoveClientTo(origin, destiny int64)
	ConfigureDumbLayer(outputChannel dispatchers.Dispatcher[T], inputChannel dispatchers.Reciver[T]) errorinfo.GatewayError
}
