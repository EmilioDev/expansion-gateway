package errorinfo

type DebugableError interface {
	SetStackTrace(stackTrace []string) GatewayError
}
