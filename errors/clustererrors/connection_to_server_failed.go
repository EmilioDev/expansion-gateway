package clustererrors

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type ConnectionToServerFailedError struct {
	ClusterError
	Server string
}

func (err ConnectionToServerFailedError) Error() string {
	return fmt.Sprintf("Failed to connect to server %s of kind %d", err.Server, err.TargetKind)
}

func (err ConnectionToServerFailedError) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}

func CreateConnectionToServerFailedError(file string, index uint16, serverPath string, kindOfServer enums.ClusterMember_Kind, isServer bool) ConnectionToServerFailedError {
	return ConnectionToServerFailedError{
		CreateClusterError(file, "connection to cluster member failed", index, 14, kindOfServer, isServer),
		serverPath,
	}
}
