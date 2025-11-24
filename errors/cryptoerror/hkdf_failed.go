package cryptoerror

import (
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/errorinfo"
)

type HkdfFailed struct {
	errors.BaseError
}

func (err *HkdfFailed) Error() string {
	return "HMAC Key Derivation Function (HKDF) failed"
}

func (err *HkdfFailed) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return err
}

func CreateHKDFfailedError(file string, index uint16) *HkdfFailed {
	return &HkdfFailed{
		errors.CreateBaseError(file, "HKDF failed", index, 22),
	}
}
