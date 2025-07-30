package dispatchers

import (
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/packets"
	"hash/fnv"
)

type PacketDispatcher struct {
	shards []chan packets.Packet
	count  int
}

func (d *PacketDispatcher) Dispatch(pkt packets.Packet) {
	index := hashPacket(pkt) % uint32(d.count)
	d.shards[index] <- pkt
}

func hashPacket(pkt packets.Packet) uint32 {
	h := fnv.New32a()
	id := helpers.ConvertInt64Into8Bytes(pkt.GetSender())

	h.Write(id[:]) // safe since sender is int64

	return h.Sum32()
}
