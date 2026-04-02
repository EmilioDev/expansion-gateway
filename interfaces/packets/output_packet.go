package packets

import (
	"expansion-gateway/internal/structs/tries"
)

type OutputPacket interface {
	GetKey() tries.SubscriptionKey
	GetPayload() []byte
}
