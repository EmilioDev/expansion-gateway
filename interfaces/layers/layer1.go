package layers

import "expansion-gateway/interfaces"

type Layer1 interface {
	interfaces.OneWayPipe
	Layer
}
