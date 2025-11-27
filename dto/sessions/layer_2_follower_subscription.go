// file: /dto/sessions/layer_2_follower_subscription.go
package sessions

import (
	"expansion-gateway/crypto"
	"expansion-gateway/dto/cryptodto"
	"expansion-gateway/enums"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/errorinfo"
)

type Layer2FollowerSubscription struct {
	Challenge          []byte
	UserId             int64
	RequestedSessionID int64
	ClientType         enums.ClientType
	ClientVersion      byte
	Encryption         enums.EncryptionAlgorithm
	ProtocolVersion    enums.ProtocolVersion
	SessionResume      bool
	EphemeralKeyPair   *cryptodto.EphemeralKeysDto
}

func (subscription *Layer2FollowerSubscription) GetChallengeAsInt32Array() []int32 {
	result := make([]int32, 0, len(subscription.Challenge))

	for _, val := range subscription.Challenge {
		result = append(result, int32(val))
	}

	return result
}

// Generates a new ephemeral key pair
func (subscription *Layer2FollowerSubscription) GenerateEphemeralKeyPair() errorinfo.GatewayError {
	if priv, pub, err := crypto.GenerateX25519KeyPair(); err == nil {
		subscription.EphemeralKeyPair = cryptodto.GenerateNewEphemeralKeysDto(priv, pub)
	} else {
		return err
	}

	return nil
}

func GenerateLayer2FollowerSubscription(
	userId,
	requestedSessionId int64,
	clientType enums.ClientType,
	clientVersion byte,
	encryption enums.EncryptionAlgorithm,
	protocolVersion enums.ProtocolVersion,
	sessionResume bool,
) *Layer2FollowerSubscription {
	var challenge []byte = []byte{}

	if tempChallenge, err := helpers.GenerateChallengeNonce(); err == nil {
		challenge = append(challenge, tempChallenge...)
	} else {
		challenge = append(challenge, helpers.GetDefaultChallengeNonce()...)
	}

	return &Layer2FollowerSubscription{
		Challenge:          challenge,
		UserId:             userId,
		RequestedSessionID: requestedSessionId,
		ClientType:         clientType,
		ClientVersion:      clientVersion,
		Encryption:         encryption,
		ProtocolVersion:    protocolVersion,
		SessionResume:      sessionResume,
	}
}
