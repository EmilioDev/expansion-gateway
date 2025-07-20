package errors

import (
	"expansion-gateway/enums"
	"fmt"
)

type PacketWithInvalidProtocolVersion struct {
	PacketError
	RequestedProtocolVersion byte
}

func (err PacketWithInvalidProtocolVersion) Error() string {
	return fmt.Sprintf("Packet of type %s with invalid protocol version of %d", enums.GetNameOfPacketType(err.PacketType), err.RequestedProtocolVersion)
}

func CreatePacketWithInvalidProtocolVersion(file string, line uint16, packetType enums.PacketType, requestedProtocolVersion byte) PacketWithInvalidProtocolVersion {
	return PacketWithInvalidProtocolVersion{
		CreatePacketError(file, "Invalid protocol version", line, 2, packetType),
		requestedProtocolVersion,
	}
}
