// file: /enums/flow_violation_kind.go
package enums

type FlowViolationKind byte

const (
	INVALID_HELLO FlowViolationKind = iota
	CLIENT_SENT_CHALLENGE
	INVALID_PACKET
)
