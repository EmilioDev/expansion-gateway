package helpers

import "expansion-gateway/interfaces/errorinfo"

func WithStackTrace(e errorinfo.DebugableError, skip int) errorinfo.GatewayError {
	return e.SetStackTrace(GetStackTrace(skip + 1)) // +1 to skip helper itself
}
