package nats

import (
	"expansion-gateway/internal/structs/tries"
)

type NatsDataDto struct {
	Payload []byte
	Key     tries.SubscriptionKey
}

func (data *NatsDataDto) GetPayload() []byte {
	return data.Payload
}

func (data *NatsDataDto) GetKey() tries.SubscriptionKey {
	return data.Key
}

func CreateNewNatsDataTransferRecipe(key tries.SubscriptionKey, payload []byte) *NatsDataDto {
	return &NatsDataDto{
		Payload: payload,
		Key:     key,
	}
}
