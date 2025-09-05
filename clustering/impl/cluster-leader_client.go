// file: /clustering/impl/cluster-leader_client.go
package impl

import (
	"context"
	dto "expansion-gateway/dto/clusters"
	"expansion-gateway/enums"
	"expansion-gateway/errors/clustererrors"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/errorinfo"
	"flag"
	"time"

	"expansion-gateway/clustering/grpc"

	google "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClusterLeader_Client struct {
	connection *google.ClientConn
	client     grpc.ExpansionGatewayClusterLeaderClient
}

func CreateClusterLeaderClient() *ClusterLeader_Client {
	return &ClusterLeader_Client{
		connection: nil,
		client:     nil,
	}
}

func (client *ClusterLeader_Client) Connect(source string) errorinfo.GatewayError {
	address := flag.String("addr", source, "the address to connect to")

	if conn, err := google.NewClient(*address, google.WithTransportCredentials(insecure.NewCredentials())); err == nil {
		client.connection = conn
		client.client = grpc.NewExpansionGatewayClusterLeaderClient(conn)

		return nil
	}

	return helpers.WithStackTrace(clustererrors.CreateConnectionToServerFailedError(
		"/clustering/impl/cluster-leader_client.go",
		34,
		source,
		enums.ClusterLeader,
		false), 0)
}

func (client *ClusterLeader_Client) Disconnect() errorinfo.GatewayError {
	if client.client != nil {
		client.client = nil
	}

	if client.connection != nil {
		client.connection.Close()
		client.connection = nil
	}

	return nil
}

func (client *ClusterLeader_Client) isReady() errorinfo.GatewayError {
	const filePath string = "/clustering/impl/cluster-leader_client.go"

	if client.client == nil {
		return clustererrors.CreateClientNotReadyError(filePath, 55, enums.ClusterLeader)
	}

	if client.connection == nil {
		return clustererrors.CreateClientNotReadyError(filePath, 59, enums.ClusterLeader)
	}

	return nil
}

func (client *ClusterLeader_Client) Subscribe(grpcCurrentNodePath string) (*dto.ClusterMemberSubscriptionResult, errorinfo.GatewayError) {
	if err := client.isReady(); err != nil {
		return nil, err
	}

	const filePath string = "/clustering/impl/cluster-leader_client.go"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if res, err := client.client.Subscribe(ctx, &grpc.FollowerSubscriptionData{
		GrpcServicePath: grpcCurrentNodePath,
	}); err == nil {
		if res.SubscriptionBody != nil {
			return &dto.ClusterMemberSubscriptionResult{
				ServerID:       res.SubscriptionBody.ServerID,
				HealthyTimeout: res.SubscriptionBody.HealthyTimeout,
			}, nil
		} else {
			return nil, clustererrors.CreateNoPayloadError(filePath, 87, enums.ClusterLeader, false)
		}
	}

	return nil, clustererrors.CreateOperationFailedError(filePath, 91, enums.ClusterLeader, false)
}

func (client *ClusterLeader_Client) Unsubscribe(serverId int64) errorinfo.GatewayError {
	if err := client.isReady(); err != nil {
		return err
	}

	const filePath string = "/clustering/impl/cluster-leader_client.go"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if res, err := client.client.Unsubscribe(ctx, &grpc.FollowerUnsubscriptionData{ServerID: serverId}); err == nil {
		if res.Success {
			return nil
		}
	}

	return clustererrors.CreateOperationFailedError(filePath, 109, enums.ClusterLeader, false)
}

func (client *ClusterLeader_Client) HealthCheck(serverId,
	messagesSinceLastCheck,
	epoch int64,
	activeSessions int32,
	healthy bool) errorinfo.GatewayError {
	if err := client.isReady(); err != nil {
		return err
	}

	const filePath string = "/clustering/impl/cluster-leader_client.go"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if res, err := client.client.HealthCheck(ctx,
		&grpc.HealthCheckData{
			ServerId:               serverId,
			ActiveSessions:         activeSessions,
			MessagesSinceLastCheck: messagesSinceLastCheck,
			Epoch:                  epoch,
			Healthy:                healthy,
		}); err == nil {
		if res.Success {
			return nil
		}
	}

	return clustererrors.CreateOperationFailedError(filePath, 138, enums.ClusterLeader, false)
}
