package auth

import (
	"expansion-gateway/interfaces/errorinfo"
)

type ConnectWithMissingFields struct {
	AuthError
	MissingFieldName string
}

func (err *ConnectWithMissingFields) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return err
}

func GenerateConnectWithEphemeralKeyMissing(file string, index uint16, connectionId int64) *ConnectWithMissingFields {
	return &ConnectWithMissingFields{
		AuthError: CreateAuthError(
			file,
			"connect packet with client ephemeral key missing",
			index,
			20,
			connectionId),
		MissingFieldName: "ClientEphemeralKey",
	}
}
