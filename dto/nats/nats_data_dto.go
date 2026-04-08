package nats

import (
	"expansion-gateway/internal/structs/tries"
)

type NatsDataDto struct {
	Payload []byte
	Key     tries.SubscriptionKey
	sender  int64
}

func (data *NatsDataDto) GetPayload() []byte {
	return data.Payload
}

func (data *NatsDataDto) GetKey() tries.SubscriptionKey {
	return data.Key
}

func (data *NatsDataDto) GetSender() int64 {
	return data.sender
}

func CreateNewNatsDataBasicTransferRecipe(key tries.SubscriptionKey, payload []byte) *NatsDataDto {
	return &NatsDataDto{
		Payload: payload,
		Key:     key,
		sender:  0,
	}
}

func CreateNewNatsDataTransferRecipe(key tries.SubscriptionKey, owner int64, payload []byte) *NatsDataDto {
	return &NatsDataDto{
		Payload: payload,
		Key:     key,
		sender:  owner,
	}
}
