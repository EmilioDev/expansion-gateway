package cryptodto

import (
	"expansion-gateway/helpers"
	"sync"
)

// container of the ephemeral keys
type EphemeralKeysDto struct {
	private        [32]byte     // the private key
	public         [32]byte     // the public key
	privateKeySync sync.RWMutex // mutex for accesing the private key
}

// returns the public key, this key can be sent to the client
func (dto *EphemeralKeysDto) GetPublicKey() [32]byte {
	return dto.public
}

// returns the private key, use only to generate the final key, do not share it
func (dto *EphemeralKeysDto) GetPrivateKey() [32]byte {
	dto.privateKeySync.RLock()
	defer dto.privateKeySync.RUnlock()

	return dto.private
}

// zero the private key
func (dto *EphemeralKeysDto) ErasePrivateKey() {
	dto.privateKeySync.Lock()
	helpers.ZeroKey(&dto.private)
	dto.privateKeySync.Unlock()
}

func GenerateNewEphemeralKeysDto(privateKey, publicKey [32]byte) *EphemeralKeysDto {
	return &EphemeralKeysDto{
		private:        privateKey,
		public:         publicKey,
		privateKeySync: sync.RWMutex{},
	}
}
