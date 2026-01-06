package errors

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type InvalidDisconnectReason struct {
	PacketError
	disconnect_reason byte
}

func (err *InvalidDisconnectReason) Error() string {
	return fmt.Sprintf("Disconnect packet has invalid disconnect reason of %d", err.disconnect_reason)
}

// Creates an invalid size error
func CreateInvalidDisconnectReasonError(file string, line uint16, reason byte) *InvalidDisconnectReason {
	return &InvalidDisconnectReason{
		CreatePacketError(file, "Disconnect packet with invalid reason", line, 23, enums.DISCONNECT),
		reason,
	}
}

func (err *InvalidDisconnectReason) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return err
}
