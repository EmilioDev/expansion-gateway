package layererrors

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type LayerClosedError struct {
	LayerError
}

func (err LayerClosedError) Error() string {
	return fmt.Sprintf("Layer of kind %s has a request, but this layer is closed", enums.LayerKindToString(err.SourceLayer))
}

func (err LayerClosedError) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}

func CreateLayerClosed_LayerError(file string, index uint16, layerType enums.LayerKind) LayerClosedError {
	return LayerClosedError{
		CreateLayerError(file, "layer closed", index, 9, layerType),
	}
}
