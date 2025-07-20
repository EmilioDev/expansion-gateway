package errors

import (
	"expansion-gateway/enums"
	"fmt"
)

type PacketWithInvalidClientType struct {
	PacketError
	SpecifiedClientType byte // the type of client that was specified
}

func (err PacketWithInvalidClientType) Error() string {
	return fmt.Sprintf("Packet of type %s with invalid client type of %d", enums.GetNameOfPacketType(err.PacketType), err.SpecifiedClientType)
}

func CreatePacketWithInvalidClientType(file string, line uint16, packetType enums.PacketType, requestedClientType byte) PacketWithInvalidClientType {
	return PacketWithInvalidClientType{
		CreatePacketError(file, "Invalid client type", line, 3, packetType),
		requestedClientType,
	}
}
