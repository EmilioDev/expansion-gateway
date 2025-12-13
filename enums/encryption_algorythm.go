package enums

type EncryptionAlgorithm byte

const (
	NoEncryptionAlgorithm EncryptionAlgorithm = iota
	XChaCha20
	AES_GCM
	AES_CTR_PLUS_HMAC_SHA256
)

func IsValidEncryptionAlgorythm(candidate byte) bool {
	return candidate <= byte(AES_CTR_PLUS_HMAC_SHA256)
}
