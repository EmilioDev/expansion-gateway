// file: /interfaces/layers/smart_layer.go
package layers

import "expansion-gateway/interfaces/errorinfo"

type SmartLayer interface {
	Layer
	// method for configuring layer 1
	ConfigureFirstLayer(layer Layer1) errorinfo.GatewayError
	ConfigureThirdLayer(layer Layer3) errorinfo.GatewayError
}
