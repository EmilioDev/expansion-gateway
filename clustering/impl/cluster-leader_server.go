package impl

import (
	"context"
	"expansion-gateway/clustering/grpc"
)

type ClusterLeader_Server struct {
	grpc.UnimplementedExpansionGatewayClusterLeaderServer
}

func (server *ClusterLeader_Server) Subscribe(context context.Context,
	data *grpc.FollowerSubscriptionData) (*grpc.SubscriptionResult, error) {
	return &grpc.SubscriptionResult{
		ErrorCode:        1,
		SubscriptionBody: nil,
	}, nil
}

func (server *ClusterLeader_Server) Unsubscribe(context context.Context, data *grpc.FollowerUnsubscriptionData) (*grpc.ServerOperationResult, error) {
	return &grpc.ServerOperationResult{
		Success: false,
	}, nil
}

func (server *ClusterLeader_Server) HealthCheck(context context.Context, data *grpc.HealthCheckData) (*grpc.ServerOperationResult, error) {
	return &grpc.ServerOperationResult{
		Success: false,
	}, nil
}
