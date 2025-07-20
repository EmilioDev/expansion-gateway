package enums

type EncryptationAlgorythm byte

const (
	NoEncryptation EncryptationAlgorythm = iota
	Salsa20
)
