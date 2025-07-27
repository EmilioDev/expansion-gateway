package layers

import (
	"expansion-gateway/interfaces/pipes"
)

type Layer3 interface {
	Layer
	pipes.OneWayPipe
}
