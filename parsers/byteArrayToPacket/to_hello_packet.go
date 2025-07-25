package bytearraytopacket

import (
	"expansion-gateway/dto"
	"expansion-gateway/enums"
	"expansion-gateway/errors"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

func ToHelloPacket(byteArray *[]byte, connectionID int64) (packets.Packet, errorinfo.GatewayError) {
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
	answer.VariableHeader.SessionResume = (currentByte & 0x02) == 0x02

	if !answer.VariableHeader.PayloadEncrypted && !answer.VariableHeader.SessionResume && byteArraySize != 5 {
		return nil, errors.CreateInvalidPacketSizeError(filePath, 65, enums.HELLO, byteArraySize)
	}

	// encryptation algorythm
	index := 5
	if answer.VariableHeader.PayloadEncrypted {
		if byteArraySize <= index {
			return nil, errors.CreateInvalidPacketSizeError(filePath, 72, enums.HELLO, byteArraySize)
		}

		currentByte = (*byteArray)[index]

		if !enums.IsValidEncryptionAlgorythm(currentByte) {
			return nil, errors.CreatePacketWithInvalidEncryptionAlgorythm(filePath, 78, enums.HELLO, currentByte)
		}

		answer.VariableHeader.Encryptation = enums.EncryptionAlgorithm(currentByte)
		index++
	}

	// requested session id to resume
	if answer.VariableHeader.SessionResume {
		if byteArraySize <= index+7 {
			return nil, errors.CreateInvalidPacketSizeError(filePath, 87, enums.HELLO, byteArraySize)
		}

		answer.VariableHeader.PretendedUserID = helpers.Convert8BytesIntoSingleInt64(
			(*byteArray)[index],
			(*byteArray)[index+1],
			(*byteArray)[index+2],
			(*byteArray)[index+3],
			(*byteArray)[index+4],
			(*byteArray)[index+5],
			(*byteArray)[index+6],
			(*byteArray)[index+7])

		//index+=8 // not needed if this is the last thing to do
	}

	return &answer, nil
}

func checkFlags(flagsByte byte) bool {
	return flagsByte <= 3
}
