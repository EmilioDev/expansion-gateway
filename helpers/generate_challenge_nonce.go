package helpers

import (
	"crypto/rand"
	"expansion-gateway/errors/cryptoerror"
	"expansion-gateway/interfaces/errorinfo"
)

func GenerateChallengeNonce() ([]byte, errorinfo.GatewayError) {
	const challengeSize int = 32 // the challenge has 32 bytes, this is for better understanding of code
	challenge := [challengeSize]byte{}

	// we fill the bytes
	if _, err := rand.Read(challenge[:]); err != nil {
		return nil, cryptoerror.CreateErrorWhileGeneratingChallenge(
			"/helpers/generate_challenge_nonce.go",
			14,
			err,
		)
	}

	return challenge[:], nil
}

func GetDefaultChallengeNonce() []byte {
	const CHALLENGE_LEN int = 32
	answer := make([]byte, 0, CHALLENGE_LEN)

	for x := 0; x < CHALLENGE_LEN; x++ {
		answer = append(answer, byte(x))
	}

	return answer
}
