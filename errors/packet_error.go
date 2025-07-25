package errors

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type PacketError struct {
	BaseError
	PacketType enums.PacketType
}

func (err PacketError) Error() string {
	return fmt.Sprintf("%s packet type: %s", err.BaseError.Error(), enums.GetNameOfPacketType(err.PacketType))
}

func CreatePacketError(file, reason string, line uint16, errorCode byte, packetType enums.PacketType) PacketError {
	return PacketError{
		CreateBaseError(file, reason, line, errorCode),
		packetType,
	}
}

func (err PacketError) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}
