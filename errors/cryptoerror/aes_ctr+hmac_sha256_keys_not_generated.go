package cryptoerror

import (
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/errorinfo"
)

type AesCtrHmacSha256KeysNotGenerated struct {
	errors.BaseError
}

func (err *AesCtrHmacSha256KeysNotGenerated) Error() string {
	return "AES-CTR + HMAC-SHA256 key generation failed"
}

func (err *AesCtrHmacSha256KeysNotGenerated) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return err
}

func CreateAesCtrHmacSha256KeysNotGenerated(file string, index uint16) *AesCtrHmacSha256KeysNotGenerated {
	return &AesCtrHmacSha256KeysNotGenerated{
		errors.CreateBaseError(file, "key generation failed", index, 31),
	}
}
