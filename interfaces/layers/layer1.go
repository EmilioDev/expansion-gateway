package layers

import (
	"expansion-gateway/interfaces/pipes"
)

type Layer1 interface {
	pipes.OneWayPipe
	Layer
}
