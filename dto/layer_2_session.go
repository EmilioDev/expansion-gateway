// file: /dto/layer_2_session.go
package dto

import (
	"expansion-gateway/config"
	"expansion-gateway/enums"
	"sync"
	"sync/atomic"
)

type Layer2Session struct {
	// atomic fields
	state              atomic.Int32
	requestedSessionId atomic.Int64
	protocolVersion    atomic.Int32
	clientType         atomic.Int32
	clientVersion      atomic.Uint32 // byte stored as uint32
	encryption         atomic.Int32
	sessionResume      atomic.Bool

	// protected by their own mutex
	challengeMu sync.RWMutex
	challenge   []byte

	subsMu        sync.RWMutex
	subscriptions map[string]struct{}

	// for the bulk updaters
	bulkUpdaterMutex sync.Mutex

	// timeout tracker
	timeoutTracker *TimeoutTracker

	// configuration object
	configuration *config.Configuration

	// encryption key
	encryptionKey []byte
}

// GenerateNewLayer2Session creates a new Layer2Session with default values
func GenerateNewLayer2Session(config *config.Configuration) *Layer2Session {
	s := &Layer2Session{
		challenge:        nil,
		subscriptions:    make(map[string]struct{}),
		challengeMu:      sync.RWMutex{},
		subsMu:           sync.RWMutex{},
		bulkUpdaterMutex: sync.Mutex{},
		timeoutTracker:   NewTimeoutTracker(config.GetSessionTimeout()),
		configuration:    config,
	}

	s.state.Store(int32(enums.HELLO_RECEIVED))
	s.requestedSessionId.Store(0)
	s.protocolVersion.Store(int32(enums.V1))
	s.clientType.Store(int32(enums.GODOT_CLIENT))
	s.clientVersion.Store(uint32(0))
	s.encryption.Store(int32(enums.NoEncryptionAlgorithm))
	s.sessionResume.Store(false)

	return s
}

// ===== State =====
func (s *Layer2Session) GetState() enums.ConnectionState {
	return enums.ConnectionState(s.state.Load())
}

func (s *Layer2Session) SetState(newState enums.ConnectionState) {
	s.state.Store(int32(newState))
}

// ===== Challenge =====
func (s *Layer2Session) GetChallenge() []byte {
	s.challengeMu.RLock()
	defer s.challengeMu.RUnlock()

	if s.challenge == nil {
		return nil
	}

	// return copy to avoid external modification
	cpy := make([]byte, len(s.challenge))
	copy(cpy, s.challenge)

	return cpy
}

func (s *Layer2Session) SetChallenge(newChallenge *[]byte) {
	s.challengeMu.Lock()
	defer s.challengeMu.Unlock()

	if newChallenge == nil {
		s.challenge = nil
		return
	}

	s.challenge = make([]byte, len(*newChallenge))
	copy(s.challenge, *newChallenge)
}

// ===== Subscriptions =====
func (s *Layer2Session) GetSubscriptions() []string {
	s.subsMu.RLock()
	defer s.subsMu.RUnlock()

	if len(s.subscriptions) == 0 {
		return nil
	}

	out := make([]string, 0, len(s.subscriptions))

	for t := range s.subscriptions {
		out = append(out, t)
	}

	return out
}

func (s *Layer2Session) HasSubscription(topic string) bool {
	s.subsMu.RLock()
	defer s.subsMu.RUnlock()

	_, hasSubscription := s.subscriptions[topic]

	return hasSubscription
}

func (s *Layer2Session) AddSubscription(topic string) {
	s.subsMu.Lock()
	defer s.subsMu.Unlock()

	s.subscriptions[topic] = struct{}{}
}

func (s *Layer2Session) ClearSubscriptions() {
	s.subsMu.Lock()
	defer s.subsMu.Unlock()

	s.subscriptions = make(map[string]struct{})
}

func (s *Layer2Session) RemoveSubscription(topic string) {
	s.subsMu.Lock()
	defer s.subsMu.Unlock()

	delete(s.subscriptions, topic)
}

// ===== Requested Session Id =====
func (s *Layer2Session) GetRequestedSessionId() int64 {
	return s.requestedSessionId.Load()
}

func (s *Layer2Session) SetRequestedSessionId(id int64) {
	s.requestedSessionId.Store(id)
}

// ===== Protocol Version =====
func (s *Layer2Session) GetProtocolVersion() enums.ProtocolVersion {
	return enums.ProtocolVersion(s.protocolVersion.Load())
}

func (s *Layer2Session) SetProtocolVersion(v enums.ProtocolVersion) {
	s.protocolVersion.Store(int32(v))
}

// ===== Client Type =====
func (s *Layer2Session) GetClientType() enums.ClientType {
	return enums.ClientType(s.clientType.Load())
}

func (s *Layer2Session) SetClientType(t enums.ClientType) {
	s.clientType.Store(int32(t))
}

// ===== Client Version =====
func (s *Layer2Session) GetClientVersion() byte {
	return byte(s.clientVersion.Load())
}

func (s *Layer2Session) SetClientVersion(v byte) {
	s.clientVersion.Store(uint32(v))
}

// ===== Encryption =====
func (s *Layer2Session) GetEncryption() enums.EncryptionAlgorithm {
	return enums.EncryptionAlgorithm(s.encryption.Load())
}

func (s *Layer2Session) SetEncryption(e enums.EncryptionAlgorithm) {
	s.encryption.Store(int32(e))
}

// ===== Session Resume =====
func (s *Layer2Session) GetSessionResume() bool {
	return s.sessionResume.Load()
}

func (s *Layer2Session) SetSessionResume(resume bool) {
	s.sessionResume.Store(resume)
}

// ==== Session Timeout ====

// TimeoutTracker returns the pointer to the tracker's object.
func (s *Layer2Session) TimeoutTracker() *TimeoutTracker {
	return s.timeoutTracker
}

// RefreshActivity marks activity now.
func (s *Layer2Session) RefreshActivity() {
	if s.timeoutTracker != nil {
		s.timeoutTracker.Refresh()
	}
}

// ==== Bulk updaters ====
func (session *Layer2Session) UpdateFromHelloPacket(packet *HelloPacket) {
	// lock packet to avoid race conditions
	session.bulkUpdaterMutex.Lock()
	defer session.bulkUpdaterMutex.Unlock()

	// encryption used in the payload
	if !packet.VariableHeader.PayloadEncrypted {
		session.encryption.Store(int32(enums.NoEncryptionAlgorithm))
	} else if enums.IsValidEncryptionAlgorythm(byte(packet.VariableHeader.Encryption)) {
		session.encryption.Store(int32(packet.VariableHeader.Encryption))
	} else {
		session.encryption.Store(int32(enums.NoEncryptionAlgorithm))
	}

	// session resume
	session.sessionResume.Store(packet.VariableHeader.SessionResume)

	// client version
	session.clientVersion.Store(uint32(packet.VariableHeader.ClientVersion))

	// client type
	session.clientType.Store(int32(packet.VariableHeader.ClientType))

	// protocol version
	session.protocolVersion.Store(int32(packet.VariableHeader.ProtocolVersion))

	// requested session id to resume
	if session.GetSessionResume() {
		session.requestedSessionId.Store(packet.VariableHeader.PretendedUserID)
	} else {
		session.requestedSessionId.Store(0)
	}
}

func (session *Layer2Session) GetEd25519PublicKey() []byte {
	switch session.GetClientType() {
	case enums.CLI_TOOL:
		return *session.configuration.GetCliEd25519PublicKey()

	case enums.GODOT_CLIENT:
		return *session.configuration.GetGodotEd25519PublicKey()

	default:
		return []byte{}
	}
}

func (s *Layer2Session) GetConfiguration() *config.Configuration {
	return s.configuration
}
