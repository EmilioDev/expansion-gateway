package dto

import (
	"expansion-gateway/enums"
	"expansion-gateway/helpers"
	errors "expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
	"expansion-gateway/internal/structs/tries"
)

type SubscribePacket struct {
	Key                   tries.SubscriptionKey
	susbscriptionPacketId int32
	owner                 int64
}

func CreateSubscribePacket(key tries.SubscriptionKey, packetId int32, owner int64) *SubscribePacket {
	return &SubscribePacket{
		Key:                   key,
		owner:                 owner,
		susbscriptionPacketId: packetId,
	}
}

func (packet *SubscribePacket) GetPacketType() enums.PacketType {
	return enums.SUBSCRIBE
}

func (packet *SubscribePacket) GetSubscriptionID() int32 {
	return packet.susbscriptionPacketId
}

func (packet *SubscribePacket) GetPacketHeader() packets.PacketHeader {
	return nil
}

func (packet *SubscribePacket) GetPayload() string {
	return ""
}

func (packet *SubscribePacket) GetSender() int64 {
	return packet.owner
}

func (packet *SubscribePacket) Marshal() ([]byte, errors.GatewayError) {
	kl := packet.Key.KeyLength()
	output := make([]byte, 0, 1+4+4+kl)

	subId := helpers.ConvertInt32Into4Bytes(packet.susbscriptionPacketId)

	key := packet.Key.ToByteArray()
	keyLen := helpers.ConvertInt32Into4Bytes(int32(kl))

	output = append(output, byte(enums.SUBSCRIBE))
	output = append(output, subId[:]...)
	output = append(output, keyLen[:]...)
	output = append(output, key...)

	return output, nil
}
