package engines

import (
	"crypto/aes"
	"crypto/cipher"

	"expansion-gateway/enums"
	"expansion-gateway/errors/cryptoerror"
	"expansion-gateway/interfaces/errorinfo"
)

type AESGCMCryptoEngine struct {
	aead         cipher.AEAD
	connectionID int64
}

func NewAESGCMCipher(sharedSecret [32]byte, connectionID int64) (*AESGCMCryptoEngine, errorinfo.GatewayError) {
	block, err := aes.NewCipher(sharedSecret[:])
	const filePath string = "/crypto/engines/aes_gcm_crypto_engine.go"

	if err != nil {
		return nil, cryptoerror.CreateCryptoEngineNotGeneratedError(
			filePath,
			17,
			err,
		)
	}

	aead, err := cipher.NewGCM(block)

	if err != nil {
		return nil, cryptoerror.CreateCryptoEngineNotGeneratedError(
			filePath,
			28,
			err,
		)
	}

	return &AESGCMCryptoEngine{
		aead:         aead,
		connectionID: connectionID,
	}, nil
}

func (engine *AESGCMCryptoEngine) Encrypt(counter uint64, data []byte) ([]byte, errorinfo.GatewayError) {
	nonce := buildNonce(engine.connectionID, counter, engine.aead.NonceSize())
	return engine.aead.Seal(nil, nonce, data, nil), nil
}

func (engine *AESGCMCryptoEngine) Decrypt(counter uint64, data []byte) ([]byte, errorinfo.GatewayError) {
	nonce := buildNonce(engine.connectionID, counter, engine.aead.NonceSize())

	if res, err := engine.aead.Open(nil, nonce, data, nil); err == nil {
		return res, nil
	} else {
		return nil, cryptoerror.CreateDecryptionFailedError(
			"/crypto/engines/aes_gcm_crypto_engine.go",
			52,
			err,
		)
	}
}

func (engine *AESGCMCryptoEngine) EncryptionSupported() enums.EncryptionAlgorithm {
	return enums.AES_GCM
}
