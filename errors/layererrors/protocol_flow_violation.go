package layererrors

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
)

type ProtocolFlowViolation struct {
	LayerError
	Kind enums.FlowViolationKind
}

func (err ProtocolFlowViolation) Error() string {
	return "A valid packet has been recived but the sender did not followed the protocol flow"
}

func CreateProtocolFlowViolation_LayerError(file string, index uint16, layerType enums.LayerKind, violation enums.FlowViolationKind) ProtocolFlowViolation {
	return ProtocolFlowViolation{
		CreateLayerError(file, "protocol flow violation", index, 13, layerType),
		violation,
	}
}

func (err ProtocolFlowViolation) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}
