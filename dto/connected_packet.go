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
	encryptionKey   []byte                    // the key used by the encryption algorithm
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

	answer.encryptionKey = []byte{}

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
	output := make([]byte, 0, 1+8+2+4+len(packet.encryptionKey)+1+8+2+2+2)

	sessionId := helpers.ConvertInt64Into8Bytes(packet.sessionID)

	output[0] = byte(enums.CONNECTED) // packet id
	copy(output[1:], sessionId[:])    // session id

	size := 9

	// encryption
	if packet.encryption != enums.NoEncryptionAlgorithm {
		output[9] = 25
		output[10] = byte(packet.encryption)

		keySize := helpers.ConvertInt32Into4Bytes(int32(len(packet.encryptionKey)))

		copy(output[11:], keySize[:])
		copy(output[15:], packet.encryptionKey)
		size += 6 + len(packet.encryptionKey)
	}

	// timeout
	output[size] = 13

	timeout := helpers.ConvertInt64Into8Bytes(packet.sessionTimeout)
	copy(output[size+1:], timeout[:])

	size += 9

	// protocol version
	output[size] = 48
	output[size+1] = byte(packet.protocolVersion)

	size += 2

	// client type
	output[size] = 31
	output[size+1] = byte(packet.clientType)

	size += 2

	// client version
	output[size] = 82
	output[size+1] = packet.clientVersion

	size += 2

	return output[:size], nil
}

func (packet *ConnectedPacket) GetSender() int64 {
	return packet.sessionID
}
