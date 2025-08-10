// file: /enums/connection_state.go
package enums

type ConnectionState byte

const (
	HELLO_RECEIVED ConnectionState = iota
	CHALLENGE_SENT
	RECEIVED_CONNECT
	SESSION_CONNECTED
	REDIRECTING
	DISCONNECTING
)
