package dto

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

type PingPacket struct {
	sessionID int64
}

func CreatePingPacket(sessionId int64) *PingPacket {
	return &PingPacket{
		sessionID: sessionId,
	}
}

func (packet *PingPacket) GetPacketType() enums.PacketType {
	return enums.PING
}

func (packet *PingPacket) GetPacketHeader() packets.PacketHeader {
	return nil
}

func (packet *PingPacket) GetPayload() string {
	return ""
}

func (packet *PingPacket) GetSender() int64 {
	return packet.sessionID
}

func (packet *PingPacket) Marshal() ([]byte, errorinfo.GatewayError) {
	return []byte{byte(enums.PING)}, nil
}
