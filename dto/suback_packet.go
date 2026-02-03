package dto

import (
	"expansion-gateway/enums"
	"expansion-gateway/helpers"
	errors "expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

type SubackPacket struct {
	susbscriptionPacketId int32
	owner                 int64
}

func CreateSubackPacket(owner int64, pid int32) *SubackPacket {
	return &SubackPacket{
		owner:                 owner,
		susbscriptionPacketId: pid,
	}
}

func (packet *SubackPacket) GetPacketType() enums.PacketType {
	return enums.SUBACK
}

func (packet *SubackPacket) GetPacketHeader() packets.PacketHeader {
	return nil
}

func (packet *SubackPacket) GetPayload() string {
	return ""
}

func (packet *SubackPacket) GetSender() int64 {
	return packet.owner
}

func (packet *SubackPacket) Marshal() ([]byte, errors.GatewayError) {
	output := make([]byte, 0, 1+4)

	subId := helpers.ConvertInt32Into4Bytes(packet.susbscriptionPacketId)

	output = append(output, byte(enums.SUBACK))
	output = append(output, subId[:]...)

	return output, nil
}
