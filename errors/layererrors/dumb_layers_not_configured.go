package layererrors

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/layers"
)

type DumbLayersNotConfigured struct {
	LayerError
	Layer1 layers.Layer1
	Layer3 layers.Layer3
}

func CreateDumbLayersNotConfigured_LayerError(file string, index uint16, layerType enums.LayerKind, layer1 layers.Layer1, layer3 layers.Layer3) DumbLayersNotConfigured {
	return DumbLayersNotConfigured{
		CreateLayerError(file, "layer channel closed", index, 11, layerType),
		layer1,
		layer3,
	}
}

func (err DumbLayersNotConfigured) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}
