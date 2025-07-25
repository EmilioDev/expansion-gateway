package errorinfo

type GatewayError interface {
	error
	GetErrorCode() byte
	GetStackTrace() []string
}
