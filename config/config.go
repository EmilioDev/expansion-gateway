package config

import (
	"expansion-gateway/helpers"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	dotenv "github.com/joho/godotenv"
)

type Configuration struct {
	port                  uint16        // the port used by the layer 1 to listen for connections
	bufferSize            int           // the max size of each packet
	shardCount            int           // number of channels to be used between the dispatchers and the receivers
	shardBufferSize       int           // the number of packets that should be buffered in the packet receivers between the layers
	timeout               int           // the connection timeout of each client to the layer 1
	sessionTimeout        time.Duration // time each session has to do any activity before being declared as inactive and then deleted (seconds)
	sessionWatcherPeriod  time.Duration // the watcher period, the time between each check to see if there is an idle session (milliseconds)
	godotEd25519PublicKey []byte        // the public key used in the Ed25519 authentication for godot
	cliEd25519PublicKey   []byte        // the public key used in the Ed25519 authentication for the cli tool
	grpcLeaderPath        string        // path to the leader of the cluster, if this field is empty, this node is a leader in the cluster
	grpcCurrentPath       string        // path to access to this node
	kcpPath               string        // the kcp path that should be used when redirecting to this node
}

// Initialices this module
func (conf *Configuration) Initialize() {
	dotenv.Load()

	defaultEd25519PublicKey := [32]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}

	// defaults
	conf.port = 7000
	conf.bufferSize = 4096
	conf.shardCount = 8
	conf.shardBufferSize = 1024
	conf.timeout = 1
	conf.sessionTimeout = 30 * time.Second      // default: 30s session timeout
	conf.sessionWatcherPeriod = 1 * time.Second // default: check every 1s
	conf.godotEd25519PublicKey = defaultEd25519PublicKey[:]
	conf.cliEd25519PublicKey = defaultEd25519PublicKey[:]
	conf.grpcLeaderPath = ""
	conf.grpcCurrentPath = "127.0.0.1:4500"
	conf.kcpPath = "127.0.0.1:7000"

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

	// CLI Ed25519 Public Key
	if s := os.Getenv("CLI_AUTH_KEY"); s != "" {
		conf.cliEd25519PublicKey = parseKey(s, defaultEd25519PublicKey)
	}

	// Godot Ed25519 Public Key
	if s := os.Getenv("GODOT_AUTH_KEY"); s != "" {
		conf.godotEd25519PublicKey = parseKey(s, defaultEd25519PublicKey)
	}

	// grpc cluster leader path
	if s := os.Getenv("CLUSTER_LEADER_GRPC"); s != "" {
		conf.grpcLeaderPath = s
	}

	// path to the current grpc server
	if s := os.Getenv("GRPC_SERVER_PATH"); s != "" {
		conf.grpcCurrentPath = s
	}

	// path that should be used when redirecting to this node in kcp
	if s := os.Getenv("NODE_KCP_PATH"); s != "" {
		conf.kcpPath = s
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

func (conf *Configuration) GetCliEd25519PublicKey() *[]byte {
	return &conf.cliEd25519PublicKey
}

func (conf *Configuration) GetGodotEd25519PublicKey() *[]byte {
	return &conf.godotEd25519PublicKey
}

func (conf *Configuration) AreWeClusterLeaders() bool {
	return conf.grpcLeaderPath == ""
}

func (conf *Configuration) GetGrpcClusterLeaderPath() string {
	return conf.grpcLeaderPath
}

func (conf *Configuration) GetGrpcCurrentServerPath() string {
	return conf.grpcCurrentPath
}

// ==== extras and local helpers ====

// parses a ed25519 key string into a byte array
func parseKey(input string, defaultValue [32]byte) []byte {
	if input == "" {
		return defaultValue[:]
	}

	answer := []byte{}

	theNumbers := strings.FieldsFunc(input, func(r rune) bool {
		switch r {
		case '.', ':', ' ', ',', '-':
			return true

		default:
			return false
		}
	})

	if len(theNumbers) < 32 {
		return defaultValue[:]
	}

	// the key always should be 32 bytes length
	for x := 0; x < 32; x++ {
		if currentNumber, err := strconv.Atoi(theNumbers[x]); err == nil {
			answer = append(answer, byte(currentNumber))
		} else {
			return defaultValue[:]
		}
	}

	return answer
}
