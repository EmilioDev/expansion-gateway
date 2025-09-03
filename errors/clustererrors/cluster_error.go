package clustererrors

import (
	"expansion-gateway/enums"
	"expansion-gateway/errors"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type ClusterError struct {
	errors.BaseError
	TargetKind enums.ClusterMember_Kind
	IsServer   bool
}

func (err ClusterError) Error() string {
	return fmt.Sprintf("Cluster member (Kind: %d) (Is Server: %t) had an error.", err.TargetKind, err.IsServer)
}

func (err ClusterError) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}

func CreateClusterError(file, description string, index uint16, errorCode byte, kindOfServer enums.ClusterMember_Kind, isServer bool) ClusterError {
	return ClusterError{
		errors.CreateBaseError(file, description, index, errorCode),
		kindOfServer,
		isServer,
	}
}
