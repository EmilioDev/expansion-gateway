package parsers

import (
	"expansion-gateway/dto"
	"expansion-gateway/errors"
)

func ParseByteArrayToPacket(byteArray *[]byte, connectionID int64) (*dto.Packet, errors.GatewayError) {
	answer := dto.CreateDefaultPacket(connectionID)

	return answer, nil
}
