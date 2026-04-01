package engines

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"expansion-gateway/enums"
	"expansion-gateway/errors/cryptoerror"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
	"io"

	"golang.org/x/crypto/hkdf"
)

type AESCTRplusHMACSHA256CryptoEngine struct {
	encKey       []byte
	macKey       []byte
	connectionID int64
}

func NewAESCTRHMACCryptoEngine(shared [32]byte, connectionID int64) (*AESCTRplusHMACSHA256CryptoEngine, errorinfo.GatewayError) {
	keyName := fmt.Sprintf("key-identifier-%d", connectionID)
	const filePath string = "/crypto/engines/aes_gcm_crypto_engine.go"
	hk := hkdf.New(sha256.New, shared[:], nil, []byte(keyName))

	var encKey [32]byte
	var macKey [32]byte

	if _, err := io.ReadFull(hk, encKey[:]); err != nil {
		return nil, cryptoerror.CreateCryptoEngineNotGeneratedError(
			filePath,
			31,
			err,
		)
	}

	if _, err := io.ReadFull(hk, macKey[:]); err != nil {
		return nil, cryptoerror.CreateCryptoEngineNotGeneratedError(
			filePath,
			39,
			err,
		)
	}

	return &AESCTRplusHMACSHA256CryptoEngine{
		encKey:       encKey[:],
		macKey:       macKey[:],
		connectionID: connectionID,
	}, nil
}

func (engine *AESCTRplusHMACSHA256CryptoEngine) Encrypt(counter uint64, data []byte) ([]byte, errorinfo.GatewayError) {
	block, err := aes.NewCipher(engine.encKey)

	if err != nil {
		return nil, cryptoerror.CreateEncryptionFailedError(
			"/crypto/engines/aes_ctr_plus_hmac_sha256_crypto_engine.go",
			55,
			err,
		)
	}

	iv := buildNonce(engine.connectionID, counter, block.BlockSize())

	stream := cipher.NewCTR(block, iv)

	ciphertext := make([]byte, len(data))
	stream.XORKeyStream(ciphertext, data)

	mac := hmac.New(sha256.New, engine.macKey)
	mac.Write(iv)
	mac.Write(ciphertext)
	tag := mac.Sum(nil)

	return append(ciphertext, tag...), nil
}

func (c *AESCTRplusHMACSHA256CryptoEngine) Decrypt(counter uint64, data []byte) ([]byte, errorinfo.GatewayError) {
	block, err := aes.NewCipher(c.encKey)
	const filePath string = "/crypto/engines/aes_ctr_plus_hmac_sha256_crypto_engine.go"

	if err != nil {
		return nil, cryptoerror.CreateDecryptionFailedError(
			filePath,
			81,
			err,
		)
	}

	tagSize := 32
	if len(data) < tagSize {
		return nil, cryptoerror.CreateDecryptionFailedError(
			filePath,
			63,
			fmt.Errorf("invalid data"),
		)
	}

	ciphertext := data[:len(data)-tagSize]
	receivedTag := data[len(data)-tagSize:]

	iv := buildNonce(c.connectionID, counter, block.BlockSize())

	mac := hmac.New(sha256.New, c.macKey)
	mac.Write(iv)
	mac.Write(ciphertext)
	expectedTag := mac.Sum(nil)

	if !hmac.Equal(receivedTag, expectedTag) {
		return nil, cryptoerror.CreateDecryptionFailedError(
			filePath,
			85,
			fmt.Errorf("invalid MAC"),
		)
	}

	stream := cipher.NewCTR(block, iv)

	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

func (engine *AESCTRplusHMACSHA256CryptoEngine) EncryptionSupported() enums.EncryptionAlgorithm {
	return enums.AES_CTR_PLUS_HMAC_SHA256
}
