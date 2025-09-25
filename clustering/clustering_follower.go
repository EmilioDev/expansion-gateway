package clustering

import (
	"expansion-gateway/config"
	"sync"
)

type ClusteringFollower struct {
	ClusterNode
}

func (cluster *ClusteringFollower) IsServer() bool {
	return true
}

func (cluster *ClusteringFollower) init() {
	//
}

func (cluster *ClusteringFollower) closeClient() {
	//
}

func CreateClusteringFollower(waiter *sync.WaitGroup, config *config.Configuration) *ClusteringFollower {
	result := ClusteringFollower{}

	result.ClusterNode = CreateBaseClusterNode(config, waiter, result.init, result.closeClient)

	return &result
}
