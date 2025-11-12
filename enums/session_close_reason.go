// file: /enums/session_close_reason.go
package enums

type SessionCloseReason byte

const (
	CloseReasonUnknown SessionCloseReason = iota
	CloseReasonChallengeTimeout
	CloseReasonIdleTimeout
	CloseReasonManual
	CloseReasonInvalidPacket
	CloseReasonProtocolViolation
	CloseReasonInternalError
	CloseReasonFailedAuthentication
	CloseReasonSessionIdTakenByOtherConnection
	CloseReasonConnectionUnauthorized
	CloseReasonUserRedirected
)
