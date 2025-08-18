// file: /dto/connect_packet.go
package dto

import (
	"expansion-gateway/enums"
	errors "expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

type ConnectPacket struct {
	OwnerID   int64    // the id of the session of the user that owns this packet
	Signature [64]byte // the signature of the data the user is going to send
}

func (packet ConnectPacket) GetPacketType() enums.PacketType {
	return enums.CONNECT
}

func (packet ConnectPacket) GetPacketHeader() packets.PacketHeader {
	return nil
}

func (packet ConnectPacket) GetPayload() string {
	return ""
}

func (packet ConnectPacket) Marshal() ([]byte, errors.GatewayError) {
	answer := make([]byte, 0, 65)

	signature := packet.Signature[:]

	answer = append(answer, byte(enums.CONNECT))
	answer = append(answer, signature...)

	return answer, nil
}

func (packet ConnectPacket) GetSender() int64 {
	return packet.OwnerID
}

func CreateConnectPacket(ownerId int64, signature [64]byte) *ConnectPacket {
	return &ConnectPacket{
		OwnerID:   ownerId,
		Signature: signature,
	}
}
