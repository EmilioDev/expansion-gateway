package dispatchers

import "expansion-gateway/interfaces/packets"

type PackerReciver struct {
	Shards []chan packets.Packet
}

func (r *PackerReciver) GetShard(index int) <-chan packets.Packet {
	return r.Shards[index]
}

func (r *PackerReciver) ShardCount() int {
	return len(r.Shards)
}

func (r *PackerReciver) TotalPending() int {
	cont := 0

	for _, channel := range r.Shards {
		cont += len(channel)
	}

	return cont
}
