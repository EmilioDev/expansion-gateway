package dto

import (
	"expansion-gateway/enums"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
	"expansion-gateway/interfaces/sessions"
)

type ConnectedPacket struct {
	sessionTimeout  int64                     // the timeout duration in milliseconds
	sessionID       int64                     // the session id
	sessionResume   bool                      // if this is a resume of a previously existing session, or a new fresh one
	state           enums.ConnectionState     // the state of this session
	protocolVersion enums.ProtocolVersion     // the version of the protocol this user will use
	clientType      enums.ClientType          // the kind of client this user is
	clientVersion   byte                      // the version this client is using
	encryption      enums.EncryptionAlgorithm // the encryption algorithm used in the payload
}

func CreateNewConnectedPacket(sessionId int64, session sessions.Session) *ConnectedPacket {
	answer := ConnectedPacket{
		sessionTimeout:  session.GetConfiguration().GetSessionTimeout().Milliseconds(),
		sessionID:       sessionId,
		sessionResume:   session.GetSessionResume(),
		state:           session.GetState(),
		protocolVersion: session.GetProtocolVersion(),
		clientType:      session.GetClientType(),
		clientVersion:   session.GetClientVersion(),
		encryption:      session.GetEncryption(),
	}

	return &answer
}

func (packet *ConnectedPacket) GetPacketType() enums.PacketType {
	return enums.CONNECTED
}

func (packet *ConnectedPacket) GetPacketHeader() packets.PacketHeader {
	return nil
}

func (packet *ConnectedPacket) GetPayload() string {
	return ""
}

func (packet *ConnectedPacket) Marshal() ([]byte, errorinfo.GatewayError) {
	output := make([]byte, 0, 1+8+2+4+1+8+2+2+2)

	sessionId := helpers.ConvertInt64Into8Bytes(packet.sessionID)

	output = append(output, byte(enums.CONNECTED))
	output = append(output, sessionId[:]...) // session id

	// encryption
	if packet.encryption != enums.NoEncryptionAlgorithm {
		output = append(output, 25)
		output = append(output, byte(packet.encryption))
	}

	output = append(output, 13)
	timeout := helpers.ConvertInt64Into8Bytes(packet.sessionTimeout)
	output = append(output, timeout[:]...)

	// protocol version
	output = append(output, 48)
	output = append(output, byte(packet.protocolVersion))

	// client type
	output = append(output, 31)
	output = append(output, byte(packet.clientType))

	// client version
	output = append(output, 82)
	output = append(output, packet.clientVersion)

	return output, nil
}

func (packet *ConnectedPacket) GetSender() int64 {
	return packet.sessionID
}
