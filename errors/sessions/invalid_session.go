package sessions

import (
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type InvalidSession struct {
	*SessionError
}

func (err *InvalidSession) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return err
}

func (err *InvalidSession) Error() string {
	return fmt.Sprintf("Error with session %d", err.sessionID)
}

func CreateInvalidSessionError(file string, index uint16, sessionId int64) *InvalidSession {
	return &InvalidSession{
		SessionError: CreateNewSessionError(file, "Session does not exist or is invalid", index, 24, sessionId),
	}
}
