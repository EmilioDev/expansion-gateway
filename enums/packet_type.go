package enums

type PacketType byte

const (
	HELLO PacketType = iota
	CHALLENGE
	CONNECT
	CONNECTED
	REDIRECT
	REDIRECTED
	REDIRECTED_OK
	SUBS
	PUBLISH
	PING
	ACK
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

	case REDIRECTED_OK:
		return "REDIRECTED_OK"

	case SUBS:
		return "SUBS"

	case PUBLISH:
		return "PUBLISH"

	case PING:
		return "PING"

	case ACK:
		return "ACK"

	case DISCONNECT:
		return "DISCONNECT"

	default:
		return ""
	}
}
