package dto

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

type PingAckPacket struct {
	sessionID int64
}

func CreatePingACKpacket(sessionId int64) *PingAckPacket {
	return &PingAckPacket{
		sessionID: sessionId,
	}
}

func (packet *PingAckPacket) GetPacketType() enums.PacketType {
	return enums.PINGACK
}

func (packet *PingAckPacket) GetPacketHeader() packets.PacketHeader {
	return nil
}

func (packet *PingAckPacket) GetPayload() string {
	return ""
}

func (packet *PingAckPacket) GetSender() int64 {
	return packet.sessionID
}

func (packet *PingAckPacket) Marshal() ([]byte, errorinfo.GatewayError) {
	return []byte{byte(enums.PINGACK)}, nil
}
