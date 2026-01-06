// file: /enums/session_close_reason.go
package enums

type DisconnectReason byte

const (
	CloseReasonNormalDisconnection DisconnectReason = iota
	CloseReasonUnknown
	CloseReasonChallengeTimeout
	CloseReasonIdleTimeout
	CloseReasonClosedByGateway
	CloseReasonClosedByAnotherService
	CloseReasonInvalidPacket
	CloseReasonProtocolViolation
	CloseReasonGatewayInternalError
	CloseReasonFailedAuthentication
	CloseReasonSessionIdTakenByOtherConnection
	CloseReasonConnectionUnauthorized
	CloseReasonUserRedirected
	CloseReasonClosedByAdmin
	CloseReasonConnectionLost
	CloseReasonClientInternalError
)

func IsValidDisconnectReason(reason byte) bool {
	return reason <= byte(CloseReasonClientInternalError)
}

func ByteReasonToDisconnectReason(reason byte) DisconnectReason {
	switch reason {
	case 13: // protocol violation
		return CloseReasonProtocolViolation

	case 8, 9, 11, 15, 14, 16, 17, 18: // internal error
		return CloseReasonGatewayInternalError

	case 0, 1, 2, 3, 4, 5, 6, 23, 20: // packet error
		return CloseReasonInvalidPacket

	case 19, 10: //authentication issues
		return CloseReasonConnectionUnauthorized

	case 21, 12, 22:
		return CloseReasonFailedAuthentication

	case 7: // external error
		fallthrough
	default:
		return CloseReasonUnknown
	}
}
