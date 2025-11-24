package enums

type EncryptionAlgorithm byte

const (
	XChaCha20 EncryptionAlgorithm = iota
	AES_GCM
	AES_CTR_PLUS_HMAC_SHA256
	NoEncryptionAlgorithm
)

func IsValidEncryptionAlgorythm(candidate byte) bool {
	return candidate <= byte(NoEncryptionAlgorithm)
}
