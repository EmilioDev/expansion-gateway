package dto

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

type DisconnectPacket struct {
	sessionID int64
	reason    enums.DisconnectReason
}

func CreateDisconnectPacket(sessionId int64, reason enums.DisconnectReason) *DisconnectPacket {
	return &DisconnectPacket{
		sessionID: sessionId,
		reason:    reason,
	}
}

func (packet *DisconnectPacket) GetDisconnectReason() enums.DisconnectReason {
	return packet.reason
}

func (packet *DisconnectPacket) GetPacketType() enums.PacketType {
	return enums.DISCONNECT
}

func (packet *DisconnectPacket) GetPacketHeader() packets.PacketHeader {
	return nil
}

func (packet *DisconnectPacket) GetPayload() string {
	return ""
}

func (packet *DisconnectPacket) GetSender() int64 {
	return packet.sessionID
}

func (packet *DisconnectPacket) Marshal() ([]byte, errorinfo.GatewayError) {
	return []byte{byte(enums.DISCONNECT), byte(packet.reason)}, nil
}

func (packet *DisconnectPacket) GetRawPayload() []byte {
	return []byte{}
}

func (packet *DisconnectPacket) GetIdentifier() string {
	return ""
}

func (packet *DisconnectPacket) SetNewOwner(newOwner int64) {
	packet.sessionID = newOwner
}
