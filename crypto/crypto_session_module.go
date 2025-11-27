package crypto

import (
	"crypto/sha256"
	"expansion-gateway/dto/cryptodto"
	"expansion-gateway/enums"
	"expansion-gateway/errors/cryptoerror"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
	"io"
	"sync/atomic"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/hkdf"
)

// the module of the sessions in Layer 2 responsible for handling cryptographics
type CryptoSessionModule struct {
	encryptionMethod atomic.Int32 // the encryption method currently used
	key              []byte       // the key used to encrypt content
	ephemeralKeys    atomic.Pointer[cryptodto.EphemeralKeysDto]
}

// returns a new encryption module for the layer 2 sessions
func CreateNewCryptoSessionModule() *CryptoSessionModule {
	answer := CryptoSessionModule{
		encryptionMethod: atomic.Int32{},
		key:              []byte{},
		ephemeralKeys:    atomic.Pointer[cryptodto.EphemeralKeysDto]{},
	}

	answer.encryptionMethod.Store(int32(enums.NoEncryptionAlgorithm))

	return &answer
}

// ===== Encryption method =====

// sets the encryption method used in this session
func (module *CryptoSessionModule) SetEncryptionAlgorithm(method enums.EncryptionAlgorithm) {
	module.encryptionMethod.Store(int32(method))
}

// gets the encryption method used in this session
func (module *CryptoSessionModule) GetEncryptionAlgorithm() enums.EncryptionAlgorithm {
	return enums.EncryptionAlgorithm(module.encryptionMethod.Load())
}

// ===== Encryption Key =====

// gets the current key stored
func (module *CryptoSessionModule) GetKey() []byte {
	return module.key
}

// generate a new key or password and stores it
func (module *CryptoSessionModule) GenerateNewKey(clientEphemeralPubKey []byte) errorinfo.GatewayError {
	ephemeralKeys := module.ephemeralKeys.Load()
	const filePath string = "/crypto/crypto_session_module.go"

	if ephemeralKeys == nil {
		return cryptoerror.CreateEphemeralKeysNotGeneratedError(filePath, 62)
	}

	serverPrivateKey := ephemeralKeys.GetPrivateKey()

	if sharedSecret, err := curve25519.X25519(serverPrivateKey[:], clientEphemeralPubKey); err == nil {
		// --- Derive final session key via HKDF-SHA256 ---
		keyName := fmt.Sprintf(
			"key-identifier-%d.%d.%d",
			helpers.GenerateRandomInt64(),
			helpers.GenerateRandomInt64(),
			helpers.GenerateRandomInt64())
		h := hkdf.New(sha256.New, sharedSecret, nil, []byte(keyName))

		sessionKey := make([]byte, 32)

		if _, err2 := io.ReadFull(h, sessionKey); err2 != nil {
			return cryptoerror.CreateHKDFfailedError(filePath, 75)
		}

		// there you have it, the final key, securely generated
		module.key = sessionKey
	}

	return nil
}

// ===== Ephemeral Keys =====

// gets the current ephemeral key pair
func (module *CryptoSessionModule) GetEphemeralKeys() *cryptodto.EphemeralKeysDto {
	return module.ephemeralKeys.Load()
}

// deletes the current ephemeral key pair
func (module *CryptoSessionModule) DeleteEphemeralKeys() {
	if keys := module.ephemeralKeys.Load(); keys != nil {
		keys.ErasePrivateKey()
		module.ephemeralKeys.Store(nil)
	}
}

// returns true if there are any ephemeral key pair stored
func (module *CryptoSessionModule) HasEphemeralKeys() bool {
	key := module.ephemeralKeys.Load()

	return key != nil
}

// generates new ephemeral keys
func (module *CryptoSessionModule) GenerateEphemeralKeys() errorinfo.GatewayError {
	if priv, pub, err := GenerateX25519KeyPair(); err != nil {
		return err
	} else {
		module.ephemeralKeys.Store(cryptodto.GenerateNewEphemeralKeysDto(priv, pub))
	}

	return nil
}

// sets the ephemeral keys from an external source
func (module *CryptoSessionModule) SetEphemeralKeys(keys *cryptodto.EphemeralKeysDto) {
	module.ephemeralKeys.Store(keys)
}

// ===== Encrypt / Decrypt

// encrypt a string using the values stored in the module
func (module *CryptoSessionModule) Encrypt(textToEncrypt *string) []byte {
	return []byte{}
}

// decrypts a string using the values stored in the module
func (module *CryptoSessionModule) Decrypt(streamToDecrypt []byte) string {
	return ""
}

// ===== Connected data =====

// this function returns the data that needs to be sent for the encryption section
// in a connected packet
func (module *CryptoSessionModule) GetDataForConnectedPacket() []byte {
	answer := make([]byte, 0, 200)

	switch enums.EncryptionAlgorithm(module.encryptionMethod.Load()) {
	case enums.NoEncryptionAlgorithm: // do nothing
		return []byte{}

	case enums.XChaCha20:
		answer = append(answer, 25)
		answer = append(answer, byte(enums.XChaCha20))

		// continue here...
	}

	return answer
}
