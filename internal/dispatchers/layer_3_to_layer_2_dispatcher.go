package dispatchers

import (
	"expansion-gateway/interfaces/packets"
	"sync/atomic"
)

type NatsDtoDispatcher struct {
	shards []chan packets.OutputPacket
	count  int
	index  atomic.Uint32
}

func (dispatcher *NatsDtoDispatcher) Dispatch(data packets.OutputPacket) {
	shardIndex := dispatcher.index.Load()
	index := shardIndex % uint32(dispatcher.count)

	dispatcher.shards[index] <- data

	dispatcher.index.Store(index + 1)
}

func CreateNatsDtoDispatcher(shards []chan packets.OutputPacket) *NatsDtoDispatcher {
	return &NatsDtoDispatcher{
		shards: shards,
		count:  len(shards),
		index:  atomic.Uint32{},
	}
}
