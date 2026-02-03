package nats

import (
	"expansion-gateway/config"
	"expansion-gateway/interfaces/dispatchers"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

type NatsBasicHandlerLayer3 struct {
	// implement Layer3 interface
}

func (layer *NatsBasicHandlerLayer3) Publish(data packets.Packet) errorinfo.GatewayError {
	return nil
}

func (layer *NatsBasicHandlerLayer3) SubscribeTo(topic string) errorinfo.GatewayError {
	return nil
}

func (layer *NatsBasicHandlerLayer3) UnsubscribeTo(topic string) errorinfo.GatewayError {
	return nil
}

func (layer *NatsBasicHandlerLayer3) ConfigureDumbLayer(outputChannel dispatchers.Dispatcher, inputChannel dispatchers.Reciver) errorinfo.GatewayError {
	return nil
}

func (layer *NatsBasicHandlerLayer3) CloseSession(sessionId int64) errorinfo.GatewayError {
	return nil
}

func (layer *NatsBasicHandlerLayer3) MoveClientTo(origin, destiny int64) {
	//
}

func (layer *NatsBasicHandlerLayer3) Start() errorinfo.GatewayError {
	return nil
}

func (layer *NatsBasicHandlerLayer3) Stop() errorinfo.GatewayError {
	return nil
}

func (layer *NatsBasicHandlerLayer3) IsWorking() bool {
	return true
}

func GenerateNewNatsLayer3Output(configuration *config.Configuration) *NatsBasicHandlerLayer3 {
	return &NatsBasicHandlerLayer3{
		//
	}
}
