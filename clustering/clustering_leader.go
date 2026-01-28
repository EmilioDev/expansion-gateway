// file: /clustering/cluster_leader.go
package clustering

import (
	"expansion-gateway/clustering/impl"
	"expansion-gateway/config"
	"expansion-gateway/dto/clusters"
	res "expansion-gateway/dto/clusters/results"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/layers"
	"expansion-gateway/internal/structs"
	"fmt"
	"sync"
	"time"
)

type ClusteringLeader struct {
	ClusterNode                                                                 // base
	Clients     *structs.SessionsDictionary[*clusters.ClusterFollowerContainer] // clients
	thisGateway layers.Layer2Leader                                             // reference to the inner gateway
}

func (cluster *ClusteringLeader) IsLeader() bool {
	return true
}

func (cluster *ClusteringLeader) subscribe(path string) (*res.ClusterMemberSubscriptionResult, errorinfo.GatewayError) {
	candidate := clusters.GetNewClusterFollowerContainer()

	if err := candidate.Client.Connect(path); err != nil {
		return nil, err
	}

	index := cluster.Clients.Add(candidate)
	cluster.NextEpoch()

	return &res.ClusterMemberSubscriptionResult{
		ServerID:       index,
		HealthyTimeout: cluster.sessionsTimeout,
	}, nil
}

func (cluster *ClusteringLeader) unsubscribe(clientId int64) (bool, errorinfo.GatewayError) {
	if client, exists := cluster.Clients.GetExists(clientId); exists {
		client.Client.Disconnect()
		cluster.Clients.Delete(clientId)
		cluster.NextEpoch()

		return true, nil
	}

	return false, nil
}

func (cluster *ClusteringLeader) healthCheck(
	serverID,
	messagesSinceLastCheck,
	epoch int64,
	activeSessions int32,
	cpuUsage,
	ramUsage float32,
	healthy bool,
) (bool, errorinfo.GatewayError) {
	if member, exists := cluster.Clients.GetExists(serverID); exists {
		member.UpdateStatus(messagesSinceLastCheck, epoch, activeSessions, ramUsage, cpuUsage, healthy)
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
	fmt.Println("cluster leader is running")

	for range ticker.C {
		if !cluster.isWorking.Load() {
			return
		}

		keys = cluster.Clients.Keys()

		for _, key := range keys {
			if client, exists := cluster.Clients.GetExists(key); exists {
				if client.IsHealthy() && client.SecondsSinceLastUpdate() > timeout {
					client.SetHealthStatus(false)
				}
			}
		}
	}
}

// closes all the clients and then removes them all from the list
func (cluster *ClusteringLeader) close() {
	keys := cluster.Clients.Keys()

	for _, key := range keys {
		if client, exists := cluster.Clients.GetExists(key); exists {
			client.Client.DropClient()
			client.Client.Disconnect()
		}
	}

	cluster.Clients.Clear()
}

func (cluster *ClusteringLeader) markUserAsRedirected(userId int64) {
	cluster.thisGateway.MarkUserAsRedirected(userId)
}

// creates a new cluster server
func CreateClusteringLeader(
	waiter *sync.WaitGroup,
	config *config.Configuration,
	gateway layers.Layer2Leader) *ClusteringLeader {
	result := ClusteringLeader{}

	result.thisGateway = gateway

	result.ClusterNode = CreateBaseClusterNode(config, waiter, result.checkClients, result.close)
	result.Clients = structs.CreateNewSessionDictionary[*clusters.ClusterFollowerContainer]()

	implementation := impl.GenerateClusterLeaderServer(
		result.subscribe,
		result.unsubscribe,
		result.healthCheck,
		result.markUserAsRedirected)
	implementation.RegisterToGrpcServer(result.server)

	return &result
}
