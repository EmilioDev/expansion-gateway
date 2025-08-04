package dispatchers

import "expansion-gateway/interfaces/packets"

type Reciver interface {
	GetShard(index int) <-chan packets.Packet
	ShardCount() int
}
