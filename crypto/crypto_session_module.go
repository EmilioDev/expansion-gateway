package crypto

import (
	"expansion-gateway/crypto/engines"
	"expansion-gateway/dto/cryptodto"
	"expansion-gateway/enums"
	"expansion-gateway/errors/cryptoerror"
	"expansion-gateway/interfaces/crypto"
	"expansion-gateway/interfaces/errorinfo"
	"sync"
	"sync/atomic"

	"golang.org/x/crypto/curve25519"
)

// the module of the sessions in Layer 2 responsible for handling cryptographics
type CryptoSessionModule struct {
	encryptionMethod  atomic.Int32                               // the encryption method currently used
	ephemeralKeys     atomic.Pointer[cryptodto.EphemeralKeysDto] // ephemeral keys used for shared secret generation
	cryptoEngine      crypto.CryptoEngine                        // the crypto-engine used for decrypt/encrypt data
	criptoEngineMutex sync.RWMutex                               // cripto engine mutex for avoiding collisions while changing the engine
	counter           atomic.Uint64                              // the counter used to encrypt data
}

// returns a new encryption module for the layer 2 sessions
func CreateNewCryptoSessionModule() *CryptoSessionModule {
	answer := CryptoSessionModule{
		encryptionMethod:  atomic.Int32{},
		ephemeralKeys:     atomic.Pointer[cryptodto.EphemeralKeysDto]{},
		cryptoEngine:      engines.NewNoneCryptoEngine(),
		criptoEngineMutex: sync.RWMutex{},
		counter:           atomic.Uint64{},
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

// generate a new key or password and stores it
func (module *CryptoSessionModule) GenerateNewKey(clientEphemeralPubKey [32]byte, connectionID int64) errorinfo.GatewayError {
	const filePath string = "/crypto/crypto_session_module.go"
	ephemeralKeys := module.ephemeralKeys.Load()

	if ephemeralKeys == nil {
		return cryptoerror.CreateEphemeralKeysNotGeneratedError(filePath, 62)
	}

	privateKey := ephemeralKeys.GetPrivateKey()
	var shared [32]byte

	curve25519.ScalarMult(&shared, &privateKey, &clientEphemeralPubKey)

	encryptionSupported := module.GetEncryptionAlgorithm()

	switch encryptionSupported {
	case enums.AES_CTR_PLUS_HMAC_SHA256:
		module.setAESCtrPlusHmacSha256Engine(connectionID, shared)

	case enums.AES_GCM:
		module.setAESGCMEngine(connectionID, shared)

	case enums.XChaCha20:
		module.setChacha20Engine(connectionID, shared)

	default:
		// do nothing
	}

	return nil
}

func (module *CryptoSessionModule) setAESGCMEngine(connectionID int64, sharedSecret [32]byte) {
	module.criptoEngineMutex.Lock()
	defer module.criptoEngineMutex.Unlock()

	if engine, err := engines.NewAESGCMCipher(sharedSecret, connectionID); err == nil {
		module.cryptoEngine = engine
	}
}

func (module *CryptoSessionModule) setChacha20Engine(connectionID int64, sharedSecret [32]byte) {
	module.criptoEngineMutex.Lock()
	defer module.criptoEngineMutex.Unlock()

	if engine, err := engines.NewChacha20CryptoEngine(sharedSecret, connectionID); err == nil {
		module.cryptoEngine = engine
	}
}

func (module *CryptoSessionModule) setAESCtrPlusHmacSha256Engine(connectionID int64, sharedSecret [32]byte) {
	module.criptoEngineMutex.Lock()
	defer module.criptoEngineMutex.Unlock()

	if engine, err := engines.NewAESCTRHMACCryptoEngine(sharedSecret, connectionID); err == nil {
		module.cryptoEngine = engine
	}
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

// ===== Encrypt / Decrypt =====

// encrypt a string using the values stored in the module
func (module *CryptoSessionModule) Encrypt(textToEncrypt *[]byte) ([]byte, uint64) {
	module.criptoEngineMutex.RLock()
	defer module.criptoEngineMutex.RUnlock()

	counter := module.counter.Load()
	module.counter.Add(1)

	input := *textToEncrypt

	if res, err := module.cryptoEngine.Encrypt(counter, input); err == nil {
		return res, counter
	}

	return input, 0
}

// decrypts a string using the values stored in the module
func (module *CryptoSessionModule) Decrypt(streamToDecrypt *[]byte, counter uint64) ([]byte, errorinfo.GatewayError) {
	module.criptoEngineMutex.RLock()
	defer module.criptoEngineMutex.RUnlock()

	input := *streamToDecrypt

	if res, err := module.cryptoEngine.Decrypt(counter, input); err == nil {
		return res, nil
	} else {
		return nil, err
	}
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
