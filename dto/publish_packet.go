package dto

import (
	"expansion-gateway/enums"
	"expansion-gateway/helpers"
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
	keyLen := packet.Key.KeyLength()
	payloadLen := len(packet.payload)
	result := make([]byte, 0, 1+1+4+keyLen+4+payloadLen)

	// publish identifier
	result = append(result, byte(enums.PUBLISH))

	// flags
	if packet.publishPacketId == 0 {
		result = append(result, 0)
	} else {
		identifier := helpers.ConvertInt32Into4Bytes(packet.publishPacketId)

		result = append(result, 1)
		result = append(result, identifier[:]...)
	}

	// key
	keyLenArray := helpers.ConvertInt32Into4Bytes(int32(keyLen))
	key := packet.Key.ToByteArray()

	result = append(result, keyLenArray[:]...)
	result = append(result, key...)

	// payload
	payloadLenArray := helpers.ConvertInt32Into4Bytes(int32(payloadLen))

	result = append(result, payloadLenArray[:]...)
	result = append(result, packet.payload...)

	return result, nil
}

func (packet *PublishPacket) NeedsAcknowledgement() bool {
	return packet.publishPacketId != 0
}
