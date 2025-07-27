package dto

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

type InvalidPacket struct {
	Sender int64
}

func (packet InvalidPacket) GetPacketType() enums.PacketType {
	return enums.NONE
}

func (packet InvalidPacket) GetPacketHeader() packets.PacketHeader {
	return nil
}

func (packet InvalidPacket) GetPayload() string {
	return ""
}

func (packet InvalidPacket) Marshal() ([]byte, errorinfo.GatewayError) {
	return []byte{}, nil
}

func (packet InvalidPacket) GetSender() int64 {
	return packet.Sender
}

// creates an invalid packet
func CreateInvalidPacket(sender int64) InvalidPacket {
	return InvalidPacket{
		sender,
	}
}
