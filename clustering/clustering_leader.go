// file: /clustering/cluster_server.go
package clustering

import (
	"expansion-gateway/clustering/impl"
	"expansion-gateway/config"
	"expansion-gateway/dto/clusters"
	res "expansion-gateway/dto/clusters/results"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/internal/structs"
	"sync"
	"time"
)

type ClusteringLeader struct {
	ClusterNode                                                                 // base
	clients     *structs.SessionsDictionary[*clusters.ClusterFollowerContainer] // clients
}

func (cluster *ClusteringLeader) IsServer() bool {
	return true
}

func (cluster *ClusteringLeader) subscribe(path string) (*res.ClusterMemberSubscriptionResult, errorinfo.GatewayError) {
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

func (cluster *ClusteringLeader) unsubscribe(clientId int64) (bool, errorinfo.GatewayError) {
	if client, exists := cluster.clients.GetExists(clientId); exists {
		client.Client.Disconnect()
		cluster.clients.Delete(clientId)
		cluster.NextEpoch()

		return true, nil
	}

	return false, nil
}

func (cluster *ClusteringLeader) healthCheck(serverID, messagesSinceLastCheck, epoch int64, activeSessions int32, healthy bool) (bool, errorinfo.GatewayError) {
	if member, exists := cluster.clients.GetExists(serverID); exists {
		member.UpdateStatus(messagesSinceLastCheck, epoch, activeSessions, healthy)
		return true, nil
	}

	return false, nil
}

func (cluster *ClusteringLeader) checkClients() {
	cluster.wg.Add(1)
	defer cluster.wg.Done()

	timeout := cluster.sessionsTimeout

	period := time.Duration(timeout) * time.Second
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	var keys []int64

	for range ticker.C {
		if !cluster.isWorking.Load() {
			return
		}

		keys = cluster.clients.Keys()

		for _, key := range keys {
			if client, exists := cluster.clients.GetExists(key); exists {
				if client.IsHealthy() && client.SecondsSinceLastUpdate() > timeout {
					client.SetHealthStatus(false)
				}
			}
		}
	}
}

// closes all the clients and then removes them all from the list
func (cluster *ClusteringLeader) close() {
	keys := cluster.clients.Keys()

	for _, key := range keys {
		if client, exists := cluster.clients.GetExists(key); exists {
			client.Client.DropClient()
			client.Client.Disconnect()
		}
	}

	cluster.clients.Clear()
}

// creates a new cluster server
func CreateClusteringLeader(waiter *sync.WaitGroup, config *config.Configuration) *ClusteringLeader {
	result := ClusteringLeader{}

	result.ClusterNode = CreateBaseClusterNode(config, waiter, result.checkClients, result.close)
	result.clients = structs.CreateNewSessionDictionary[*clusters.ClusterFollowerContainer]()

	implementation := impl.GenerateClusterLeaderServer(result.subscribe, result.unsubscribe, result.healthCheck)
	implementation.RegisterToGrpcServer(result.server)

	return &result
}
