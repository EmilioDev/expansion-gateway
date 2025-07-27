package layers

import (
	"expansion-gateway/interfaces/pipes"
)

type Layer2 interface {
	Layer
	pipes.TwoWaysPipe
}
