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
	bytes := helpers.ConvertInt64Into8Bytes(packet.Subscription.SubscriptionID)
	pathSizeAsInt32 := int32(len(packet.Subscription.GatewayPath))
	pathSize := helpers.ConvertInt32Into4Bytes(pathSizeAsInt32)
	pathLimit := 45 + int(pathSizeAsInt32)

	if len(packet.Subscription.SessionEphemeralKey) == 32 { // if you have public password
		// easy to put 77 here, but this way you know where each number comes from
		output := make([]byte, 0, 1+8+32+4+len(packet.Subscription.GatewayPath)+32)

		output[0] = byte(enums.REDIRECT)

		copy(output[1:9], bytes[:])                                         // subscription id
		copy(output[9:41], packet.Subscription.Challenge)                   // challenge
		copy(output[41:45], pathSize[:])                                    // gateway path length
		copy(output[45:pathLimit], []byte(packet.Subscription.GatewayPath)) // gateway path
		copy(output[pathLimit:], packet.Subscription.SessionEphemeralKey)   // session ephemeral public key

		return output, nil
	} else { // if not
		// easy to put 45 here, but this way you know where each number comes from
		output := make([]byte, 0, 1+8+32+4+len(packet.Subscription.GatewayPath))

		output[0] = byte(enums.REDIRECT)

		copy(output[1:9], bytes[:])                                         // subscription id
		copy(output[9:41], packet.Subscription.Challenge)                   // challenge
		copy(output[41:45], pathSize[:])                                    // gateway path length
		copy(output[45:pathLimit], []byte(packet.Subscription.GatewayPath)) // gateway path

		return output, nil
	}
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
