package crypto

import "expansion-gateway/enums"

type EncryptedMessage struct {
	EncryptionAlgorythm enums.EncryptionAlgorithm // the algorithm used to generate this message
	CipherText          []byte                    // the text, but ciphered
	Nonce               []byte                    // the nonce used in the encryption
	AAD                 []byte                    // any aditional authenticated data (aad) used in the encryption
}
