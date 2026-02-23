package dto

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
	"expansion-gateway/internal/structs/tries"
)

type PublishPacket struct {
	Key             tries.SubscriptionKey
	publishPacketId int32
	owner           int64
	payload         []byte
}

func CreatePublishPacket(key tries.SubscriptionKey, packetId int32, owner int64, payload []byte) *PublishPacket {
	return &PublishPacket{
		Key:             key,
		publishPacketId: packetId,
		owner:           owner,
		payload:         payload,
	}
}

func CreateBasicPublishPacket(key tries.SubscriptionKey, payload []byte) *PublishPacket {
	return &PublishPacket{
		Key:             key,
		publishPacketId: 0,
		owner:           0,
		payload:         payload,
	}
}

func (packet *PublishPacket) GetPacketType() enums.PacketType {
	return enums.PUBLISH
}

func (packet *PublishPacket) GetPublishPacketID() int32 {
	return packet.publishPacketId
}

func (packet *PublishPacket) GetPacketHeader() packets.PacketHeader {
	return nil
}

func (packet *PublishPacket) GetPayload() string {
	return string(packet.payload)
}

func (packet *PublishPacket) GetRawPayload() []byte {
	return packet.payload
}

func (packet *PublishPacket) GetIdentifier() string {
	return packet.Key.ToString()
}

func (packet *PublishPacket) SetNewOwner(newOwner int64) {
	packet.owner = newOwner
}

func (packet *PublishPacket) GetSender() int64 {
	return packet.owner
}

func (packet *PublishPacket) Marshal() ([]byte, errorinfo.GatewayError) {
	return []byte{}, nil
}
