// file: /dto/challenge_packet.go
package dto

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

type ChallengePacket struct {
	BaseHeader
	UserId    int64
	Challenge []byte
}

func (packet ChallengePacket) GetPacketType() enums.PacketType {
	return enums.CHALLENGE
}

func (packet ChallengePacket) GetPacketHeader() packets.PacketHeader {
	return nil
}

func (packet ChallengePacket) GetPayload() string {
	return ""
}

func (packet ChallengePacket) Marshal() ([]byte, errorinfo.GatewayError) {
	output := make([]byte, 1+len(packet.Challenge))

	output[0] = byte(enums.CHALLENGE)
	copy(output[1:], packet.Challenge)

	return output, nil
}

func (packet ChallengePacket) GetSender() int64 {
	return packet.UserId
}

func GenerateChallengePacket(userId int64, challenge *[]byte) *ChallengePacket {
	return &ChallengePacket{
		BaseHeader{},
		userId,
		*challenge,
	}
}
