package config

import (
	"expansion-gateway/helpers"
	"fmt"
	"os"
	"strconv"
	"time"

	dotenv "github.com/joho/godotenv"
)

type Configuration struct {
	port                 uint16        // the port used by the layer 1 to listen for connections
	bufferSize           int           // the max size of each packet
	shardCount           int           // number of channels to be used between the dispatchers and the receivers
	shardBufferSize      int           // the number of packets that should be buffered in the packet receivers between the layers
	timeout              int           // the connection timeout of each client to the layer 1
	sessionTimeout       time.Duration // time each session has to do any activity before being declared as inactive and then deleted (seconds)
	sessionWatcherPeriod time.Duration // the watcher period, the time between each check to see if there is an idle session (milliseconds)
}

// Initialices this module
func (conf *Configuration) Initialize() {
	dotenv.Load()

	// defaults
	conf.port = 7000
	conf.bufferSize = 4096
	conf.shardCount = 8
	conf.shardBufferSize = 1024
	conf.timeout = 1
	conf.sessionTimeout = 30 * time.Second      // default: 30s session timeout
	conf.sessionWatcherPeriod = 1 * time.Second // default: check every 1s

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

	// connection timeout
	if val := os.Getenv("CONNECTION_TIMEOUT"); val != "" {
		if num, err := strconv.Atoi(val); err == nil && num > 0 {
			conf.timeout = num
		}
	}

	// session timeout (seconds)
	if s := os.Getenv("SESSION_TIMEOUT_SECONDS"); s != "" {
		if sec, err := strconv.Atoi(s); err == nil && sec >= 0 {
			conf.sessionTimeout = time.Duration(sec) * time.Second
		}
	}

	// session watcher period (seconds)
	if s := os.Getenv("SESSION_WATCHER_PERIOD_MS"); s != "" {
		// allow ms, parse int
		if ms, err := strconv.Atoi(s); err == nil && ms > 0 {
			conf.sessionWatcherPeriod = time.Duration(ms) * time.Millisecond
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

func (conf *Configuration) GetConnectionTimeout() int {
	return conf.timeout
}

func (conf *Configuration) GetSessionTimeout() time.Duration {
	return conf.sessionTimeout
}

func (conf *Configuration) GetSessionWatcherPeriod() time.Duration {
	return conf.sessionWatcherPeriod
}
