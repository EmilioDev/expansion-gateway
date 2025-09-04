package impl

import (
	"expansion-gateway/dto"
	"expansion-gateway/interfaces/errorinfo"
)

type ClusterLeaderSubscribeCallback func(string) (*dto.ClusterMemberSubscriptionResult, errorinfo.GatewayError)
