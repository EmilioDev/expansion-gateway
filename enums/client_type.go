package enums

type ClientType byte

const (
	GODOT_CLIENT ClientType = iota
	CLI_TOOL
	maxClientType
)

func IsValidClientType(candidateToClientType byte) bool {
	return candidateToClientType < byte(maxClientType)
}
