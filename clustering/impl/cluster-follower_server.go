// file: /clustering/impl/cluster-follower_server.go
package impl

import (
	"context"
	"expansion-gateway/clustering/grpc"
)

type ClusterFollower_Server struct {
	grpc.UnimplementedExpansionGatewayClusterFollowerServer
}

func (server *ClusterFollower_Server) CheckIfFollowerIsOnline(context context.Context, empty *grpc.Empty) (*grpc.FollowerOperationResult, error) {
	return &grpc.FollowerOperationResult{
		Success: true,
	}, nil
}

func (server *ClusterFollower_Server) RequestAcceptClient(context context.Context, data *grpc.SubscriptionRequestData) (*grpc.RedirectionRequestResult, error) {
	return &grpc.RedirectionRequestResult{
		ErrorCode: 1,
		Body: &grpc.SubscriptionRequestResultBody{
			Challenge: []int32{0, 0, 0, 0, 0, 0, 0, 0},
		},
	}, nil
}
