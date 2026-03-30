package engines

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
)

type NoneCryptoEngine struct{}

func (engine *NoneCryptoEngine) Encrypt(counter uint64, data []byte) ([]byte, errorinfo.GatewayError) {
	return data, nil
}

func (engine *NoneCryptoEngine) Decrypt(counter uint64, data []byte) ([]byte, errorinfo.GatewayError) {
	return data, nil
}

func (engine *NoneCryptoEngine) EncryptionSupported() enums.EncryptionAlgorithm {
	return enums.NoEncryptionAlgorithm
}

func NewNoneCryptoEngine() *NoneCryptoEngine {
	return &NoneCryptoEngine{}
}
