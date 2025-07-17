package errors

import "expansion-gateway/enums"

type InvalidPacket struct {
	PacketError
}

func (err InvalidPacket) Error() string {
	return "Invalid packet"
}

func CreateInvalidPacketError(file string, line uint16) InvalidPacket {
	return InvalidPacket{
		CreatePacketError(file, "Invalid packet", line, 0, enums.NONE),
	}
}
