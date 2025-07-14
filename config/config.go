package config

import (
	"expansion-gateway/helpers"
	"fmt"
	"os"
	"strconv"

	dotenv "github.com/joho/godotenv"
)

type Configuration struct {
	Port uint16
}

func (conf *Configuration) Initialize() {
	dotenv.Load()

	// port
	stringPort := os.Getenv("PORT")

	if candidateToPort, err := strconv.Atoi(stringPort); err == nil {
		if candidateToPort > 1024 {
			conf.Port = helpers.ConvertIntToUint16(candidateToPort)
			return
		}
	}

	conf.Port = 7000
}

// returns the server address to be used in this server
func (conf *Configuration) GetServerAddress() string {
	return fmt.Sprintf("0.0.0.0:%d", conf.Port)
}
