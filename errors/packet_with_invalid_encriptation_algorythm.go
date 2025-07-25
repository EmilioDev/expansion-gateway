package errors

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type PacketWithInvalidEncryptionAlgorythm struct {
	PacketError
	RequestedEncryptAlgorythm byte
}

func (err PacketWithInvalidEncryptionAlgorythm) Error() string {
	return fmt.Sprintf("Packet of type %s with invalid encryption algorythm of %d", enums.GetNameOfPacketType(err.PacketType), err.RequestedEncryptAlgorythm)
}

func CreatePacketWithInvalidEncryptionAlgorythm(file string, line uint16, packetType enums.PacketType, requestedEncryptionAlgorythm byte) PacketWithInvalidEncryptionAlgorythm {
	return PacketWithInvalidEncryptionAlgorythm{
		CreatePacketError(file, "Invalid client type", line, 6, packetType),
		requestedEncryptionAlgorythm,
	}
}

func (err PacketWithInvalidEncryptionAlgorythm) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}
