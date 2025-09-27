// file: /clustering/cluster_follower.go
package clustering

import (
	"expansion-gateway/clustering/impl"
	"expansion-gateway/config"
	"expansion-gateway/dto/clusters/results"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/layers"
	"log"
	"sync"
	"time"
)

type ClusteringFollower struct {
	ClusterNode                                         // base
	leader                   *impl.ClusterLeader_Client // client used to interact with the cluster leader
	memberID                 int64                      // the id of this member in the cluster
	thisGateway              layers.Layer2              // the reference to the kcp gateway
	failedConsecutiveschecks int                        // the number of consecutives health checks that had failed
}

func (cluster *ClusteringFollower) IsLeader() bool {
	return false
}

// ==== init & close callbacks

// this method is called when the server starts, it configures the timer that periodically
// sends health checks to the cluster leader
func (cluster *ClusteringFollower) initCallback() {
	// we configure the wait group to wait for this goroutine first
	cluster.wg.Add(1)
	defer cluster.wg.Done()

	time.Sleep(time.Second * 10) // wait until this server starts

	// a bit defensive, yes, but mot bad
	if !cluster.isWorking.Load() {
		log.Fatalf("the follower with path of %s did not started", cluster.grpcCurrentServerPath)
	}

	// we attempt subscription
	if res, err := cluster.leader.Subscribe(cluster.grpcCurrentServerPath); err == nil {
		cluster.sessionsTimeout = res.HealthyTimeout
		cluster.memberID = res.ServerID

		// subcription ok, we start the update loop setting
		period := time.Duration(cluster.sessionsTimeout/2) * time.Second
		ticker := time.NewTicker(period)
		defer ticker.Stop()

		cluster.failedConsecutiveschecks = 0

		// and we have the tick loop here
		for range ticker.C {
			if !cluster.isWorking.Load() {
				return
			}

			// we send the health check to the leader
			if err := cluster.leader.HealthCheck(
				cluster.memberID,
				cluster.MessagesCounter.Load(),
				cluster.epoch.Load(),
				cluster.thisGateway.GetActiveSessions(),
				true,
			); err == nil {
				cluster.failedConsecutiveschecks = 0
			} else {
				cluster.failedConsecutiveschecks++

				// if we reach the consecutive 10 failed health checks, we collapse
				if cluster.failedConsecutiveschecks >= 10 {
					cluster.Stop()
					log.Fatalf("cluster follower at %s has reached the 10 consecutives failed health checks", cluster.grpcCurrentServerPath)
				}
			}
		}
	} else {
		// the subscription has failed, we notify and then exit
		log.Fatalf("the cluster leader rejected subscription of a follower with path of %s", cluster.grpcCurrentServerPath)
	}
}

func (cluster *ClusteringFollower) closeClientCallback() {
	cluster.thisGateway.Stop()
}

// ==== server callbacks ====

func (cluster *ClusteringFollower) acceptClientCallback(userID, requestedSessionId int64, clientType, clientVersion, encryption, protocolVersion int32, sessionResume bool) (*results.ClustersSubscriptionRequestBody, errorinfo.GatewayError) {
	return nil, nil
}

func (cluster *ClusteringFollower) hasThisSessionCallback(sessionId int64) (bool, errorinfo.GatewayError) {
	return cluster.thisGateway.HasSession(sessionId), nil
}

func (cluster *ClusteringFollower) requestExitCallback() errorinfo.GatewayError {
	return nil
}

func CreateClusteringFollower(waiter *sync.WaitGroup, config *config.Configuration, gateway layers.Layer2) *ClusteringFollower {
	if config.AreWeClusterLeaders() {
		log.Fatalln("a follower cannot be configured to be a leader of the cluster")
	}

	result := ClusteringFollower{}

	result.ClusterNode = CreateBaseClusterNode(config, waiter, result.initCallback, result.closeClientCallback)

	leader := impl.CreateClusterLeaderClient()
	leaderSrc := config.GetGrpcClusterLeaderPath()

	if err := leader.Connect(leaderSrc); err != nil {
		log.Fatalf("a follower has failed on conecting to the leader at %s", leaderSrc)
	}

	result.leader = leader
	result.thisGateway = gateway

	implementation := impl.CreateClusterFollowerServer(result.acceptClientCallback, result.hasThisSessionCallback, result.requestExitCallback)
	implementation.RegisterToGrpcServer(result.server)

	return &result
}
