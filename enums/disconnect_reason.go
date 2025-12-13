package enums

type DisconnectReason byte

const (
	NORMAL_DISCONNECTION DisconnectReason = iota
	CLIENT_DISCONNECTED
)
