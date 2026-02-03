package bytearraytopacket

import (
	"expansion-gateway/dto"
	"expansion-gateway/enums"
	"expansion-gateway/errors"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
	"expansion-gateway/internal/structs/tries"
)

func ToSubscribePacket(byteArray *[]byte, connectionID int64) (packets.Packet, errorinfo.GatewayError) {
	packetSize := len(*byteArray)
	input := *byteArray
	const filePath string = "/parsers/byteArrayToPacket/to_subscribe_packet.go"

	if packetSize < 9 {
		return nil, errors.CreateInvalidPacketSizeError(filePath, 16, enums.SUBSCRIBE, packetSize)
	}

	packetId := helpers.Convert4bytesIntoInt32([4]byte{input[1], input[2], input[3], input[4]})
	subLen := helpers.Convert4bytesIntoInt32([4]byte{input[5], input[6], input[7], input[8]})
	keyBytes := 9 + subLen

	if keyBytes != int32(packetSize) {
		return nil, errors.CreateInvalidPacketSizeError(filePath, 23, enums.SUBSCRIBE, packetSize)
	}

	rawKey := string(input[9:keyBytes])

	if key, err := tries.ConvertStringToSubscriptionKey(rawKey); err == nil {
		return dto.CreateSubscribePacket(key, packetId, connectionID), nil
	} else {
		return nil, err
	}
}
