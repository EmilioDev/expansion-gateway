package dispatchers

import "expansion-gateway/interfaces/packets"

type Reciver interface {
	GetShard(index int) <-chan packets.Packet // returns the shard at the given index
	ShardCount() int                          // number of shards this receiver has
	TotalPending() int                        // number of messages still pending to be processed
}
