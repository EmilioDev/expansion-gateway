package others

import (
	"expansion-gateway/config"
	disp "expansion-gateway/interfaces/dispatchers"
	"expansion-gateway/interfaces/packets"
	"expansion-gateway/internal/dispatchers"
)

func NewShardedPacketDispatcher(cfg *config.Configuration) (disp.Dispatcher[packets.Packet], disp.Reciver[packets.Packet]) {
	shardCount := cfg.GetShardCount()
	bufferSize := cfg.GetShardBufferSize()

	shards := make([]chan packets.Packet, shardCount)
	for i := range shards {
		shards[i] = make(chan packets.Packet, bufferSize)
	}

	dispatcher := &dispatchers.PacketDispatcher{
		Shards: shards,
		Count:  shardCount,
	}

	reciver := &dispatchers.PackerReciver{
		Shards: shards,
	}

	return dispatcher, reciver
}
