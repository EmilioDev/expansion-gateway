package cryptoerror

import (
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/errorinfo"
)

type CryptoEngineNotGenerated struct {
	errors.BaseError
	innerError error
}

func (err *CryptoEngineNotGenerated) Error() string {
	return "Error whith ephemeral keys, they are not generated yet"
}

func (err *CryptoEngineNotGenerated) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return err
}

func CreateCryptoEngineNotGeneratedError(file string, index uint16, innerError error) *CryptoEngineNotGenerated {
	return &CryptoEngineNotGenerated{
		BaseError:  errors.CreateBaseError(file, "crypto engine not generated", index, 28),
		innerError: innerError,
	}
}
