package sessions

import (
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type SessionError struct {
	errors.BaseError
	sessionID int64
}

func (err *SessionError) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return err
}

func (err *SessionError) Error() string {
	return fmt.Sprintf("Error with session %d", err.sessionID)
}

func CreateNewSessionError(file, description string, index uint16, errorCode byte, sessionId int64) *SessionError {
	return &SessionError{
		BaseError: errors.CreateBaseError(file, description, index, errorCode),
		sessionID: sessionId,
	}
}
