// file: /interfaces/layers/smart_layer.go
package layers

import "expansion-gateway/interfaces/errorinfo"

type SmartLayer interface {
	Layer
	// method for configuring layer 1
	ConfigureFirstLayer(layer Layer1) errorinfo.GatewayError

	// method for configuring the third layer
	ConfigureThirdLayer(layer Layer3) errorinfo.GatewayError

	// wait untill all actions have been completed
	Wait()
}
