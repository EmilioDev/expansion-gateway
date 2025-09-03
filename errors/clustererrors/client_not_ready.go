package clustererrors

import (
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
)

type ClientNotReady struct {
	ClusterError
}

func (err ClientNotReady) Error() string {
	return fmt.Sprintf("Cluster client (Kind: %d) is not ready yet to operate.", err.TargetKind)
}

func (err ClientNotReady) SetStackTrace(stackTrace []string) errorinfo.GatewayError {
	err.StackTrace = stackTrace
	return &err
}

func CreateClientNotReadyError(file string, index uint16, kindOfServer enums.ClusterMember_Kind) ClientNotReady {
	return ClientNotReady{
		CreateClusterError(file, "client not ready yet to operate", index, 15, kindOfServer, false),
	}
}
