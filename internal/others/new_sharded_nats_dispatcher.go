package others

import (
	disp "expansion-gateway/interfaces/dispatchers"
	"expansion-gateway/interfaces/packets"
	"expansion-gateway/internal/dispatchers"
)

func NewShardedNatsDispatcher(shardCount, bufferSize int) (disp.Dispatcher[packets.OutputPacket], disp.Reciver[packets.OutputPacket]) {
	shards := make([]chan packets.OutputPacket, shardCount)

	for i := range shards {
		shards[i] = make(chan packets.OutputPacket, bufferSize)
	}

	return dispatchers.CreateNatsDtoDispatcher(shards), dispatchers.CreateNewShardsDtoReceiver(shards)
}
