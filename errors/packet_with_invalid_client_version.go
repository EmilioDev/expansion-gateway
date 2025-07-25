package errors

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type PacketWithInvalidClientVersion struct {
	PacketError
	RequestedVersion byte
	ClientType       enums.ClientType
}

func (err PacketWithInvalidClientVersion) Error() string {
	return fmt.Sprintf("Packet of type %s and client type of %d with invalid version of %d", enums.GetNameOfPacketType(err.PacketType), byte(err.ClientType), err.RequestedVersion)
}

func CreatePacketWithInvalidClientVersion(file string, line uint16, packetType enums.PacketType, clientType enums.ClientType, requestedVersion byte) PacketWithInvalidClientVersion {
	return PacketWithInvalidClientVersion{
		CreatePacketError(file, "Invalid client version", line, 4, packetType),
		requestedVersion,
		clientType,
	}
}

func (err PacketWithInvalidClientVersion) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}
