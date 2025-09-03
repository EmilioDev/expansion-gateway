package impl

import (
	"context"
	"expansion-gateway/clustering/grpc"
	"expansion-gateway/dto"
	"expansion-gateway/enums"
	"expansion-gateway/errors/clustererrors"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/errorinfo"
	"flag"
	"math"
	"time"

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
	address := flag.String("addr", source, "the address to connect to")

	if conn, err := google.NewClient(*address, google.WithTransportCredentials(insecure.NewCredentials())); err == nil {
		c.connection = conn
		c.client = grpc.NewExpansionGatewayClusterFollowerClient(conn)

		return nil
	}

	return helpers.WithStackTrace(clustererrors.CreateConnectionToServerFailedError(
		"/clustering/impl/cluster-leader_client.go",
		30,
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

func (c *ClusterFollower_Client) isReady() errorinfo.GatewayError {
	const filePath string = "/clustering/impl/cluster-leader_client.go"

	if c.client == nil {
		return clustererrors.CreateClientNotReadyError(filePath, 55, enums.ClusterLeader)
	}

	if c.connection == nil {
		return clustererrors.CreateClientNotReadyError(filePath, 59, enums.ClusterLeader)
	}

	return nil
}

func (c *ClusterFollower_Client) CheckIfFollowerIsOnline() (bool, errorinfo.GatewayError) {
	if err := c.isReady(); err != nil {
		return false, err
	}

	const filePath string = "/clustering/impl/cluster-follower_client.go"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if res, err := c.client.CheckIfFollowerIsOnline(ctx, &grpc.Empty{}); err == nil {
		return res.Success, nil
	}

	return false, clustererrors.CreateOperationFailedError(filePath,
		87,
		enums.ClusterFollower,
		false)
}

func (c *ClusterFollower_Client) RequestAcceptClient(userID int64, userFrame *dto.SessionFrame) (*dto.ClusterUserSubscriptionResult, errorinfo.GatewayError) {
	if err := c.isReady(); err != nil {
		return nil, err
	}

	const filePath string = "/clustering/impl/cluster-follower_client.go"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if res, err := c.client.RequestAcceptClient(ctx, userFrame.ToSubscriptionRequestData(userID)); err == nil {
		if res.Body == nil {
			return nil, clustererrors.CreateNoPayloadError(filePath,
				104,
				enums.ClusterFollower,
				false)
		}

		response := dto.ClusterUserSubscriptionResult{
			Challenge:      []byte{},
			SubscriptionID: res.Body.SubscriptionID,
		}

		iterations := len(res.Body.Challenge)

		for x := 0; x < iterations; x++ {
			toConvert := res.Body.Challenge[x]

			if toConvert < 0 {
				toConvert = 0
			} else if toConvert > math.MaxUint8 {
				toConvert = 255
			}

			response.Challenge = append(response.Challenge, byte(toConvert))
		}

		return &response, nil
	}

	return nil, clustererrors.CreateOperationFailedError(filePath,
		133,
		enums.ClusterFollower,
		false)
}

func (c *ClusterFollower_Client) HasThisSession(sessionID int64) (bool, errorinfo.GatewayError) {
	if err := c.isReady(); err != nil {
		return false, err
	}

	const filePath string = "/clustering/impl/cluster-follower_client.go"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if res, err := c.client.HasSession(ctx, &grpc.FollowerSessionID{SessionID: sessionID}); err == nil {
		return res.Success, nil
	}

	return false, clustererrors.CreateOperationFailedError(filePath,
		153,
		enums.ClusterFollower,
		false)
}
