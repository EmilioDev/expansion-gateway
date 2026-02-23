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

func ToPublishPacket(byteArray *[]byte, connectionID int64) (packets.Packet, errorinfo.GatewayError) {
	input := *byteArray
	streamSize := len(input)
	const filePath string = "/parsers/byteArrayToPacket/to_publish_packet.go"

	var packetId int32 = 0

	if streamSize < 6 {
		return nil, errors.CreateInvalidPacketSizeError(filePath, 18, enums.PUBLISH, streamSize)
	}

	index := 1

	// packet id
	if input[index] == 1 {
		idArray := [4]byte{input[index+1], input[index+2], input[index+3], input[index+4]}
		packetId = helpers.Convert4bytesIntoInt32(idArray)
		index += 5
	}

	// key
	if index+3 >= streamSize {
		return nil, errors.CreateInvalidPacketSizeError(filePath, 33, enums.PUBLISH, streamSize)
	}

	keySizeArray := [4]byte{input[index], input[index+1], input[index+2], input[index+3]}
	index += 4
	keySize := int(helpers.Convert4bytesIntoInt32(keySizeArray))

	if index+keySize >= streamSize {
		return nil, errors.CreateInvalidPacketSizeError(filePath, 41, enums.PUBLISH, streamSize)
	}

	var key tries.SubscriptionKey

	if k, err := tries.ConvertStringToSubscriptionKey(string(input[index : index+keySize])); err == nil {
		key = k
	} else {
		return nil, err
	}

	index += keySize

	// payload
	if index+3 >= streamSize {
		return nil, errors.CreateInvalidPacketSizeError(filePath, 56, enums.PUBLISH, streamSize)
	}

	payloadSizeArray := [4]byte{input[index], input[index+1], input[index+2], input[index+3]}
	payloadSize := int(helpers.Convert4bytesIntoInt32(payloadSizeArray))
	index += 4

	if payloadSize+index != streamSize {
		return nil, errors.CreateInvalidPacketSizeError(filePath, 64, enums.PUBLISH, streamSize)
	}

	payload := input[index : index+payloadSize]

	// return
	return dto.CreatePublishPacket(key, packetId, connectionID, payload), nil
}
