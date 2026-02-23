// file: /dto/connect_packet.go
package dto

import (
	"expansion-gateway/enums"
	errors "expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

type ConnectPacket struct {
	OwnerID            int64     // the id of the session of the user that owns this packet
	Signature          [64]byte  // the signature of the data the user is going to send
	ClientEphemeralKey *[32]byte // the ephemeral key of this client
}

func (packet *ConnectPacket) GetPacketType() enums.PacketType {
	return enums.CONNECT
}

func (packet *ConnectPacket) GetPacketHeader() packets.PacketHeader {
	return nil
}

func (packet *ConnectPacket) GetPayload() string {
	return ""
}

func (packet *ConnectPacket) Marshal() ([]byte, errors.GatewayError) {
	answer := make([]byte, 0, 65)

	signature := packet.Signature[:]

	answer = append(answer, byte(enums.CONNECT))
	answer = append(answer, signature...)

	return answer, nil
}

func (packet *ConnectPacket) GetSender() int64 {
	return packet.OwnerID
}

func (packet *ConnectPacket) GetRawPayload() []byte {
	return []byte{}
}

func (packet *ConnectPacket) GetIdentifier() string {
	return ""
}

func (packet *ConnectPacket) SetNewOwner(newOwner int64) {
	packet.OwnerID = newOwner
}

func CreateConnectPacketWithoutEphemeralKey(ownerId int64, signature [64]byte) *ConnectPacket {
	return &ConnectPacket{
		OwnerID:            ownerId,
		Signature:          signature,
		ClientEphemeralKey: nil,
	}
}

func CreateConnectPacketWithEphemeralKey(ownerId int64, signature [64]byte, key [32]byte) *ConnectPacket {
	return &ConnectPacket{
		OwnerID:            ownerId,
		Signature:          signature,
		ClientEphemeralKey: &key,
	}
}
