// file: /clustering/impl/cluster-follower_client.go
package impl

import (
	"context"
	"expansion-gateway/clustering/grpc"
	clusters "expansion-gateway/dto/clusters/results"
	dto "expansion-gateway/dto/sessions"
	"expansion-gateway/enums"
	"expansion-gateway/errors/clustererrors"
	"expansion-gateway/helpers"
	"expansion-gateway/helpers/constants"
	"expansion-gateway/interfaces/errorinfo"
	"flag"
	"fmt"

	google "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClusterFollower_Client struct {
	connection *google.ClientConn
	client     grpc.ExpansionGatewayClusterFollowerClient
}

func CreateClusterFollowerClient() *ClusterFollower_Client {
	return &ClusterFollower_Client{
		connection: nil,
		client:     nil,
	}
}

func (c *ClusterFollower_Client) Connect(source string) errorinfo.GatewayError {
	addr := fmt.Sprintf("addr%d", helpers.GenerateRandomInt64())
	address := flag.String(addr, source, "the address to connect to")

	if conn, err := google.NewClient(*address, google.WithTransportCredentials(insecure.NewCredentials())); err == nil {
		c.connection = conn
		c.client = grpc.NewExpansionGatewayClusterFollowerClient(conn)

		return nil
	}

	return helpers.WithStackTrace(clustererrors.CreateConnectionToServerFailedError(
		"/clustering/impl/cluster-leader_client.go",
		43,
		source,
		enums.ClusterLeader,
		false), 0)
}

func (c *ClusterFollower_Client) Disconnect() errorinfo.GatewayError {
	if c.client != nil {
		c.client = nil
	}

	if c.connection != nil {
		c.connection.Close()
		c.connection = nil
	}

	return nil
}

func (c *ClusterFollower_Client) DropClient() (bool, errorinfo.GatewayError) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.CLUSTER_REQUEST_TIMEOUT)
	defer cancel()
	const filePath string = "/clustering/impl/cluster-leader_client.go"

	if res, err := c.client.DropClient(ctx, &grpc.Empty{}); err != nil {
		return false, clustererrors.CreateOperationFailedError(filePath,
			69,
			enums.ClusterFollower,
			false)
	} else {
		return res.Success, nil
	}
}

func (c *ClusterFollower_Client) isReady() errorinfo.GatewayError {
	const filePath string = "/clustering/impl/cluster-leader_client.go"

	if c.client == nil {
		return clustererrors.CreateClientNotReadyError(filePath, 83, enums.ClusterLeader)
	}

	if c.connection == nil {
		return clustererrors.CreateClientNotReadyError(filePath, 87, enums.ClusterLeader)
	}

	return nil
}

func (c *ClusterFollower_Client) CheckIfFollowerIsOnline() (bool, errorinfo.GatewayError) {
	if err := c.isReady(); err != nil {
		return false, err
	}

	const filePath string = "/clustering/impl/cluster-follower_client.go"
	ctx, cancel := context.WithTimeout(context.Background(), constants.CLUSTER_REQUEST_TIMEOUT)
	defer cancel()

	if res, err := c.client.CheckIfFollowerIsOnline(ctx, &grpc.Empty{}); err == nil {
		return res.Success, nil
	}

	return false, clustererrors.CreateOperationFailedError(filePath,
		106,
		enums.ClusterFollower,
		false)
}

func (c *ClusterFollower_Client) RequestAcceptClient(userID int64, userFrame *dto.SessionFrame) (*clusters.ClusterUserSubscriptionResult, errorinfo.GatewayError) {
	if err := c.isReady(); err != nil {
		return nil, err
	}

	const filePath string = "/clustering/impl/cluster-follower_client.go"
	ctx, cancel := context.WithTimeout(context.Background(), constants.CLUSTER_REQUEST_TIMEOUT)
	defer cancel()

	if res, err := c.client.RequestAcceptClient(ctx, userFrame.ToSubscriptionRequestData(userID)); err == nil {
		if res.Body == nil {
			return nil, clustererrors.CreateNoPayloadError(filePath,
				122,
				enums.ClusterFollower,
				false)
		}

		response := clusters.ClusterUserSubscriptionResult{
			Challenge:           helpers.ConvertInt32ArrayIntoByteArray(res.Body.Challenge),
			SubscriptionID:      res.Body.SubscriptionID,
			GatewayPath:         res.Body.NewGatewayAddress,
			SessionEphemeralKey: helpers.ConvertInt32ArrayIntoByteArray(res.Body.ServerPublicEphemeralKey),
		}

		return &response, nil
	}

	return nil, clustererrors.CreateOperationFailedError(filePath,
		151,
		enums.ClusterFollower,
		false)
}

func (c *ClusterFollower_Client) HasThisSession(sessionID int64) (bool, errorinfo.GatewayError) {
	if err := c.isReady(); err != nil {
		return false, err
	}

	const filePath string = "/clustering/impl/cluster-follower_client.go"
	ctx, cancel := context.WithTimeout(context.Background(), constants.CLUSTER_REQUEST_TIMEOUT)
	defer cancel()

	if res, err := c.client.HasSession(ctx, &grpc.FollowerSessionID{SessionID: sessionID}); err == nil {
		return res.Success, nil
	}

	return false, clustererrors.CreateOperationFailedError(filePath,
		170,
		enums.ClusterFollower,
		false)
}
