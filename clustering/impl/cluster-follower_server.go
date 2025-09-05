// file: /clustering/impl/cluster-follower_server.go
package impl

import (
	"context"
	"expansion-gateway/clustering/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ClusterFollower_Server struct {
	grpc.UnimplementedExpansionGatewayClusterFollowerServer

	// ==== callbacks ====

	acceptClientCallback RequestAcceptClientCallback
}

// checks if this follower is online
func (server *ClusterFollower_Server) CheckIfFollowerIsOnline(context context.Context, empty *grpc.Empty) (*grpc.FollowerOperationResult, error) {
	return &grpc.FollowerOperationResult{
		Success: true,
	}, nil
}

func (server *ClusterFollower_Server) RequestAcceptClient(context context.Context, data *grpc.SubscriptionRequestData) (*grpc.RedirectionRequestResult, error) {
	if data != nil {
		if res, err := server.acceptClientCallback(
			data.UserID,
			data.ReqestedSessionID,
			data.ClientType,
			data.ClientVersion,
			data.Encryption,
			data.ProtocolVersion,
			data.SessionResume,
		); err == nil {
			return &grpc.RedirectionRequestResult{
				Success: true,
				Body: &grpc.SubscriptionRequestResultBody{
					SubscriptionID: res.SubscriptionID,
					Challenge:      res.Challenge,
				},
			}, nil
		} else {
			return nil, status.Errorf(codes.Internal, "error: %s, code: %d", err.Error(), err.GetErrorCode())
		}
	}

	return nil, status.Error(codes.InvalidArgument, "you need to specify a valid parameter for requesting to accept an user")
}
