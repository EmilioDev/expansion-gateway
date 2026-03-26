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
	counter         uint64
	useWindow       bool
	payload         []byte
}

func CreatePublishPacket(
	key tries.SubscriptionKey,
	packetId int32,
	owner int64,
	payload []byte,
	counter uint64,
	useWindow bool,
) *PublishPacket {
	return &PublishPacket{
		Key:             key,
		publishPacketId: packetId,
		owner:           owner,
		payload:         payload,
		counter:         counter,
		useWindow:       useWindow,
	}
}

func CreateBasicPublishPacket(key tries.SubscriptionKey, payload []byte) *PublishPacket {
	return &PublishPacket{
		Key:             key,
		publishPacketId: 0,
		owner:           0,
		payload:         payload,
		counter:         0,
		useWindow:       true,
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
	packetLen := keyLen + payloadLen + 1 + 1 + 4 + 4

	var flag byte = 0

	if packet.publishPacketId != 0 {
		packetLen += 4
		flag = 1
	}

	if packet.counter != 0 {
		packetLen += 8
		flag += 2
	}

	if packet.useWindow {
		flag += 4
	}

	result := make([]byte, 0, packetLen)

	// publish identifier
	result = append(result, byte(enums.PUBLISH))

	// flags
	result = append(result, flag)

	// packet identifier
	if packet.publishPacketId != 0 {
		identifier := helpers.ConvertInt32Into4Bytes(packet.publishPacketId)
		result = append(result, identifier[:]...)
	}

	// packet counter
	if packet.counter != 0 {
		counter := helpers.ConvertUInt64Into8Bytes(packet.counter)
		result = append(result, counter[:]...)
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
