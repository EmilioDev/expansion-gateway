// file: /controllers/layer_2_follower.go
package controllers

import (
	"crypto/ed25519"
	"expansion-gateway/clustering"
	"expansion-gateway/config"
	"expansion-gateway/dto"
	"expansion-gateway/dto/clusters/results"
	"expansion-gateway/dto/sessions"
	"expansion-gateway/enums"
	auth_errors "expansion-gateway/errors/auth"
	"expansion-gateway/errors/layererrors"
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

	if encryption != enums.NoEncryptionAlgorithm { // this branch is for handling payload encryption pass
		if err := subscription.GenerateEphemeralKeyPair(); err != nil {
			return nil, err
		}

		ephemeralPublicKey := subscription.EphemeralKeyPair.GetPublicKey()

		return &results.ClustersSubscriptionRequestBody{
			SubscriptionID:      subscriptionId,
			Challenge:           subscription.Challenge,
			GatewayPath:         layer.configuration.GetKcpPathToThisGateway(),
			SessionEphemeralKey: ephemeralPublicKey[:],
		}, nil
	}

	// enviar la llave publica efimera de la subscripcion como un []int32
	return &results.ClustersSubscriptionRequestBody{
		SubscriptionID:      subscriptionId,
		Challenge:           subscription.Challenge,
		GatewayPath:         layer.configuration.GetKcpPathToThisGateway(),
		SessionEphemeralKey: nil,
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
	const filePath string = "/controllers/layer_2_follower.go"

	switch packet.GetPacketType() {
	case enums.REDIRECTED:
		return layer.handleREDIRECTEDpacket(packet)
	}

	return layererrors.CreateProtocolFlowViolation_LayerError(
		filePath,
		60,
		enums.LAYER_2,
		enums.INCORRECT_PACKET_KIND,
	)
}

// ==== handler from layer 3

func (layer *Layer2Follower) handlePacketFromLayer3(packet packets.Packet) errorinfo.GatewayError {
	return nil
}

// ==== privates ====

// handler for redirected packets
func (layer *Layer2Follower) handleREDIRECTEDpacket(packet packets.Packet) errorinfo.GatewayError {
	const filePath string = "/controllers/layer_2_follower.go"

	if redirectedPacket, isRedirected := packet.(*dto.RedirectedPacket); isRedirected {
		if subscription, hasSubscription := layer.subscriptions.GetExists(redirectedPacket.SubscriptionID); hasSubscription {
			var key *[]byte = nil

			if subscription.ClientType == enums.GODOT_CLIENT {
				key = layer.configuration.GetGodotEd25519PublicKey()
			} else {
				key = layer.configuration.GetCliEd25519PublicKey()
			}

			if ed25519.Verify(*key, subscription.Challenge, redirectedPacket.Signature[:]) {
				newSession := sessions.GenerateNewLayer2Session(layer.configuration)
				newSession.UpdateFromFollowerSubscription(subscription)
				layer.subscriptions.Delete(redirectedPacket.SubscriptionID)

				sessionId := layer.sessions.Add(newSession)
				layer.approveSession(sessionId)

				return nil
			}
		}

		// this connection is unauthorized!!!
		layer.closeSession(redirectedPacket.SubscriptionID, enums.CloseReasonConnectionUnauthorized)

		return auth_errors.CreateConnectionUnauthorizedError(
			filePath,
			96,
			redirectedPacket.SubscriptionID,
		)
	}

	// this connection is not respecting the protocol
	return layererrors.CreateProtocolFlowViolation_LayerError(
		filePath,
		103,
		enums.LAYER_2,
		enums.INVALID_PACKET)
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
