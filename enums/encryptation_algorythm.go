package enums

type EncryptionAlgorithm byte

const (
	Salsa20 EncryptionAlgorithm = iota
	// new values will be added between here
	// and here
	NoEncryptionAlgorithm
)

func IsValidEncryptionAlgorythm(candidate byte) bool {
	return candidate < byte(NoEncryptionAlgorithm)
}
