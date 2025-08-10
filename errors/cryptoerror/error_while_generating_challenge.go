package cryptoerror

import (
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/errorinfo"
)

type ErrorWhileGeneratingChallenge struct {
	errors.BaseError
	InnerError error
}

func (err ErrorWhileGeneratingChallenge) Error() string {
	return "Error while generating challenge nonce for checking a client authenticity"
}

func (err ErrorWhileGeneratingChallenge) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}

func CreateErrorWhileGeneratingChallenge(file string, index uint16, theError error) ErrorWhileGeneratingChallenge {
	return ErrorWhileGeneratingChallenge{
		errors.CreateBaseError(file, "error generating the challenge nonce", index, 12),
		theError,
	}
}
