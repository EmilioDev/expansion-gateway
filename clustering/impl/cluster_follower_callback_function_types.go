package impl

import (
	"expansion-gateway/dto/clusters"
	"expansion-gateway/interfaces/errorinfo"
)

type RequestAcceptClientCallback func(int64, int64, int32, int32, int32, int32, bool) (*clusters.ClustersSubscriptionRequestBody, errorinfo.GatewayError)
