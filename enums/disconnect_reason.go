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

func (reason DisconnectReason) String() string {
	switch reason {
	case CloseReasonNormalDisconnection:
		return "CloseReasonNormalDisconnection"

	case CloseReasonUnknown:
		return "CloseReasonUnknown"

	case CloseReasonChallengeTimeout:
		return "CloseReasonChallengeTimeout"

	case CloseReasonIdleTimeout:
		return "CloseReasonIdleTimeout"

	case CloseReasonClosedByGateway:
		return "CloseReasonClosedByGateway"

	case CloseReasonClosedByAnotherService:
		return "CloseReasonClosedByAnotherService"

	case CloseReasonInvalidPacket:
		return "CloseReasonInvalidPacket"

	case CloseReasonProtocolViolation:
		return "CloseReasonProtocolViolation"

	case CloseReasonGatewayInternalError:
		return "CloseReasonGatewayInternalError"

	case CloseReasonFailedAuthentication:
		return "CloseReasonFailedAuthentication"

	case CloseReasonSessionIdTakenByOtherConnection:
		return "CloseReasonSessionIdTakenByOtherConnection"

	case CloseReasonConnectionUnauthorized:
		return "CloseReasonConnectionUnauthorized"

	case CloseReasonUserRedirected:
		return "CloseReasonUserRedirected"

	case CloseReasonClosedByAdmin:
		return "CloseReasonClosedByAdmin"

	case CloseReasonConnectionLost:
		return "CloseReasonConnectionLost"

	case CloseReasonClientInternalError:
		return "CloseReasonClientInternalError"
	}

	return "shit"
}
