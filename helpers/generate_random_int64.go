package helpers

import (
	"crypto/rand"
	"math/big"
	"time"
)

func GenerateRandomInt64() int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(27))

	if err != nil {
		return 0
	}

	if time.Now().Unix()%2 == 0 {
		return nBig.Int64()
	}

	return nBig.Int64() * -1
}
