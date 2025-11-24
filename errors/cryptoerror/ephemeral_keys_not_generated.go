package cryptoerror

import (
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/errorinfo"
)

type EphemeralKeysNotGenerated struct {
	errors.BaseError
}

func (err *EphemeralKeysNotGenerated) Error() string {
	return "Error whith ephemeral keys, they are not generated yet"
}

func (err *EphemeralKeysNotGenerated) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return err
}

func CreateEphemeralKeysNotGeneratedError(file string, index uint16) *EphemeralKeysNotGenerated {
	return &EphemeralKeysNotGenerated{
		errors.CreateBaseError(file, "error with ephemeral keys not generated", index, 21),
	}
}
