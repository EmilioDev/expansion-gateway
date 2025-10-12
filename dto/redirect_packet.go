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

func (packet RedirectPacket) GetPacketType() enums.PacketType {
	return enums.REDIRECT
}

func (packet RedirectPacket) GetPacketHeader() packets.PacketHeader {
	return nil
}

func (packet RedirectPacket) GetPayload() string {
	return ""
}

func (packet RedirectPacket) Marshal() ([]byte, errorinfo.GatewayError) {
	output := make([]byte, 0, 1+8+8+len(packet.Subscription.Challenge))

	bytes := helpers.ConvertInt64Into8Bytes(packet.UserId)

	output[0] = byte(enums.REDIRECT)
	copy(output[1:], bytes[:])

	bytes = helpers.ConvertInt64Into8Bytes(packet.Subscription.SubscriptionID)
	copy(output[9:], bytes[:])
	copy(output[17:], packet.Subscription.Challenge)

	return output, nil
}

func (packet RedirectPacket) GetSender() int64 {
	return packet.UserId
}

func CreateNewRedirectPacket(userId int64, subscription *results.ClusterUserSubscriptionResult) *RedirectPacket {
	return &RedirectPacket{
		UserId:       userId,
		Subscription: subscription,
	}
}
