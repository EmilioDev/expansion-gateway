package dispatchers

import "expansion-gateway/interfaces/packets"

type Dispatcher interface {
	Dispatch(pkt packets.Packet)
}
