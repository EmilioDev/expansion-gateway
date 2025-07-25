package enums

type EncryptionAlgorithm byte

const (
	NoEncryption EncryptionAlgorithm = iota
	Salsa20
	//
	maxEncryptionAlgorythm
)

func IsValidEncryptionAlgorythm(candidate byte) bool {
	return candidate < byte(maxEncryptionAlgorythm)
}
