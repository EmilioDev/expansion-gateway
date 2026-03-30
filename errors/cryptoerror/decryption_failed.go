package cryptoerror

import (
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/errorinfo"
)

type DecryptionFailed struct {
	errors.BaseError
	innerError error
}

func (err *DecryptionFailed) Error() string {
	return "Error whith ephemeral keys, they are not generated yet"
}

func (err *DecryptionFailed) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return err
}

func CreateDecryptionFailedError(file string, index uint16, innerError error) *DecryptionFailed {
	return &DecryptionFailed{
		BaseError:  errors.CreateBaseError(file, "crypto engine not generated", index, 30),
		innerError: innerError,
	}
}
