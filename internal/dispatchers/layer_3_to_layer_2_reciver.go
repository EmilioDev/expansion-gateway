package dispatchers

import (
	"expansion-gateway/interfaces/packets"
)

type NatsDtoReceiver struct {
	shards []chan packets.OutputPacket
}

func (receiver *NatsDtoReceiver) GetShard(index int) <-chan packets.OutputPacket {
	return receiver.shards[index]
}

func (receiver *NatsDtoReceiver) ShardCount() int {
	return len(receiver.shards)
}

func (receiver *NatsDtoReceiver) TotalPending() int {
	cont := 0

	for _, channel := range receiver.shards {
		cont += len(channel)
	}

	return cont
}

func CreateNewShardsDtoReceiver(shards []chan packets.OutputPacket) *NatsDtoReceiver {
	return &NatsDtoReceiver{
		shards: shards,
	}
}
