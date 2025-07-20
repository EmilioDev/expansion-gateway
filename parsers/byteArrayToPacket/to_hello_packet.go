package bytearraytopacket

import (
	"expansion-gateway/dto"
	"expansion-gateway/enums"
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/packets"
)

func ToHelloPacket(byteArray *[]byte, connectionID int64) (packets.Packet, errors.GatewayError) {
	byteArraySize := len(*byteArray)
	const filePath string = "/parsers/byteArrayToPacket/to_hello_packet.go"

	if byteArraySize < 5 || byteArraySize > 14 {
		return nil, errors.CreateInvalidPacketSizeError(filePath, 14, enums.HELLO, byteArraySize)
	}

	answer := dto.HelloPacket{
		Sender:         connectionID,
		VariableHeader: &dto.HelloPacketVariableHeader{},
	}

	// check protocol version
	currentByte := (*byteArray)[1]

	if enums.IsValidProtocolVersion(currentByte) {
		answer.VariableHeader.ProtocolVersion = enums.ProtocolVersion(currentByte)
	} else {
		return nil, errors.CreatePacketWithInvalidProtocolVersion(filePath, 24, enums.HELLO, currentByte)
	}

	//client type
	currentByte = (*byteArray)[2]

	if enums.IsValidClientType(currentByte) {
		answer.VariableHeader.ClientType = enums.ClientType(currentByte)
	} else {
		//
	}

	return &answer, nil
}
