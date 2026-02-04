package dto

import (
	"expansion-gateway/enums"
	"expansion-gateway/helpers"
	errors "expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
	"expansion-gateway/internal/structs/tries"
)

type UnsubscribePacket struct {
	Key                     tries.SubscriptionKey
	unsusbscriptionPacketId int32
	owner                   int64
}

// creates an UNSUBSCRIBE packet
func CreateUnsubscribePacket(key tries.SubscriptionKey, packetId int32, owner int64) *UnsubscribePacket {
	return &UnsubscribePacket{
		Key:                     key,
		owner:                   owner,
		unsusbscriptionPacketId: packetId,
	}
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

func (packet *UnsubscribePacket) GetUnsubscriptionID() int32 {
	return packet.unsusbscriptionPacketId
}

func (packet *UnsubscribePacket) Marshal() ([]byte, errors.GatewayError) {
	kl := packet.Key.KeyLength()
	output := make([]byte, 0, 1+4+4+kl)

	unsubId := helpers.ConvertInt32Into4Bytes(packet.unsusbscriptionPacketId)

	key := packet.Key.ToByteArray()
	keyLen := helpers.ConvertInt32Into4Bytes(int32(kl))

	output = append(output, byte(enums.UNSUBSCRIBE))
	output = append(output, unsubId[:]...)
	output = append(output, keyLen[:]...)
	output = append(output, key...)

	return output, nil
}
