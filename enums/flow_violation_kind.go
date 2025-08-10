package enums

type FlowViolationKind byte

const (
	INVALID_HELLO FlowViolationKind = iota
	CLIENT_SENT_CHALLENGE
)
