package errors

import (
	"expansion-gateway/enums"
	"fmt"
)

type PacketWithInvalidEncryptationAlgorythm struct {
	PacketError
	RequestedEncryptAlgorythm byte
}

func (err PacketWithInvalidEncryptationAlgorythm) Error() string {
	return fmt.Sprintf("Packet of type %s with invalid encryption algorythm of %d", enums.GetNameOfPacketType(err.PacketType), err.RequestedEncryptAlgorythm)
}

func CreatePacketWithInvalidEncryptationAlgorythm(file string, line uint16, packetType enums.PacketType, requestedEncryptionAlgorythm byte) PacketWithInvalidEncryptationAlgorythm {
	return PacketWithInvalidEncryptationAlgorythm{
		CreatePacketError(file, "Invalid client type", line, 3, packetType),
		requestedEncryptionAlgorythm,
	}
}
