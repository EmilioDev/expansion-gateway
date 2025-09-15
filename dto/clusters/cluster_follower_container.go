// file: /dto/clusters/cluster_follower_container.go
package clusters

import (
	"expansion-gateway/clustering/impl"
	"sync/atomic"
)

type ClusterFollowerContainer struct {
	healthy                atomic.Bool                  // if the gateway is healthy
	messagesSinceLastCheck atomic.Int64                 // messages received since last check
	activeSessions         atomic.Int32                 // number of active sessions
	epoch                  atomic.Int64                 // the current epoch of the gateway
	Client                 *impl.ClusterFollower_Client // the client used to interact remotelly with this cluster member
}

func GetNewClusterFollowerContainer() *ClusterFollowerContainer {
	return &ClusterFollowerContainer{
		healthy:                atomic.Bool{},
		messagesSinceLastCheck: atomic.Int64{},
		activeSessions:         atomic.Int32{},
		epoch:                  atomic.Int64{},
		Client:                 impl.CreateClusterFollowerClient(),
	}
}

func (member *ClusterFollowerContainer) IsHealthy() bool {
	return member.healthy.Load()
}

func (member *ClusterFollowerContainer) MessagesSinceLastCheck() int64 {
	return member.messagesSinceLastCheck.Load()
}

func (member *ClusterFollowerContainer) ActiveSessions() int32 {
	return member.activeSessions.Load()
}

func (member *ClusterFollowerContainer) Epoch() int64 {
	return member.epoch.Load()
}
