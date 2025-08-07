// file: /interfaces/layers/layer.go
package layers

import "expansion-gateway/interfaces/errorinfo"

// Base type of all the three layers of the gateway
type Layer interface {
	// Sends data to be processed by this layer
	Start() errorinfo.GatewayError
	Stop() errorinfo.GatewayError
	IsWorking() bool
}
