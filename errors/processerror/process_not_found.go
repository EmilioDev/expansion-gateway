package processerror

import (
	"expansion-gateway/errors"
	"fmt"
)

type ProcessNotFoundError struct {
	errors.BaseError
	pid int32
}

func (err ProcessNotFoundError) Error() string {
	return fmt.Sprintf("Process %d was not found", err.pid)
}

func CreateProcessNotFoundError(file string, index uint16, pid int32) *ProcessNotFoundError {
	return &ProcessNotFoundError{
		BaseError: errors.CreateBaseError(file, "process not found", index, 18),
		pid:       pid,
	}
}
