package enums

type PacketType byte

const (
	HELLO PacketType = iota
	CHALLENGE
	CONNECT
	CONNECTED
	REDIRECT
	DISCONNECT
	NONE
)

func GetNameOfPacketType(packetType PacketType) string {
	switch packetType {
	case HELLO:
		return "HELLO"

	case CHALLENGE:
		return "CHALLENGE"

	case CONNECT:
		return "CONNECT"

	case CONNECTED:
		return "CONNECTED"

	case REDIRECT:
		return "REDIRECT"

	case DISCONNECT:
		return "DISCONNECT"

	default:
		return ""
	}
}
