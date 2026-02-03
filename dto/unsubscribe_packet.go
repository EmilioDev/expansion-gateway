package dto

import (
	"expansion-gateway/enums"
	"expansion-gateway/helpers"
	errors "expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
	"expansion-gateway/internal/structs/tries"
)

type UnsubscribePacket struct {
	Key                   tries.SubscriptionKey
	susbscriptionPacketId int32
	owner                 int64
}

func (packet *UnsubscribePacket) GetPacketType() enums.PacketType {
	return enums.UNSUBSCRIBE
}

func (packet *UnsubscribePacket) GetPacketHeader() packets.PacketHeader {
	return nil
}

func (packet *UnsubscribePacket) GetPayload() string {
	return ""
}

func (packet *UnsubscribePacket) GetSender() int64 {
	return packet.owner
}

func (packet *UnsubscribePacket) Marshal() ([]byte, errors.GatewayError) {
	kl := packet.Key.KeyLength()
	output := make([]byte, 0, 1+4+4+kl)

	subId := helpers.ConvertInt32Into4Bytes(packet.susbscriptionPacketId)

	key := packet.Key.ToByteArray()
	keyLen := helpers.ConvertInt32Into4Bytes(int32(kl))

	output = append(output, byte(enums.UNSUBSCRIBE))
	output = append(output, subId[:]...)
	output = append(output, keyLen[:]...)
	output = append(output, key...)

	return output, nil
}
