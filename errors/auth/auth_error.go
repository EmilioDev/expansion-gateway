package auth

import (
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type AuthError struct {
	errors.BaseError
	ConnectionID int64
}

func (err AuthError) Error() string {
	return fmt.Sprintf(" %d) is not authorized", err.ConnectionID)
}

func (err AuthError) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}

func CreateAuthError(file, description string, index uint16, errorCode byte, connectionID int64) AuthError {
	return AuthError{
		BaseError:    errors.CreateBaseError(file, description, index, errorCode),
		ConnectionID: connectionID,
	}
}
