package errors

import (
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type ExternalErrorWrapper struct {
	BaseError
	InnerError error
}

// Message of the error
func (err ExternalErrorWrapper) Error() string {
	return fmt.Sprintf("Error in file %s, line %d, reason: %s, error code %d, inner error: %s", err.File, err.Index, err.Reason, err.ErrorCode, err.InnerError.Error())
}

// Creates a new error wrapper
func CreateErrorWrapper(file string, line uint16, innerError error) ExternalErrorWrapper {
	return ExternalErrorWrapper{
		CreateBaseError(file, "external error", line, 7),
		innerError,
	}
}

func (err ExternalErrorWrapper) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}
