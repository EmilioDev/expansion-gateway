package layers

import "expansion-gateway/interfaces"

type Layer2 interface {
	Layer
	interfaces.TwoWaysPipe
}
