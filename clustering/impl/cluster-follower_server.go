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

	acceptClientCallback   RequestAcceptClientCallback             // function to call when registering a client
	hasThisSessionCallback GatewayHasThisSessionRegisteredCallback // function to call if you want to check if a session is registered
	requestExitCallback    RequestExitCallback                     // function to be called when the cluster leader requests this member to go down
}

// constructor of the follower servers
func CreateClusterFollowerServer(acceptClientCallback RequestAcceptClientCallback,
	hasThisSessionCallback GatewayHasThisSessionRegisteredCallback,
	requestExitCallback RequestExitCallback) *ClusterFollower_Server {

	return &ClusterFollower_Server{
		grpc.UnimplementedExpansionGatewayClusterFollowerServer{},
		acceptClientCallback,
		hasThisSessionCallback,
		requestExitCallback,
	}
}

func (server *ClusterFollower_Server) DropClient(context context.Context, empty *grpc.Empty) (*grpc.FollowerOperationResult, error) {
	if err := server.requestExitCallback(); err != nil {
		return &grpc.FollowerOperationResult{
			Success: false,
		}, nil
	}

	return &grpc.FollowerOperationResult{
		Success: true,
	}, nil
}

// checks if this follower is online
func (server *ClusterFollower_Server) CheckIfFollowerIsOnline(context context.Context, empty *grpc.Empty) (*grpc.FollowerOperationResult, error) {
	return &grpc.FollowerOperationResult{
		Success: true,
	}, nil
}

// function called when the leader requests this gateway to handle a client
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

// function called when the leader wants to know if this follower is handling a session or not
func (server *ClusterFollower_Server) HasSession(context context.Context, data *grpc.FollowerSessionID) (*grpc.FollowerOperationResult, error) {
	if data != nil {
		if res, err := server.hasThisSessionCallback(data.SessionID); err == nil {
			return &grpc.FollowerOperationResult{
				Success: res,
			}, nil
		} else {
			return nil, status.Errorf(codes.Internal, "error: %s, code: %d", err.Error(), err.GetErrorCode())
		}
	}

	return nil, status.Error(codes.InvalidArgument, "you need to specify a valid parameter for checking if a session is registered in this gateway")
}
