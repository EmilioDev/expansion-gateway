// file: /controllers/layer_2_leader.go
package controllers

import (
	"crypto/ed25519"
	"expansion-gateway/clustering"
	"expansion-gateway/config"
	"expansion-gateway/dto"
	"expansion-gateway/dto/clusters"
	dtoSessions "expansion-gateway/dto/sessions"
	"expansion-gateway/enums"
	"expansion-gateway/errors/auth"
	"expansion-gateway/errors/layererrors"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
	"os"
)

type Layer2Leader struct {
	*Layer2Core
	clusterServer *clustering.ClusteringLeader
	pid           int32
}

func (layer *Layer2Leader) MarkUserAsRedirected(userId int64) {
	if session, sessionExists := layer.sessions.GetExists(userId); sessionExists {
		if session.GetState() == enums.REDIRECTING {
			// we mark this session as already redirected
			session.SetState(enums.REDIRECTED_TO_MEMBER)

			// we refresh the session and extend the lifetime
			session.RefreshActivity()
			session.ExtendTimeoutInterval()

			// and we disable the inputs
			layer.closeSessionInLayer3(userId, enums.CloseReasonUserRedirected)
			layer.layer1.DisableSession(userId)
		}
	}
}

func (layer *Layer2Leader) initializeCluster() errorinfo.GatewayError {
	return layer.clusterServer.Start()
}

func (layer *Layer2Leader) stopCluster() errorinfo.GatewayError {
	return layer.clusterServer.Stop()
}

// ==== layer 1 packet handlers ====

// global packet handler from layer 1
func (layer *Layer2Leader) handlePacketFromLayer1(packet packets.Packet) errorinfo.GatewayError {
	const filePath string = "/controllers/layer_2_leader.go"

	layer.clusterServer.NewMessage()

	switch packet.GetPacketType() {
	case enums.HELLO:
		if helloPacket, ok := packet.(*dto.HelloPacket); ok {
			return layer.handleHelloPacket(helloPacket)
		}

		return layererrors.CreateProtocolFlowViolation_LayerError(
			filePath,
			69,
			enums.LAYER_2,
			enums.INCORRECT_PACKET_KIND)

	case enums.CHALLENGE: // clients should never send a challenge
		return layererrors.CreateProtocolFlowViolation_LayerError(
			filePath,
			76,
			enums.LAYER_2,
			enums.CLIENT_SENT_CHALLENGE)

	case enums.NONE: //the layer 1 received an invalid packet
		return layererrors.CreateProtocolFlowViolation_LayerError(
			filePath,
			83,
			enums.LAYER_2,
			enums.INVALID_PACKET)

	case enums.CONNECT:
		if connectPacket, ok := packet.(*dto.ConnectPacket); ok {
			return layer.handleConnectPacket(connectPacket)
		}

		return layererrors.CreateProtocolFlowViolation_LayerError(
			filePath,
			94,
			enums.LAYER_2,
			enums.INCORRECT_PACKET_KIND)
	}

	return nil
}

// layer 1 connect packet handler
func (layer *Layer2Leader) handleConnectPacket(packet *dto.ConnectPacket) errorinfo.GatewayError {
	sessionId := packet.GetSender()
	const filePath string = "/controllers/layer_2_leader.go"

	if layer.isRedirecting(sessionId) {
		return nil
	}

	if session, sessionExists := layer.sessions.GetExists(sessionId); sessionExists {
		session.RefreshActivity()

		if state := session.GetState(); state == enums.CHALLENGE_SENT || state == enums.RECEIVED_CONNECT {
			session.SetState(enums.RECEIVED_CONNECT)
			publicKey := session.GetEd25519PublicKey()
			challenge := session.GetChallenge()

			if session.Encryption.GetEncryptionAlgorithm() == enums.NoEncryptionAlgorithm {
				if ed25519.Verify(publicKey, challenge, packet.Signature[:]) {
					layer.authorizeSession(sessionId) // handled here
				} else {
					// this client is unauthorized!!!
					layer.closeSession(sessionId, enums.CloseReasonFailedAuthentication) // handled
				}
			} else {
				if packet.ClientEphemeralKey == nil {
					return auth.GenerateConnectWithEphemeralKeyMissing(filePath, 125, sessionId)
				}

				clientEphemeralKey := *packet.ClientEphemeralKey

				msg := make([]byte, 0, len(clientEphemeralKey)+len(challenge))
				msg = append(msg, challenge...)
				msg = append(msg, clientEphemeralKey[:]...)

				if ed25519.Verify(publicKey, msg, packet.Signature[:]) {
					layer.authorizeSession(sessionId)                     // session authorized
					session.Encryption.GenerateNewKey(clientEphemeralKey) // final password generated
					session.Encryption.DeleteEphemeralKeys()              // ephemeral keys deleted
				} else {
					// this client is unauthorized!!!
					layer.closeSession(sessionId, enums.CloseReasonFailedAuthentication) // handled
				}
			}

			return nil
		}

		// i think, instead of doing this, we should give the connect an extra functionality:
		// if the user wants its datas again, it could just send a connect
		return layererrors.CreateProtocolFlowViolation_LayerError(
			filePath,
			133,
			enums.LAYER_2,
			enums.CLIENT_SENT_CONNECT_AT_WRONG_MOMENT)

	}

	return layererrors.CreateProtocolFlowViolation_LayerError(
		filePath,
		141,
		enums.LAYER_2,
		enums.SESSION_CLOSED)
}

// layer 1 hello packet handler
func (layer *Layer2Leader) handleHelloPacket(packet *dto.HelloPacket) errorinfo.GatewayError {
	clientId := packet.GetSender()
	const filePath string = "/controllers/layer_2_leader.go"
	var newChallenge []byte
	var err errorinfo.GatewayError = nil

	if layer.isRedirecting(clientId) {
		return nil
	}

	// if the session exists, then this is a retry packet
	if sessionStored, sessionExist := layer.sessions.GetExists(clientId); sessionExist {
		connectionState := sessionStored.GetState()
		sessionStored.RefreshActivity()

		// we check if the current state of the connection allows retrying hello
		if connectionState == enums.CHALLENGE_SENT || connectionState == enums.HELLO_RECEIVED {
			// we update the session from the hello packet (this is a retry hello)
			sessionStored.UpdateFromHelloPacket(packet)

			// we generate a new challenge nonce, and use it to generate a challenge packet,
			// send the packet to the client, and store the nonce for later check
			if newChallenge, err = helpers.GenerateChallengeNonce(); err == nil {
				sessionStored.SetChallenge(&newChallenge)
			} else {
				newChallenge = helpers.GetDefaultChallengeNonce()
				sessionStored.SetChallenge(&newChallenge)
			}

			//newChallengePacket := dto.GenerateChallengePacket(clientId, &newChallenge)
			var newChallengePacket *dto.ChallengePacket = nil

			if sessionStored.GetEncryption() == enums.NoEncryptionAlgorithm {
				// we generate the challenge with the challenge only
				newChallengePacket = dto.GenerateChallengePacket(clientId, &newChallenge)
			} else {
				// we generate the ephemeral keys and check if everything is ok so far
				if err := sessionStored.Encryption.GenerateEphemeralKeys(); err != nil {
					return err
				}

				// we generate the challenge packet with the challenge and the server ephemeral public key
				newChallengePacket = dto.GenerateChallengePacketWithServerPublicEphemeralKey(
					clientId,
					&newChallenge,
					sessionStored.Encryption.GetEphemeralKeys().GetPublicKey(),
				)
			}

			if err = layer.layer1.SendPacket(newChallengePacket); err == nil {
				if connectionState == enums.HELLO_RECEIVED {
					sessionStored.SetState(enums.CHALLENGE_SENT)
				}
			}

			return err
		}

		// if the state of the session is not CHALLENGE_SENT or HELLO_RECEIVED, then this is an invalid packet
		// apply the corresponding measure for sending invalid packet
		return layererrors.CreateProtocolFlowViolation_LayerError(filePath, 205, enums.LAYER_2, enums.INVALID_HELLO)
	}

	// the session does not exist
	// then we generate a new one
	newSession := dtoSessions.GenerateNewLayer2Session(layer.configuration)

	// we update the session from the hello packet
	newSession.UpdateFromHelloPacket(packet)

	// store the session
	layer.sessions.Store(newSession, clientId)

	// we generate a new challenge nonce, and use it to generate a challenge packet,
	// send the packet to the client, and store the nonce for later check
	if newChallenge, err = helpers.GenerateChallengeNonce(); err == nil {
		newSession.SetChallenge(&newChallenge)
		var newChallengePacket *dto.ChallengePacket = nil

		if newSession.Encryption.GetEncryptionAlgorithm() == enums.NoEncryptionAlgorithm {
			newChallengePacket = dto.GenerateChallengePacket(clientId, &newChallenge)
		} else {
			// we generate the ephemeral keys and check if everything is ok so far
			if err := newSession.Encryption.GenerateEphemeralKeys(); err != nil {
				return err
			}

			// we generate the challenge packet with the challenge and the server ephemeral public key
			newChallengePacket = dto.GenerateChallengePacketWithServerPublicEphemeralKey(
				clientId,
				&newChallenge,
				newSession.Encryption.GetEphemeralKeys().GetPublicKey(),
			)
		}

		if err2 := layer.layer1.SendPacket(newChallengePacket); err2 == nil {
			newSession.SetState(enums.CHALLENGE_SENT)
			return nil
		} else {
			return err2
		}
	}

	// if the random challenge generation failed, then we go for a manual one
	defaultChallengeNonce := helpers.GetDefaultChallengeNonce()

	newSession.SetChallenge(&defaultChallengeNonce)
	newChallengePacket := dto.GenerateChallengePacket(clientId, &defaultChallengeNonce)

	if err := layer.layer1.SendPacket(newChallengePacket); err == nil {
		newSession.SetState(enums.CHALLENGE_SENT)
	} else {
		return err
	}

	return nil
}

// ==== handler of packets from layer 3 ====

func (layer *Layer2Leader) handlePacketFromLayer3(packet packets.Packet) errorinfo.GatewayError {
	return nil
}

// ==== handling authentication ====

// handles a client that has just finished proving its identity successfully
func (layer *Layer2Leader) authorizeSession(sessionId int64) {
	if sessionToApprove, sessionExist := layer.sessions.GetExists(sessionId); sessionExist {
		switch sessionToApprove.GetState() {
		case enums.SESSION_CONNECTED:
			packet := dto.CreateNewConnectedPacket(sessionId, sessionToApprove)
			layer.layer1.SendPacket(packet)

		case enums.CHALLENGE_SENT, enums.RECEIVED_CONNECT: // the user just proved identity
			sessionToApprove.SetState(enums.RECEIVED_CONNECT)

			currentProcessData := helpers.GetResourceUsageOfProcessNoError(layer.pid)

			var selectedIndex int64 = 0

			messagesCount := layer.clusterServer.MessagesCounter.Load()
			sessionsCount := int32(layer.sessions.Len())
			cpuUsage := float32(currentProcessData.CPUusage)
			ramUsage := currentProcessData.RAMusage

			currentPoint := helpers.CalculateClusterMemberWeight(
				messagesCount,
				sessionsCount,
				cpuUsage,
				ramUsage,
				true,
			) * 1.01 // yes, intentionally, we make the leader heavier than the followers
			tempPoint := currentPoint

			layer.clusterServer.Clients.Iterate(func(index int64, data *clusters.ClusterFollowerContainer) {
				messagesCount = data.MessagesSinceLastCheck()
				sessionsCount = data.ActiveSessions()
				cpuUsage = data.CPUpercentUsage()
				ramUsage = data.RAMpercentUsage()

				tempPoint = helpers.CalculateClusterMemberWeight(
					messagesCount,
					sessionsCount,
					cpuUsage,
					ramUsage,
					data.IsHealthy(),
				)

				if tempPoint <= currentPoint {
					currentPoint = tempPoint
					selectedIndex = index
				}

			})

			if selectedIndex != 0 { // it has to be redirected
				if member, memberExists := layer.clusterServer.Clients.GetExists(selectedIndex); memberExists {
					if forwardData, err := member.Client.RequestAcceptClient(
						sessionId,
						sessionToApprove.GetFrame(),
					); err == nil {
						packet := dto.CreateNewRedirectPacket(sessionId, forwardData)
						layer.layer1.SendPacket(packet)
						sessionToApprove.SetRedirectPacket(packet)
						sessionToApprove.SetState(enums.REDIRECTING)

						return
					}
				}
			}

			// if the client reaches this point, it stays
			layer.approveSession(sessionId)

		default: // if the status is none of the expected ones
			layer.closeSession(sessionId, enums.CloseReasonProtocolViolation)
		}

	}
}

// ==== constructor ====
func CreateNewLayer2Leader(conf *config.Configuration) *Layer2Leader {
	answer := Layer2Leader{}
	core := CreateNewLayer2Core(
		conf,
		answer.handlePacketFromLayer1,
		answer.handlePacketFromLayer3,
		answer.initializeCluster,
		answer.stopCluster,
	)

	answer.Layer2Core = core
	answer.clusterServer = clustering.CreateClusteringLeader(core.wg, conf, &answer)
	answer.pid = int32(os.Getpid())

	return &answer
}
