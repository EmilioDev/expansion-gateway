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
	globalWg    *sync.WaitGroup                                                 // global wait group, to inform when it is closed
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

	return &res.ClusterMemberSubscriptionResult{
		ServerID:       index,
		HealthyTimeout: cluster.sessionsTimeout,
	}, nil
}

func (cluster *ClusteringServer) unsubscribe(clientId int64) (bool, errorinfo.GatewayError) {
	cluster.clients.Delete(clientId)
	return true, nil
}

func (cluster *ClusteringServer) healthCheck(serverID, messagesSinceLastCheck, epoch int64, activeSessions int32, healthy bool) (bool, errorinfo.GatewayError) {
	return true, nil
}

func CreateClusteringServer(waiter *sync.WaitGroup, config *config.Configuration) *ClusteringServer {
	result := ClusteringServer{}

	result.ClusterNode = CreateBaseClusterNode(config)
	result.clients = structs.CreateNewSessionDictionary[*clusters.ClusterFollowerContainer]()
	result.globalWg = waiter

	implementation := impl.GenerateClusterLeaderServer(result.subscribe, result.unsubscribe, result.healthCheck)
	implementation.RegisterToGrpcServer(result.server)

	return &result
}
