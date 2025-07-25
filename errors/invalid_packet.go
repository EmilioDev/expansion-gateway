package errors

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
)

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

func (err InvalidPacket) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}
