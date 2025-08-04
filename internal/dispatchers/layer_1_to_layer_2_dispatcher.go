package dispatchers

import (
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/packets"
	"hash/fnv"
)

type PacketDispatcher struct {
	Shards []chan packets.Packet
	Count  int
}

func (d *PacketDispatcher) Dispatch(pkt packets.Packet) {
	index := hashPacket(pkt) % uint32(d.Count)
	d.Shards[index] <- pkt
}

func hashPacket(pkt packets.Packet) uint32 {
	h := fnv.New32a()
	id := helpers.ConvertInt64Into8Bytes(pkt.GetSender())

	h.Write(id[:]) // safe since sender is int64

	return h.Sum32()
}
