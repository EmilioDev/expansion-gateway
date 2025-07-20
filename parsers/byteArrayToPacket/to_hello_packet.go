package bytearraytopacket

import (
	"expansion-gateway/dto"
	"expansion-gateway/enums"
	"expansion-gateway/errors"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/packets"
)

func ToHelloPacket(byteArray *[]byte, connectionID int64) (packets.Packet, errors.GatewayError) {
	byteArraySize := len(*byteArray)
	const filePath string = "/parsers/byteArrayToPacket/to_hello_packet.go"

	if byteArraySize < 5 || byteArraySize > 14 {
		return nil, errors.CreateInvalidPacketSizeError(filePath, 15, enums.HELLO, byteArraySize)
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
		return nil, errors.CreatePacketWithInvalidProtocolVersion(filePath, 26, enums.HELLO, currentByte)
	}

	//client type
	currentByte = (*byteArray)[2]

	if enums.IsValidClientType(currentByte) {
		answer.VariableHeader.ClientType = enums.ClientType(currentByte)
	} else {
		return nil, errors.CreatePacketWithInvalidClientType(filePath, 35, enums.HELLO, currentByte)
	}

	// client version
	currentByte = (*byteArray)[3]

	if helpers.ValidateClientVersion(answer.VariableHeader.ClientType, currentByte) {
		answer.VariableHeader.ClientVersion = currentByte
	} else {
		return nil, errors.CreatePacketWithInvalidClientVersion(filePath, 45, enums.HELLO, answer.VariableHeader.ClientType, currentByte)
	}

	// packet flags
	currentByte = (*byteArray)[4]

	if !checkFlags(currentByte) {
		return nil, errors.CreatePacketWithInvalidFlags(filePath, 54, enums.HELLO, answer.VariableHeader.ClientType, currentByte)
	}

	// from right to left:
	// the 1st bit -> payload encrypted
	// the 2nd bit -> requesting session resume

	answer.VariableHeader.PayloadEncrypted = (currentByte & 0x01) == 0x01
	answer.VariableHeader.PayloadEncrypted = (currentByte & 0x02) == 0x02

	// encryptation algorythm
	if answer.VariableHeader.PayloadEncrypted && byteArraySize < 6 {
		return nil, errors.CreateInvalidPacketSizeError(filePath, 65, enums.HELLO, byteArraySize)
	}

	currentByte = (*byteArray)[5]

	if !enums.IsValidEncryptationAlgorythm(currentByte) {
		return nil, errors.CreatePacketWithInvalidEncryptationAlgorythm(filePath, 72, enums.HELLO, currentByte)
	}

	answer.VariableHeader.Encryptation = enums.EncryptationAlgorythm(currentByte)

	return &answer, nil
}

func checkFlags(flagsByte byte) bool {
	return flagsByte <= 3
}
