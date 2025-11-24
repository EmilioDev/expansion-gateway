package auth

import "expansion-gateway/interfaces/errorinfo"

type ConnectionUnauthorizedError struct {
	AuthError
}

func (err *ConnectionUnauthorizedError) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return err
}

func CreateConnectionUnauthorizedError(file string, index uint16, connectionID int64) ConnectionUnauthorizedError {
	return ConnectionUnauthorizedError{
		AuthError: CreateAuthError(file, "connection unauthorized", index, 19, connectionID),
	}
}
