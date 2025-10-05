package sessions

import (
	"expansion-gateway/config"
	"expansion-gateway/enums"
)

type Session interface {
	GetConfiguration() *config.Configuration
	GetSessionResume() bool
	GetState() enums.ConnectionState
	SetState(enums.ConnectionState)
	GetProtocolVersion() enums.ProtocolVersion
	GetClientType() enums.ClientType
	GetEncryption() enums.EncryptionAlgorithm
}
