// file: /dto/challenge_packet.go
package dto

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

type ChallengePacket struct {
	BaseHeader
	UserId                   int64
	Challenge                []byte
	ServerPublicEphemeralKey *[32]byte
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
	var output []byte = nil

	if packet.ServerPublicEphemeralKey == nil {
		output = make([]byte, 0, 1+32) // the challenge always has 32 bytes length

		output = append(output, byte(enums.CHALLENGE)) // we mark the stream as a CHALLENGE packet
		output = append(output, packet.Challenge...)   // we set the challenge and that's it
	} else {
		output = make([]byte, 0, 1+32+32)             // the challenge always has 32 bytes length and the key also has 32 bytes length
		serverKey := *packet.ServerPublicEphemeralKey // the server public key, it has 32 bytes length

		output = append(output, byte(enums.CHALLENGE)) // we mark the stream as a CHALLENGE packet
		output = append(output, packet.Challenge...)   // we set the challenge
		output = append(output, serverKey[:]...)       // and we set the server public ephemeral key, and done.
	}

	return output, nil
}

func (packet ChallengePacket) GetSender() int64 {
	return packet.UserId
}

// generates a challenge without the the ephemeral key
func GenerateChallengePacket(userId int64, challenge *[]byte) *ChallengePacket {
	return &ChallengePacket{
		BaseHeader:               BaseHeader{},
		UserId:                   userId,
		Challenge:                *challenge,
		ServerPublicEphemeralKey: nil,
	}
}

// generates a challenge with the the server ephemeral key
func GenerateChallengePacketWithServerPublicEphemeralKey(userId int64, challenge *[]byte, serverPublicKey [32]byte) *ChallengePacket {
	return &ChallengePacket{
		BaseHeader:               BaseHeader{},
		UserId:                   userId,
		Challenge:                *challenge,
		ServerPublicEphemeralKey: &serverPublicKey,
	}
}
