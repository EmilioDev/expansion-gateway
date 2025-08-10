// file: /dto/layer_2_session.go
package dto

import (
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
}

// GenerateNewLayer2Session creates a new Layer2Session with default values
func GenerateNewLayer2Session() *Layer2Session {
	s := &Layer2Session{
		challenge:     nil,
		subscriptions: make(map[string]struct{}),
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

// ===== RequestedSessionId =====
func (s *Layer2Session) GetRequestedSessionId() int64 {
	return s.requestedSessionId.Load()
}

func (s *Layer2Session) SetRequestedSessionId(id int64) {
	s.requestedSessionId.Store(id)
}

// ===== ProtocolVersion =====
func (s *Layer2Session) GetProtocolVersion() enums.ProtocolVersion {
	return enums.ProtocolVersion(s.protocolVersion.Load())
}

func (s *Layer2Session) SetProtocolVersion(v enums.ProtocolVersion) {
	s.protocolVersion.Store(int32(v))
}

// ===== ClientType =====
func (s *Layer2Session) GetClientType() enums.ClientType {
	return enums.ClientType(s.clientType.Load())
}

func (s *Layer2Session) SetClientType(t enums.ClientType) {
	s.clientType.Store(int32(t))
}

// ===== ClientVersion =====
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

// ===== SessionResume =====
func (s *Layer2Session) GetSessionResume() bool {
	return s.sessionResume.Load()
}

func (s *Layer2Session) SetSessionResume(resume bool) {
	s.sessionResume.Store(resume)
}
