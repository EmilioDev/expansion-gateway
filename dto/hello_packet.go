package dto

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/packets"
)

type HelloPacket struct {
	Sender         int64                      // connection id assigned to this user when connected
	VariableHeader *HelloPacketVariableHeader // variable header of this packet
}

type HelloPacketVariableHeader struct {
	ProtocolVersion  enums.ProtocolVersion       // protocol version requested by the user
	ClientType       enums.ClientType            // client type (godot, cli...)
	ClientVersion    byte                        // version of this client
	PayloadEncrypted bool                        // if this user is going to encrypt the payload
	SessionResume    bool                        // if this is a resume of a previously open session
	Encryptation     enums.EncryptationAlgorythm // encryptation algorythm this user is going to use in the payload
	PretendedUserID  int64                       // the id of the previously open session this user wants to continue using
}

func (packet HelloPacket) GetPacketType() enums.PacketType {
	return enums.HELLO
}

func (packet HelloPacket) GetPacketHeader() packets.PacketHeader {
	return packet.VariableHeader
}

func (packet HelloPacket) GetPayload() string {
	return ""
}
