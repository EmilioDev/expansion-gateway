package dto

import (
	"expansion-gateway/enums"
	"expansion-gateway/helpers"
	errors "expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
	"fmt"
)

type PubackPacket struct {
	owner     int64
	publishId int32
	result    enums.PublishResult
}

func CreatePubackPacket(owner int64, publishToAcknowledge int32, result enums.PublishResult) *PubackPacket {
	return &PubackPacket{
		owner:     owner,
		publishId: publishToAcknowledge,
		result:    result,
	}
}

func (packet *PubackPacket) GetPacketType() enums.PacketType {
	return enums.PUBACK
}

func (packet *PubackPacket) GetPacketHeader() packets.PacketHeader {
	return nil
}

func (packet *PubackPacket) GetPayload() string {
	return ""
}

func (packet *PubackPacket) GetSender() int64 {
	return packet.owner
}

func (packet *PubackPacket) GetIdentifier() string {
	return fmt.Sprintf("%d", packet.publishId)
}

func (packet *PubackPacket) SetNewOwner(newOwner int64) {
	packet.owner = newOwner
}

func (packet *PubackPacket) Marshal() ([]byte, errors.GatewayError) {
	output := make([]byte, 0, 1+4+1)

	subId := helpers.ConvertInt32Into4Bytes(packet.publishId)

	output = append(output, byte(enums.PUBACK))
	output = append(output, subId[:]...)

	output = append(output, byte(packet.result))

	return output, nil
}

func (packet *PubackPacket) GetRawPayload() []byte {
	return []byte{}
}
