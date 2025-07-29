package layererrors

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type ChannelClosedInLayerError struct {
	LayerError
	ChannelClosedKind enums.LayerChannelType
}

func (err ChannelClosedInLayerError) Error() string {
	return fmt.Sprintf("Layer of kind %s had issues with %s", enums.LayerKindToString(err.SourceLayer), enums.LayerChannelTypeToString(err.ChannelClosedKind))
}

func (err ChannelClosedInLayerError) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}

func CreateChannelClosed_LayerError(file string, index uint16, layerType enums.LayerKind, kindOchannelClosedError enums.LayerChannelType) ChannelClosedInLayerError {
	return ChannelClosedInLayerError{
		CreateLayerError(file, "layer channel closed", index, 8, layerType),
		kindOchannelClosedError,
	}
}
