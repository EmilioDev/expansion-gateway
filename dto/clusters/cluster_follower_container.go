// file: /dto/clusters/cluster_follower_container.go
package clusters

import (
	"expansion-gateway/clustering/impl"
	"math"
	"sync"
	"time"
)

type ClusterFollowerContainer struct {
	lock                   sync.RWMutex                 // the mutex for the multiple updates
	healthy                bool                         // if the gateway is healthy
	messagesSinceLastCheck int64                        // messages received since last check
	activeSessions         int32                        // number of active sessions
	epoch                  int64                        // the current epoch of the gateway
	Client                 *impl.ClusterFollower_Client // the client used to interact remotelly with this cluster member
	lastUpdate             time.Time                    // the moment when the last update was done
}

func GetNewClusterFollowerContainer() *ClusterFollowerContainer {
	return &ClusterFollowerContainer{
		healthy:                true,
		messagesSinceLastCheck: 0,
		activeSessions:         0,
		epoch:                  0,
		Client:                 impl.CreateClusterFollowerClient(),
		lock:                   sync.RWMutex{},
		lastUpdate:             time.Now(),
	}
}

func (member *ClusterFollowerContainer) IsHealthy() bool {
	member.lock.RLock()
	defer member.lock.RUnlock()

	return member.healthy
}

func (member *ClusterFollowerContainer) SetHealthStatus(status bool) {
	member.lock.Lock()
	defer member.lock.Unlock()

	member.healthy = status
}

func (member *ClusterFollowerContainer) MessagesSinceLastCheck() int64 {
	member.lock.RLock()
	defer member.lock.RUnlock()

	return member.messagesSinceLastCheck
}

func (member *ClusterFollowerContainer) ActiveSessions() int32 {
	member.lock.RLock()
	defer member.lock.RUnlock()

	return member.activeSessions
}

func (member *ClusterFollowerContainer) Epoch() int64 {
	member.lock.RLock()
	defer member.lock.RUnlock()

	return member.epoch
}

func (member *ClusterFollowerContainer) UpdateStatus(messagesSinceLastCheck, epoch int64, activeSessions int32, healthy bool) {
	member.lock.Lock()
	defer member.lock.Unlock()

	member.messagesSinceLastCheck = messagesSinceLastCheck
	member.epoch = epoch
	member.activeSessions = activeSessions
	member.healthy = healthy
	member.lastUpdate = time.Now()
}

func (member *ClusterFollowerContainer) SecondsSinceLastUpdate() int64 {
	return int64(math.Round(time.Since(member.lastUpdate).Seconds()))
}
