package bytearraytopacket

import (
	"expansion-gateway/dto"
	"expansion-gateway/enums"
	"expansion-gateway/errors"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
)

// creates a redirected packet from a byte stream
func ToRedirectedPacket(byteArray *[]byte, connectionID int64) (packets.Packet, errorinfo.GatewayError) {
	input := *byteArray
	streamSize := len(input)
	const filePath string = "/parsers/byteArrayToPacket/to_redirected_packet.go"

	switch streamSize {
	case 73:
		subscriptionId := helpers.ConvertBytesArrayIntoSingleInt64(input[1:9])
		signature := [64]byte{}

		copy(signature[:], input[9:])

		return dto.CreateRedirectedPacketWithoutEphemeralKey(connectionID, signature, subscriptionId), nil

	case 105:
		subscriptionId := helpers.ConvertBytesArrayIntoSingleInt64(input[1:9])
		signature := [64]byte{}
		ephemeralPublicKey := [32]byte{}

		copy(signature[:], input[9:73])
		copy(ephemeralPublicKey[:], input[73:])

		return dto.CreateRedirectedPacketWithEphemeralKey(connectionID, signature, subscriptionId, ephemeralPublicKey), nil

	default:
		return nil, errors.CreateInvalidPacketSizeError(filePath, 22, enums.REDIRECTED, streamSize)
	}
}
