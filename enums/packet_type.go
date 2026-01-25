package enums

type PacketType byte

const (
	HELLO PacketType = iota
	CHALLENGE
	CONNECT
	CONNECTED
	REDIRECT
	REDIRECTED
	SUBSCRIBE
	SUBACK
	UNSUBSCRIBE
	UNSUBACK
	PUBLISH
	PUBACK
	PING
	PINGACK
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

	case REDIRECTED:
		return "REDIRECTED"

	case SUBSCRIBE:
		return "SUBS"

	case SUBACK:
		return "SUBACK"

	case PUBLISH:
		return "PUBLISH"

	case PUBACK:
		return "PUBACK"

	case PING:
		return "PING"

	case PINGACK:
		return "PINGACK"

	case DISCONNECT:
		return "DISCONNECT"

	default:
		return ""
	}
}
