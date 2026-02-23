package dto

import (
	"expansion-gateway/enums"
	"expansion-gateway/errors"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

type HelloPacket struct {
	Sender         int64                      // connection id assigned to this user when connected
	VariableHeader *HelloPacketVariableHeader // variable header of this packet
}

type HelloPacketVariableHeader struct {
	BaseHeader                                 // The base header struct
	ProtocolVersion  enums.ProtocolVersion     // protocol version requested by the user
	ClientType       enums.ClientType          // client type (godot, cli...)
	ClientVersion    byte                      // version of this client
	PayloadEncrypted bool                      // if this user is going to encrypt the payload
	SessionResume    bool                      // if this is a resume of a previously open session
	Encryption       enums.EncryptionAlgorithm // encryptation algorythm this user is going to use in the payload
	PretendedUserID  int64                     // the id of the previously open session this user wants to continue using
}

func CreateHelloPacket(sender int64) *HelloPacket {
	return &HelloPacket{
		Sender: sender,
		VariableHeader: &HelloPacketVariableHeader{
			BaseHeader:       BaseHeader{},
			ProtocolVersion:  enums.V1,
			ClientType:       enums.CLI_TOOL,
			ClientVersion:    1,
			PayloadEncrypted: false,
			SessionResume:    false,
			Encryption:       enums.NoEncryptionAlgorithm,
			PretendedUserID:  0,
		},
	}
}

func (packet *HelloPacket) GetPacketType() enums.PacketType {
	return enums.HELLO
}

func (packet *HelloPacket) GetPacketHeader() packets.PacketHeader {
	return packet.VariableHeader
}

func (packet *HelloPacket) GetPayload() string {
	return ""
}

func (packet *HelloPacket) GetSender() int64 {
	return packet.Sender
}

func (packet *HelloPacket) GetRawPayload() []byte {
	return []byte{}
}

func (packet *HelloPacket) Marshal() ([]byte, errorinfo.GatewayError) {
	const filePath string = "/dto/hello_packet.go"
	answer := []byte{
		byte(enums.HELLO),
		byte(packet.VariableHeader.ProtocolVersion),
		byte(packet.VariableHeader.ClientType),
		packet.VariableHeader.ClientVersion,
	}

	var currentByte byte = 0x00

	const (
		flagEncrypted     byte = 1 << 0
		flagSessionResume byte = 1 << 1
	)

	if packet.VariableHeader.PayloadEncrypted {
		currentByte |= flagEncrypted
	}

	if packet.VariableHeader.SessionResume {
		currentByte |= flagSessionResume
	}

	answer = append(answer, currentByte)

	if packet.VariableHeader.PayloadEncrypted {
		if packet.VariableHeader.Encryption < enums.NoEncryptionAlgorithm {
			answer = append(answer, byte(packet.VariableHeader.Encryption))
		} else {
			return nil, errors.CreatePacketWithInvalidEncryptionAlgorythm(
				filePath,
				72,
				enums.HELLO,
				byte(packet.VariableHeader.Encryption),
			)
		}
	}

	if packet.VariableHeader.SessionResume {
		session_id := helpers.ConvertInt64Into8Bytes(packet.VariableHeader.PretendedUserID)
		answer = append(answer, session_id[:]...)
	}

	return answer, nil
}

func (packet *HelloPacket) GetIdentifier() string {
	return ""
}

func (packet *HelloPacket) SetNewOwner(newOwner int64) {
	packet.Sender = newOwner
}
