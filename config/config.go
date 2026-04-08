package config

import (
	"encoding/base64"
	"expansion-gateway/helpers"
	"expansion-gateway/internal/structs/tries"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	dotenv "github.com/joho/godotenv"
)

type Configuration struct {
	port                     uint16                // the port used by the layer 1 to listen for connections
	bufferSize               int                   // the max size of each packet
	shardCount               int                   // number of channels to be used between the dispatchers and the receivers
	shardBufferSize          int                   // the number of packets that should be buffered in the packet receivers between the layers
	timeout                  int                   // the connection timeout of each client to the layer 1
	sessionTimeout           time.Duration         // time each session has to do any activity before being declared as inactive and then deleted (seconds)
	sessionWatcherPeriod     time.Duration         // the watcher period, the time between each check to see if there is an idle session (milliseconds)
	godotEd25519PublicKey    []byte                // the public key used in the Ed25519 authentication for godot
	cliEd25519PublicKey      []byte                // the public key used in the Ed25519 authentication for the cli tool
	grpcLeaderPath           string                // path to the leader of the cluster, if this field is empty, this node is a leader in the cluster
	grpcIpAddressWithoutPort string                // path to access to this node
	kcpPathWithoutPort       string                // the kcp path that should be used when redirecting to this node
	grpcClusterPort          uint16                // the port this cluster member is using to receive request from other members
	natsServerPath           string                // path to the NATS server where this client should publish the packets
	publishWindowSize        uint16                // the size of the window to allow a publish packet to be processed if it requests to use window range
	ecoPath                  tries.SubscriptionKey // the path that should be answered with an eco
}

// ===== environment variables, their description, and their default value
// PORT: the port used to listen for kcp connections from internet. default: 7000
// GRPC_CLUSTER_PORT: the port used by the gRPC server to interact with the cluster members. default: 40000
// BUFFER_SIZE: the number of bytes used as buffer for reading from each kcp client. default: 4096
// SHARD_COUNT: the number of sharded channels used to communicate between each layer. default: 8
// SHARD_BUFFER_SIZE: the number of messages that can be stacked in each of the sharded channels between the layers. default: 1024
// CONNECTION_TIMEOUT: the timeout in seconds used by each kcp connection. default: 1
// SESSION_TIMEOUT_SECONDS: the timeout each session has in layer 2 before being disconnected for inactivity. default: 30 seconds
// SESSION_WATCHER_PERIOD_MS: the period between each check to find if any session is inactive. default: 1 second
// CLI_AUTH_KEY: the cli auth key used to validate the cli app. default: [LKopGQQ6iytq9BywdlOfkQoIw0Kpny2xu8F1JXnDAb4=] (base64)
// GODOT_AUTH_KEY: the key used by the game to prove its identity. default: [LKopGQQ6iytq9BywdlOfkQoIw0Kpny2xu8F1JXnDAb4=] (base64)
// CLUSTER_LEADER_GRPC: the path to the cluster leader of this cluster. default: empty
// GRPC_SERVER_IP_PATH_WITHOUT_PORT: the ip, without the port, used by the grpc server of this cluster member. default: "0.0.0.0"
// NODE_KCP_PATH: the ip, without the port, that should be used when redirecting to this cluster member. default: 127.0.0.1
// NATS_PATH: path to the NATS server where to interact to
// PUBLISH_WINDOW_SIZE: the size of the window to allow publish packets. default is 10
// ECO: if something is published in this path, it will be published back to the client, without being forwarded. default is "/eco"

// Initialices this module
func (conf *Configuration) Initialize() {
	dotenv.Load()

	// private: [160 135 112 15 214 124 209 247 138 189 68 26 191 207 215 199 197 87 113 188 237 31 70 77 193 125 135 96 26 250 146 92]
	// public: LKopGQQ6iytq9BywdlOfkQoIw0Kpny2xu8F1JXnDAb4=

	encodeKey, _ := base64.StdEncoding.DecodeString("LKopGQQ6iytq9BywdlOfkQoIw0Kpny2xu8F1JXnDAb4=")
	defaultEd25519PublicKey := [32]byte(encodeKey)

	// defaults
	conf.port = 7000
	conf.bufferSize = 4096 // the size of the buffer used to read from the client
	conf.shardCount = 8    // number of channels used to interact between layers
	conf.shardBufferSize = 1024
	conf.timeout = 1
	conf.sessionTimeout = 30 * time.Second      // default: 30s session timeout
	conf.sessionWatcherPeriod = 1 * time.Second // default: check every 1s
	conf.godotEd25519PublicKey = defaultEd25519PublicKey[:]
	conf.cliEd25519PublicKey = defaultEd25519PublicKey[:]
	conf.grpcLeaderPath = ""
	conf.grpcIpAddressWithoutPort = "0.0.0.0"
	conf.kcpPathWithoutPort = "127.0.0.1"
	conf.grpcClusterPort = 40000 // yes, warhammer 40k!!!!
	conf.natsServerPath = "127.0.0.1:4222"
	conf.publishWindowSize = 10
	conf.ecoPath = tries.SubscriptionKey("/eco")

	// ECO path
	if val := os.Getenv("ECO"); val != "" {
		if path, err := tries.ConvertStringToSubscriptionKey(val); err == nil {
			conf.ecoPath = path
		}
	}

	// PUBLISH window size
	if val := os.Getenv("PUBLISH_WINDOW_SIZE"); val != "" {
		if num, err := strconv.Atoi(val); err == nil {
			conf.publishWindowSize = helpers.ConvertIntToUint16(num)
		}
	}

	// NATS server
	if val := os.Getenv("NATS_PATH"); val != "" {
		conf.natsServerPath = val
	}

	// port
	stringPort := os.Getenv("PORT")

	if candidateToPort, err := strconv.Atoi(stringPort); err == nil {
		if candidateToPort > 1024 {
			conf.port = helpers.ConvertIntToUint16(candidateToPort)
		}
	}

	// cluster grpc port
	if val := os.Getenv("GRPC_CLUSTER_PORT"); val != "" {
		if num, err := strconv.Atoi(val); err == nil && num > 1024 && num < math.MaxUint16 {
			conf.grpcClusterPort = uint16(num)
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
	if s := os.Getenv("GRPC_SERVER_IP_PATH_WITHOUT_PORT"); s != "" {
		conf.grpcIpAddressWithoutPort = helpers.RemoveCharFromString(s, byte('"'))
	}

	// path that should be used when redirecting to this node in kcp
	if s := os.Getenv("NODE_KCP_PATH"); s != "" {
		conf.kcpPathWithoutPort = s
	}
}

func (conf *Configuration) GetEcoPath() tries.SubscriptionKey {
	return conf.ecoPath
}

func (conf *Configuration) GetPublishWindowSize() uint16 {
	return conf.publishWindowSize
}

func (conf *Configuration) GetKcpPathToThisGateway() string {
	return fmt.Sprintf("%s:%d", conf.kcpPathWithoutPort, conf.port)
}

func (conf *Configuration) GetUniversalKcpPathToThisGateway() string {
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
	return fmt.Sprintf("%s:%d", conf.grpcIpAddressWithoutPort, conf.grpcClusterPort)
}

func (conf *Configuration) GetClusterGrpcPort() uint16 {
	return conf.grpcClusterPort
}

func (conf *Configuration) GetNATSserverPath() string {
	return conf.natsServerPath
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
