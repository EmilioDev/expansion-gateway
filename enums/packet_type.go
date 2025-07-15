package enums

type PacketType byte

const (
	HELLO PacketType = iota
	CHALLENGE
	CONNECT
	CONNECTED
	REDIRECT
	DISCONNECT
)
