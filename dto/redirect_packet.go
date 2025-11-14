package dto

import (
	"expansion-gateway/dto/clusters/results"
	"expansion-gateway/enums"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

type RedirectPacket struct {
	UserId       int64
	Subscription *results.ClusterUserSubscriptionResult
}

func (packet *RedirectPacket) GetPacketType() enums.PacketType {
	return enums.REDIRECT
}

func (packet *RedirectPacket) GetPacketHeader() packets.PacketHeader {
	return nil
}

func (packet *RedirectPacket) GetPayload() string {
	return ""
}

func (packet *RedirectPacket) Marshal() ([]byte, errorinfo.GatewayError) {
	output := make([]byte, 0, 1+8+32+len(packet.Subscription.GatewayPath))

	output[0] = byte(enums.REDIRECT)

	bytes := helpers.ConvertInt64Into8Bytes(packet.Subscription.SubscriptionID)

	copy(output[1:], bytes[:])                                 // subscription id
	copy(output[9:], packet.Subscription.Challenge)            // challenge
	copy(output[41:], []byte(packet.Subscription.GatewayPath)) // gateway path

	return output, nil
}

func (packet *RedirectPacket) GetSender() int64 {
	return packet.UserId
}

func CreateNewRedirectPacket(userId int64, subscription *results.ClusterUserSubscriptionResult) *RedirectPacket {
	return &RedirectPacket{
		UserId:       userId,
		Subscription: subscription,
	}
}
