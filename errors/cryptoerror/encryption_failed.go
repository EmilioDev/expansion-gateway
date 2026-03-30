package cryptoerror

import (
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/errorinfo"
)

type EncryptionFailed struct {
	errors.BaseError
	innerError error
}

func (err *EncryptionFailed) Error() string {
	return "Error whith ephemeral keys, they are not generated yet"
}

func (err *EncryptionFailed) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return err
}

func CreateEncryptionFailedError(file string, index uint16, innerError error) *EncryptionFailed {
	return &EncryptionFailed{
		BaseError:  errors.CreateBaseError(file, "crypto engine not generated", index, 29),
		innerError: innerError,
	}
}
