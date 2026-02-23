package nats

import (
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type NatsError struct {
	errors.BaseError
	innerError error
}

func (err *NatsError) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return err
}

func (err *NatsError) Error() string {
	if err.innerError != nil {
		return fmt.Sprintf("Error: %s", err.innerError.Error())
	}

	return "NATS error"
}

func CreateNatsError(file, description string, index uint16, errorCode byte, innerError error) *NatsError {
	return &NatsError{
		BaseError:  errors.CreateBaseError(file, description, index, errorCode),
		innerError: innerError,
	}
}
