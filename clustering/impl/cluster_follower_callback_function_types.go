// file: /clustering/impl/cluster_follower_callback_function_types.go
package impl

import (
	"expansion-gateway/dto/clusters"
	"expansion-gateway/interfaces/errorinfo"
)

type RequestAcceptClientCallback func(int64, int64, int32, int32, int32, int32, bool) (*clusters.ClustersSubscriptionRequestBody, errorinfo.GatewayError)
type GatewayHasThisSessionRegisteredCallback func(int64) (bool, errorinfo.GatewayError)
