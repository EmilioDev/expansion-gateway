package config

import (
	"expansion-gateway/helpers"
	"fmt"
	"os"
	"strconv"

	dotenv "github.com/joho/godotenv"
)

type Configuration struct {
	port            uint16
	bufferSize      int
	shardCount      int
	shardBufferSize int
}

// Initialices this module
func (conf *Configuration) Initialize() {
	dotenv.Load()

	// defaults
	conf.port = 7000
	conf.bufferSize = 4096
	conf.shardCount = 8
	conf.shardBufferSize = 1024

	// port
	stringPort := os.Getenv("PORT")

	if candidateToPort, err := strconv.Atoi(stringPort); err == nil {
		if candidateToPort > 1024 {
			conf.port = helpers.ConvertIntToUint16(candidateToPort)
		}
	}

	// buffer size
	stringBufferSize := os.Getenv("BUFFER_SIZE")

	if candidateToBufferSize, err := strconv.Atoi(stringBufferSize); err == nil {
		if candidateToBufferSize > 0 {
			conf.bufferSize = candidateToBufferSize
		}
	}

	// Load shard count
	if val := os.Getenv("SHARD_COUNT"); val != "" {
		if num, err := strconv.Atoi(val); err == nil && num > 0 {
			conf.shardCount = num
		}
	}

	// Load shard buffer size
	if val := os.Getenv("SHARD_BUFFER_SIZE"); val != "" {
		if num, err := strconv.Atoi(val); err == nil && num > 0 {
			conf.shardBufferSize = num
		}
	}
}

// returns the server address to be used in this server
func (conf *Configuration) GetServerAddress() string {
	return fmt.Sprintf("0.0.0.0:%d", conf.port)
}

func (conf *Configuration) GetBufferSize() int {
	return conf.bufferSize
}

func (conf *Configuration) GetShardCount() int {
	return conf.shardCount
}

func (conf *Configuration) GetShardBufferSize() int {
	return conf.shardBufferSize
}
