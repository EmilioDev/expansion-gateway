package clustererrors

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type OperationFailedError struct {
	ClusterError
}

func (err OperationFailedError) Error() string {
	return fmt.Sprintf("Cluster client (Kind: %d) is not ready yet to operate.", err.TargetKind)
}

func (err OperationFailedError) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}

func CreateOperationFailedError(file string, index uint16, kindOfServer enums.ClusterMember_Kind, isServer bool) OperationFailedError {
	return OperationFailedError{
		CreateClusterError(file, "the requested operation has failed", index, 17, kindOfServer, isServer),
	}
}
