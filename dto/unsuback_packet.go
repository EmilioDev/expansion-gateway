package dto

import (
	"expansion-gateway/enums"
	"expansion-gateway/helpers"
	errors "expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
	"fmt"
)

type UnsubackPacket struct {
	unsusbscriptionPacketId int32
	owner                   int64
}

// creates an UNSUBACK packet
func CreateUnsubackPacket(owner int64, pid int32) *UnsubackPacket {
	return &UnsubackPacket{
		owner:                   owner,
		unsusbscriptionPacketId: pid,
	}
}

func (packet *UnsubackPacket) GetPacketType() enums.PacketType {
	return enums.UNSUBACK
}

func (packet *UnsubackPacket) GetPacketHeader() packets.PacketHeader {
	return nil
}

func (packet *UnsubackPacket) GetPayload() string {
	return ""
}

func (packet *UnsubackPacket) GetSender() int64 {
	return packet.owner
}

func (packet *UnsubackPacket) GetRawPayload() []byte {
	return []byte{}
}

func (packet *UnsubackPacket) GetIdentifier() string {
	return fmt.Sprintf("%d", packet.unsusbscriptionPacketId)
}

func (packet *UnsubackPacket) SetNewOwner(newOwner int64) {
	packet.owner = newOwner
}

func (packet *UnsubackPacket) Marshal() ([]byte, errors.GatewayError) {
	output := make([]byte, 0, 1+4)

	unsubId := helpers.ConvertInt32Into4Bytes(packet.unsusbscriptionPacketId)

	output = append(output, byte(enums.UNSUBACK))
	output = append(output, unsubId[:]...)

	return output, nil
}
