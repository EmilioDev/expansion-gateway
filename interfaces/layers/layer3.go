package layers

import "expansion-gateway/interfaces"

type Layer3 interface {
	Layer
	interfaces.OneWayPipe
}
