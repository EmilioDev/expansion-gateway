package errors

import (
	"expansion-gateway/enums"
	"fmt"
)

type PacketWithInvalidFlags struct {
	PacketError
	FlagsByte byte
}

func (err PacketWithInvalidFlags) Error() string {
	return fmt.Sprintf("Packet of type %s with invalid flags of %d", enums.GetNameOfPacketType(err.PacketType), err.FlagsByte)
}

func CreatePacketWithInvalidFlags(file string, line uint16, packetType enums.PacketType, clientType enums.ClientType, flagsByte byte) PacketWithInvalidFlags {
	return PacketWithInvalidFlags{
		CreatePacketError(file, "Invalid flags", line, 5, packetType),
		flagsByte,
	}
}
