package dto

import (
	"expansion-gateway/enums"
	"expansion-gateway/helpers"
	errors "expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

type RedirectedPacket struct {
	OwnerID                  int64    // the id of the session of the user that owns this packet
	Signature                [64]byte // the signature of the data the user is going to send
	SubscriptionID           int64    // the id of the subscription this packet is validating
	ClientPublicEphemeralKey []byte   // the public ephemeral key of the client, used to generate the final key
}

func (packet *RedirectedPacket) GetPacketType() enums.PacketType {
	return enums.REDIRECTED
}

func (packet *RedirectedPacket) GetPacketHeader() packets.PacketHeader {
	return nil
}

func (packet *RedirectedPacket) GetPayload() string {
	return ""
}

func (packet *RedirectedPacket) Marshal() ([]byte, errors.GatewayError) {
	if len(packet.ClientPublicEphemeralKey) == 32 {
		answer := make([]byte, 0, 1+8+64+32) //yep, you can put directly 105, but with this, conbined with above, you can understand better

		signature := packet.Signature[:]
		subscription := helpers.ConvertInt64Into8Bytes(packet.SubscriptionID)

		answer = append(answer, byte(enums.REDIRECTED))
		answer = append(answer, subscription[:]...)
		answer = append(answer, signature...)
		answer = append(answer, packet.ClientPublicEphemeralKey...)

		return answer, nil
	} else {
		answer := make([]byte, 0, 1+8+64) //yep, you can put directly 73, but with this, conbined with above, you can understand better

		signature := packet.Signature[:]
		subscription := helpers.ConvertInt64Into8Bytes(packet.SubscriptionID)

		answer = append(answer, byte(enums.REDIRECTED))
		answer = append(answer, subscription[:]...)
		answer = append(answer, signature...)

		return answer, nil
	}
}

func (packet *RedirectedPacket) GetSender() int64 {
	return packet.OwnerID
}

func CreateRedirectedPacketWithoutEphemeralKey(
	owner int64,
	signature [64]byte,
	subscriptionId int64) *RedirectedPacket {
	return &RedirectedPacket{
		OwnerID:                  owner,
		Signature:                signature,
		SubscriptionID:           subscriptionId,
		ClientPublicEphemeralKey: nil,
	}
}

func CreateRedirectedPacketWithEphemeralKey(
	owner int64,
	signature [64]byte,
	subscriptionId int64,
	clientPublicEphemeralKey [32]byte) *RedirectedPacket {
	return &RedirectedPacket{
		OwnerID:                  owner,
		Signature:                signature,
		SubscriptionID:           subscriptionId,
		ClientPublicEphemeralKey: clientPublicEphemeralKey[:],
	}
}
