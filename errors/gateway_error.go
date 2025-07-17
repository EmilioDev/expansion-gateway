package errors

type GatewayError interface {
	error
	GetErrorCode() byte
}
