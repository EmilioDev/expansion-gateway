package nats

import (
	"expansion-gateway/interfaces/errorinfo"
)

type ConnectionToNatsFailed struct {
	*NatsError
}

func (err *ConnectionToNatsFailed) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return err
}

func CreateConnectionToNatsFailedError(file string, index uint16, innerError error) *ConnectionToNatsFailed {
	return &ConnectionToNatsFailed{
		NatsError: CreateNatsError(file, "connection to NATS failed", index, 26, innerError),
	}
}
