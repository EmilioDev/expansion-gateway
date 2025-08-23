package dto

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

type ConnectedPacket struct {
	sessionTimeout  int64                     // the timeout duration in milliseconds
	sessionID       int64                     // the session id
	sessionResume   bool                      // if this is a resume of a previously existing session, or a new fresh one
	state           enums.ConnectionState     // the state of this session
	protocolVersion enums.ProtocolVersion     // the version of the protocol this user will use
	clientType      enums.ClientType          // the kind of client this user is
	encryption      enums.EncryptionAlgorithm // the encryption algorithm used in the payload
	encryptionKey   []byte                    // the key used by the encryption algorithm
}

func CreateNewConnectPacket(sessionId int64, session *Layer2Session) *ConnectedPacket {
	answer := ConnectedPacket{
		sessionTimeout:  session.GetConfiguration().GetSessionTimeout().Milliseconds(),
		sessionID:       sessionId,
		sessionResume:   session.GetSessionResume(),
		state:           session.GetState(),
		protocolVersion: session.GetProtocolVersion(),
		clientType:      session.GetClientType(),
		encryption:      session.GetEncryption(),
	}

	answer.encryptionKey = []byte{}

	return &answer
}

func (packet ConnectedPacket) GetPacketType() enums.PacketType {
	return enums.CONNECTED
}

func (packet ConnectedPacket) GetPacketHeader() packets.PacketHeader {
	return nil
}

func (packet ConnectedPacket) GetPayload() string {
	return ""
}

func (packet ConnectedPacket) Marshal() ([]byte, errorinfo.GatewayError) {
	return []byte{}, nil
}

func (packet ConnectedPacket) GetSender() int64 {
	return packet.sessionID
}
