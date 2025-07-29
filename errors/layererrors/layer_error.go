package layererrors

import (
	"expansion-gateway/enums"
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type LayerError struct {
	errors.BaseError
	SourceLayer enums.LayerKind
}

func (err LayerError) Error() string {
	return fmt.Sprintf("Layer %d thrown an error", byte(err.SourceLayer))
}

func CreateLayerError(file, reason string, index uint16, errorCode byte, layerType enums.LayerKind) LayerError {
	return LayerError{
		errors.CreateBaseError(file, reason, index, errorCode),
		layerType,
	}
}

func (err LayerError) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}
