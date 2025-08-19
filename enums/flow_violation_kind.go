// file: /enums/flow_violation_kind.go
package enums

type FlowViolationKind byte

const (
	INVALID_HELLO FlowViolationKind = iota
	CLIENT_SENT_CHALLENGE
	INVALID_PACKET
	INCORRECT_PACKET_KIND
	CLIENT_SENT_CONNECT_AT_WRONG_MOMENT
	SESSION_CLOSED
)
