package helpers

import (
	"crypto/rand"
	"expansion-gateway/errors/cryptoerror"
	"expansion-gateway/interfaces/errorinfo"
)

func GenerateChallengeNonce() ([]byte, errorinfo.GatewayError) {
	const challengeSize = 32
	challenge := make([]byte, challengeSize)
	_, err := rand.Read(challenge)

	if err != nil {
		return nil, cryptoerror.CreateErrorWhileGeneratingChallenge(
			"/helpers/generate_challenge_nonce.go",
			14,
			err,
		)
	}

	return challenge, nil
}

func GetDefaultChallengeNonce() []byte {
	const CHALLENGE_LEN int = 32
	answer := make([]byte, 0, CHALLENGE_LEN)

	for x := 0; x < CHALLENGE_LEN; x++ {
		answer = append(answer, byte(x))
	}

	return answer
}
