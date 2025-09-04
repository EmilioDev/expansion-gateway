package impl

import (
	"expansion-gateway/dto"
	"expansion-gateway/interfaces/errorinfo"
)

type ClusterLeaderSubscribeCallback func(string) (*dto.ClusterMemberSubscriptionResult, errorinfo.GatewayError)
type ClusterLeaderUnsubscribeCallback func(int64) (bool, errorinfo.GatewayError)
type ClusterLeaderHealthCheckCallback func(int64, int64, int64, int32, bool) (bool, errorinfo.GatewayError)
