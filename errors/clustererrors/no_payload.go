package clustererrors

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type NoPayloadError struct {
	ClusterError
}

func (err NoPayloadError) Error() string {
	return fmt.Sprintf("Cluster client (Kind: %d) is not ready yet to operate.", err.TargetKind)
}

func (err NoPayloadError) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}

func CreateNoPayloadError(file string, index uint16, kindOfServer enums.ClusterMember_Kind, isServer bool) NoPayloadError {
	return NoPayloadError{
		CreateClusterError(file, "the response had no payload to read", index, 16, kindOfServer, isServer),
	}
}
