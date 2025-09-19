package clustering

import (
	"expansion-gateway/clustering/impl"
	"expansion-gateway/config"
	"expansion-gateway/dto/clusters"
	res "expansion-gateway/dto/clusters/results"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/internal/structs"
	"sync"
)

type ClusteringServer struct {
	ClusterNode                                                                 // base
	clients     *structs.SessionsDictionary[*clusters.ClusterFollowerContainer] // clients
}

func (cluster *ClusteringServer) IsServer() bool {
	return true
}

func (cluster *ClusteringServer) subscribe(path string) (*res.ClusterMemberSubscriptionResult, errorinfo.GatewayError) {
	candidate := clusters.GetNewClusterFollowerContainer()

	if err := candidate.Client.Connect(path); err != nil {
		return nil, err
	}

	index := cluster.clients.Add(candidate)
	cluster.NextEpoch()

	return &res.ClusterMemberSubscriptionResult{
		ServerID:       index,
		HealthyTimeout: cluster.sessionsTimeout,
	}, nil
}

func (cluster *ClusteringServer) unsubscribe(clientId int64) (bool, errorinfo.GatewayError) {
	cluster.clients.Delete(clientId)
	cluster.NextEpoch()

	return true, nil
}

func (cluster *ClusteringServer) healthCheck(serverID, messagesSinceLastCheck, epoch int64, activeSessions int32, healthy bool) (bool, errorinfo.GatewayError) {
	if member, exists := cluster.clients.GetExists(serverID); exists {
		member.UpdateStatus(messagesSinceLastCheck, epoch, activeSessions, healthy)
		return true, nil
	}

	return false, nil
}

func CreateClusteringServer(waiter *sync.WaitGroup, config *config.Configuration) *ClusteringServer {
	result := ClusteringServer{}

	result.ClusterNode = CreateBaseClusterNode(config, waiter)
	result.clients = structs.CreateNewSessionDictionary[*clusters.ClusterFollowerContainer]()

	implementation := impl.GenerateClusterLeaderServer(result.subscribe, result.unsubscribe, result.healthCheck)
	implementation.RegisterToGrpcServer(result.server)

	return &result
}
