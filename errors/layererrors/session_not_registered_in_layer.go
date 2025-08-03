package layererrors

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type SessionNotRegisteredInLayer struct {
	LayerError
	SessionId int64
}

func (err SessionNotRegisteredInLayer) Error() string {
	return fmt.Sprintf("Layer of kind %s has a request to session %d, but this session is not registered", enums.LayerKindToString(err.SourceLayer), err.SessionId)
}

func (err SessionNotRegisteredInLayer) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}

func CreateSessionNotRegistered_LayerError(file string, index uint16, layerType enums.LayerKind, sessionId int64) SessionNotRegisteredInLayer {
	return SessionNotRegisteredInLayer{
		CreateLayerError(file, "layer closed", index, 9, layerType),
		sessionId,
	}
}
