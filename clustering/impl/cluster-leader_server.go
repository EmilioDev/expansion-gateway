// file: /clustering/impl/cluster-leader_server.go
package impl

import (
	"context"
	"expansion-gateway/clustering/grpc"

	rpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ClusterLeader_Server struct {
	grpc.UnimplementedExpansionGatewayClusterLeaderServer // grpc layer

	// ==== callbacks ====

	subscribeCallback      ClusterLeaderSubscribeCallback      // subscribe callback
	unsubscribeCallback    ClusterLeaderUnsubscribeCallback    // unsubscribe callback
	healthCheckCallback    ClusterLeaderHealthCheckCallback    // health check callback
	userRedirectedCallback ClusterLeaderUserRedirectedCallback // user redirected callback
}

// creates a new server for a cluster leader
func GenerateClusterLeaderServer(
	subscribeCallback ClusterLeaderSubscribeCallback,
	unsubscribeCallback ClusterLeaderUnsubscribeCallback,
	healthCheckCallback ClusterLeaderHealthCheckCallback,
	userRedirectedCallback ClusterLeaderUserRedirectedCallback) *ClusterLeader_Server {
	return &ClusterLeader_Server{
		grpc.UnimplementedExpansionGatewayClusterLeaderServer{},
		subscribeCallback,
		unsubscribeCallback,
		healthCheckCallback,
		userRedirectedCallback,
	}
}

// subscribes to this server
func (sv *ClusterLeader_Server) Subscribe(context context.Context,
	data *grpc.FollowerSubscriptionData) (*grpc.SubscriptionResult, error) {
	if data != nil {
		if res, err := sv.subscribeCallback(data.GrpcServicePath); err == nil {
			return &grpc.SubscriptionResult{
				ErrorCode: 0,
				SubscriptionBody: &grpc.SubscriptionResultBody{
					ServerID:       res.ServerID,
					HealthyTimeout: res.HealthyTimeout,
				},
			}, nil
		} else {
			return nil, status.Errorf(codes.Internal, "error: %s, code: %d", err.Error(), err.GetErrorCode())
		}
	}

	return nil, status.Error(codes.InvalidArgument, "you need to specify a valid parameter for subscription")
}

// removes a subscription to this server
func (sv *ClusterLeader_Server) Unsubscribe(context context.Context, data *grpc.FollowerUnsubscriptionData) (*grpc.ServerOperationResult, error) {
	if data != nil {
		if res, err := sv.unsubscribeCallback(data.ServerID); err == nil {
			return &grpc.ServerOperationResult{
				Success: res,
			}, nil
		} else {
			return nil, status.Errorf(codes.Internal, "error: %s, code: %d", err.Error(), err.GetErrorCode())
		}
	}

	return nil, status.Error(codes.InvalidArgument, "you need to specify a valid parameter for the unsubscription")
}

// sends a health check to this server
func (sv *ClusterLeader_Server) HealthCheck(context context.Context, data *grpc.HealthCheckData) (*grpc.ServerOperationResult, error) {
	if data != nil {
		if res, err := sv.healthCheckCallback(
			data.ServerId,
			data.MessagesSinceLastCheck,
			data.Epoch,
			data.ActiveSessions,
			data.CpuPercentUsage,
			data.RamPercentUsage,
			data.Healthy,
		); err == nil {
			return &grpc.ServerOperationResult{
				Success: res,
			}, nil
		} else {
			return nil, status.Errorf(codes.Internal, "error: %s, code: %d", err.Error(), err.GetErrorCode())
		}
	}

	return nil, status.Error(codes.InvalidArgument, "you need to specify a valid parameter for the health check")
}

func (sv *ClusterLeader_Server) RegisterToGrpcServer(server *rpc.Server) {
	grpc.RegisterExpansionGatewayClusterLeaderServer(server, sv)
}

func (sv *ClusterLeader_Server) UserIsRedirected(ctx context.Context, data *grpc.UserRedirectedData) (*grpc.ServerOperationResult, error) {
	sv.userRedirectedCallback(data.UserID)

	return &grpc.ServerOperationResult{
		Success: true,
	}, nil
}
