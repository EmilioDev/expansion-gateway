package errors

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type InvalidPacketSize struct {
	PacketError
	CurrentSize int
}

func (err InvalidPacketSize) Error() string {
	return fmt.Sprintf("Packet of type %s with invalid size of %d", enums.GetNameOfPacketType(err.PacketType), err.CurrentSize)
}

// Creates an invalid size error
func CreateInvalidPacketSizeError(file string, line uint16, packetType enums.PacketType, currentSize int) InvalidPacketSize {
	return InvalidPacketSize{
		CreatePacketError(file, "Packet with invalid size", line, 1, packetType),
		currentSize,
	}
}

func (err InvalidPacketSize) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}
