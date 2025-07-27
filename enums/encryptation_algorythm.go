package enums

type EncryptionAlgorithm byte

const (
	NoEncryption EncryptionAlgorithm = iota
	Salsa20
	//new values will be added here...
	MaxEncryptionAlgorythm
)

func IsValidEncryptionAlgorythm(candidate byte) bool {
	return candidate < byte(MaxEncryptionAlgorythm)
}
