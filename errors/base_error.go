package errors

import "fmt"

type BaseError struct {
	File      string
	Line      uint16
	Reason    string
	ErrorCode byte
}

func CreateBaseError(file, reason string, line uint16, errorCode byte) BaseError {
	return BaseError{
		File:      file,
		Line:      line,
		Reason:    reason,
		ErrorCode: errorCode,
	}
}

// Message of the error
func (err *BaseError) Error() string {
	return fmt.Sprintf("Error in file %s, line %d, reason: %s, error code %d", err.File, err.Line, err.Reason, err.ErrorCode)
}

func (err BaseError) GetErrorCode() byte {
	return err.ErrorCode
}
