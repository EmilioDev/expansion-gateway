package enums

type EncryptationAlgorythm byte

const (
	NoEncryptation EncryptationAlgorythm = iota
	Salsa20
	//
	maxEncryptationAlgorythm
)

func IsValidEncryptationAlgorythm(candidate byte) bool {
	return candidate < byte(maxEncryptationAlgorythm)
}
