package engines

import (
	"crypto/cipher"

	"expansion-gateway/enums"
	"expansion-gateway/errors/cryptoerror"
	"expansion-gateway/interfaces/errorinfo"

	"golang.org/x/crypto/chacha20poly1305"
)

type Chacha20CryptoEngine struct {
	aead         cipher.AEAD
	connectionID int64
}

func NewChacha20CryptoEngine(shared [32]byte, connectionID int64) (*Chacha20CryptoEngine, errorinfo.GatewayError) {
	if aead, err := chacha20poly1305.NewX(shared[:]); err == nil {
		return &Chacha20CryptoEngine{
			aead:         aead,
			connectionID: connectionID,
		}, nil
	} else {
		return nil, cryptoerror.CreateCryptoEngineNotGeneratedError(
			"/crypto/engines/chacha20_crypto_engine.go",
			17,
			err,
		)
	}
}

func (engine *Chacha20CryptoEngine) Encrypt(counter uint64, data []byte) ([]byte, errorinfo.GatewayError) {
	nonce := buildNonce(engine.connectionID, counter, chacha20poly1305.NonceSizeX)
	return engine.aead.Seal(nil, nonce, data, nil), nil
}

func (engine *Chacha20CryptoEngine) Decrypt(counter uint64, data []byte) ([]byte, errorinfo.GatewayError) {
	nonce := buildNonce(engine.connectionID, counter, chacha20poly1305.NonceSizeX)

	if res, err := engine.aead.Open(nil, nonce, data, nil); err == nil {
		return res, nil
	} else {
		return nil, cryptoerror.CreateDecryptionFailedError(
			"/crypto/engines/chacha20_crypto_engine.go",
			41,
			err,
		)
	}
}

func (engine *Chacha20CryptoEngine) EncryptionSupported() enums.EncryptionAlgorithm {
	return enums.XChaCha20
}
