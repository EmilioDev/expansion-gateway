package dto

import "expansion-gateway/enums"

type Packet struct {
	Sender     int64
	PacketType enums.PacketType
	Payload    string
}

func CreateDefaultPacket(sender int64) *Packet {
	return &Packet{
		Sender:     sender,
		PacketType: enums.HELLO,
		Payload:    "",
	}
}
