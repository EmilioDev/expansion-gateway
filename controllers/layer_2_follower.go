// file: /controllers/layer_2_follower.go
package controllers

import (
	"expansion-gateway/clustering"
	"expansion-gateway/config"
	"expansion-gateway/dto/clusters/results"
	"expansion-gateway/dto/sessions"
	"expansion-gateway/enums"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
	"expansion-gateway/internal/structs"
)

type Layer2Follower struct {
	*Layer2Core
	clusterServer *clustering.ClusteringFollower
	subscriptions *structs.SessionsDictionary[*sessions.Layer2FollowerSubscription]
}

func (layer *Layer2Follower) GenerateUserSubscription(
	userID,
	requestedSessionID int64,
	clientType enums.ClientType,
	clientVersion byte,
	encryption enums.EncryptionAlgorithm,
	protocolVersion enums.ProtocolVersion,
	sessionResume bool) (*results.ClustersSubscriptionRequestBody, errorinfo.GatewayError) {
	subscription := sessions.GenerateLayer2FollowerSubscription(userID, requestedSessionID, clientType, clientVersion, encryption, protocolVersion, sessionResume)
	subscriptionId := layer.subscriptions.Add(subscription)

	return &results.ClustersSubscriptionRequestBody{
		SubscriptionID: subscriptionId,
		Challenge:      subscription.GetChallengeAsInt32Array(),
	}, nil
}

func (layer *Layer2Follower) initializeCluster() errorinfo.GatewayError {
	return layer.clusterServer.Start()
}

func (layer *Layer2Follower) stopCluster() errorinfo.GatewayError {
	return layer.clusterServer.Stop()
}

// ==== handler from layer 1 ====
func (layer *Layer2Follower) handlePacketFromLayer1(packet packets.Packet) errorinfo.GatewayError {
	return nil
}

// ==== handler from layer 3

func (layer *Layer2Follower) handlePacketFromLayer3(packet packets.Packet) errorinfo.GatewayError {
	return nil
}

// ==== constructor ====

func CreateNewLayer2Follower(conf *config.Configuration) *Layer2Follower {
	answer := &Layer2Follower{}
	core := CreateNewLayer2Core(
		conf,
		answer.handlePacketFromLayer1,
		answer.handlePacketFromLayer3,
		answer.initializeCluster,
		answer.stopCluster,
	)

	answer.Layer2Core = core
	answer.clusterServer = clustering.CreateClusteringFollower(core.wg, conf, answer)
	answer.subscriptions = structs.CreateNewSessionDictionary[*sessions.Layer2FollowerSubscription]()

	return answer
}
