package enums

type ProtocolVersion byte

const (
	V1 ProtocolVersion = iota
)

func IsValidProtocolVersion(candidateToProtocolVersion byte) bool {
	return candidateToProtocolVersion == 0
}
