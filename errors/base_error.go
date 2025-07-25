package errors

import (
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type BaseError struct {
	File       string   // File where this error had happen
	Index      uint16   // Line where this error had happen, now it is an index
	Reason     string   // The reason of the error
	ErrorCode  byte     // An error code, for an easy identification of the error type without needing to check if this error is of one type or another
	StackTrace []string // The stackTrace of this error
}

func CreateBaseError(file, reason string, index uint16, errorCode byte) BaseError {
	return BaseError{
		File:       file,
		Index:      index,
		Reason:     reason,
		ErrorCode:  errorCode,
		StackTrace: []string{},
	}
}

// Message of the error
func (err BaseError) Error() string {
	return fmt.Sprintf("Error in file %s, line %d, reason: %s, error code %d", err.File, err.Index, err.Reason, err.ErrorCode)
}

func (err BaseError) GetErrorCode() byte {
	return err.ErrorCode
}

func (err BaseError) GetStackTrace() []string {
	return err.StackTrace
}

func (err BaseError) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}
