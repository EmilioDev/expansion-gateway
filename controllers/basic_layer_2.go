// file: /controllers/basic_layer_2.go
package controllers

import (
	"crypto/ed25519"
	"expansion-gateway/config"
	"expansion-gateway/dto"
	"expansion-gateway/enums"
	"expansion-gateway/errors/layererrors"
	"expansion-gateway/helpers"
	disp "expansion-gateway/interfaces/dispatchers"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/layers"
	"expansion-gateway/interfaces/packets"
	"expansion-gateway/internal/others"
	"expansion-gateway/internal/structs"
	"sync"
	"sync/atomic"
	"time"
)

type BasicLayer2 struct {
	layer1        layers.Layer1
	layer3        layers.Layer3
	working       *atomic.Bool
	configuration *config.Configuration
	layer1Reciver disp.Reciver
	sessions      *structs.SessionsDictionary[*dto.Layer2Session]
	wg            *sync.WaitGroup
}

func (layer BasicLayer2) ConfigureFirstLayer(target layers.Layer1) errorinfo.GatewayError {
	layer.layer1 = target

	dispatcher, reciver := others.NewShardedDispatcher(layer.configuration)

	layer.layer1Reciver = reciver

	return layer.layer1.ConfigureDumbLayer(dispatcher)
}

func (layer BasicLayer2) ConfigureThirdLayer(target layers.Layer3) errorinfo.GatewayError {
	layer.layer3 = target
	return nil
}

func (layer BasicLayer2) IsWorking() bool {
	return layer.working.Load()
}

func (layer BasicLayer2) Start() errorinfo.GatewayError {
	if layer.working.Load() {
		return nil
	}

	if layer.layer1 == nil || layer.layer3 == nil {
		return layererrors.CreateDumbLayersNotConfigured_LayerError(
			"/controllers/basic_layer_2.go",
			53,
			enums.LAYER_2,
			layer.layer1,
			layer.layer3)
	}

	// Start Layer 1
	if err := layer.layer1.Start(); err != nil {
		return err
	}

	// Start Layer 3 (if applicable)
	if err := layer.layer3.Start(); err != nil {
		return err
	}

	layer.working.Store(true)

	layer.initializeLayer1Listeners()
	layer.initializeLayer3Listeners()

	// start session timeout manager
	layer.wg.Add(1)
	go layer.sessionTimeoutWatcher()

	return nil
}

func (layer BasicLayer2) Stop() errorinfo.GatewayError {
	layer.working.Store(false)

	if layer.layer1 != nil {
		if err := layer.layer1.Stop(); err != nil {
			return err
		}
	}

	if layer.layer3 != nil {
		if err := layer.layer3.Stop(); err != nil {
			return err
		}
	}

	layer.wg.Wait()

	return nil
}

func (layer *BasicLayer2) initializeLayer1Listeners() {
	shards := layer.layer1Reciver.ShardCount()

	for x := 0; x < shards; x++ {
		layer.wg.Add(1)
		go layer.listenLayer1(x)
	}
}

func (layer *BasicLayer2) initializeLayer3Listeners() {
	// Reserved for later
}

// ==== packet listeners ====

// layer 1 packet listener
func (layer *BasicLayer2) listenLayer1(shardIndex int) {
	channel := layer.layer1Reciver.GetShard(shardIndex)
	defer layer.wg.Done()

	for layer.IsWorking() {
		select {
		case packet, ok := <-channel:
			if !ok {
				return
			}

			if err := layer.handlePacketFromLayer1(packet); err != nil {
				sessionToClose := packet.GetSender()

				switch err.GetErrorCode() {
				case 13: // protocol violation
					layer.closeSession(sessionToClose, enums.CloseReasonProtocolViolation)

				case 8, 9, 10, 11, 12:
					layer.closeSession(sessionToClose, enums.CloseReasonInternalError)

				case 0, 1, 2, 3, 4, 5, 6: // packet error
					layer.closeSession(sessionToClose, enums.CloseReasonInvalidPacket)

				case 7: // external error
					fallthrough
				default:
					layer.closeSession(sessionToClose, enums.CloseReasonUnknown)
				}
			}

		default:
			time.Sleep(time.Millisecond * 10) // Yield CPU, prevent tight loop
		}
	}
}

// ==== layer 1 packet handlers ====

// global packet handler from layer 1
func (layer *BasicLayer2) handlePacketFromLayer1(packet packets.Packet) errorinfo.GatewayError {
	const filePath string = "/controllers/basic_layer_2.go"

	switch packet.GetPacketType() {
	case enums.HELLO:
		if helloPacket, ok := packet.(*dto.HelloPacket); ok {
			return layer.handleHelloPacket(helloPacket)
		}

		return layererrors.CreateProtocolFlowViolation_LayerError(
			filePath,
			171,
			enums.LAYER_2,
			enums.INCORRECT_PACKET_KIND)

	case enums.CHALLENGE: // clients should never send a challenge
		return layererrors.CreateProtocolFlowViolation_LayerError(
			filePath,
			178,
			enums.LAYER_2,
			enums.CLIENT_SENT_CHALLENGE)

	case enums.NONE: //the layer 1 received an invalid packet
		return layererrors.CreateProtocolFlowViolation_LayerError(
			filePath,
			185,
			enums.LAYER_2,
			enums.INVALID_PACKET)

	case enums.CONNECT:
		if connectPacket, ok := packet.(*dto.ConnectPacket); ok {
			return layer.handleConnectPacket(connectPacket)
		}

		return layererrors.CreateProtocolFlowViolation_LayerError(
			filePath,
			196,
			enums.LAYER_2,
			enums.INCORRECT_PACKET_KIND)
	}

	return nil
}

// layer 1 connect packet handler
func (layer *BasicLayer2) handleConnectPacket(packet *dto.ConnectPacket) errorinfo.GatewayError {
	sessionId := packet.GetSender()
	const filePath string = "/controllers/basic_layer_2.go"

	if session, sessionExists := layer.sessions.GetExists(sessionId); sessionExists {
		session.RefreshActivity()

		if state := session.GetState(); state == enums.CHALLENGE_SENT || state == enums.RECEIVED_CONNECT {
			publicKey := session.GetEd25519PublicKey()
			challenge := session.GetChallenge()

			if ok := ed25519.Verify(publicKey, challenge, packet.Signature[:]); ok {
				layer.authorizeSession(sessionId) // handled here
			} else {
				// this client is unauthorized!!!
				layer.closeSession(sessionId, enums.CloseReasonFailedAuthentication) // handled
			}

			return nil
		}

		// i think, instead of doing this, we should give the connect an extra functionality:
		// if the user wants its datas again, it could just send a connect
		return layererrors.CreateProtocolFlowViolation_LayerError(
			filePath,
			231,
			enums.LAYER_2,
			enums.CLIENT_SENT_CONNECT_AT_WRONG_MOMENT)

	}

	return layererrors.CreateProtocolFlowViolation_LayerError(
		filePath,
		239,
		enums.LAYER_2,
		enums.SESSION_CLOSED)
}

// layer 1 hello packet handler
func (layer *BasicLayer2) handleHelloPacket(packet *dto.HelloPacket) errorinfo.GatewayError {
	clientId := packet.GetSender()
	const filePath string = "/controllers/basic_layer_2.go"
	var newChallenge []byte
	var err errorinfo.GatewayError = nil

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

			newChallengePacket := dto.GenerateChallengePacket(clientId, &newChallenge)

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
	newSession := dto.GenerateNewLayer2Session(layer.configuration)

	// we update the session from the hello packet
	newSession.UpdateFromHelloPacket(packet)

	// store the session
	layer.sessions.Store(newSession, clientId)

	// we generate a new challenge nonce, and use it to generate a challenge packet,
	// send the packet to the client, and store the nonce for later check
	if newChallenge, err = helpers.GenerateChallengeNonce(); err == nil {
		newSession.SetChallenge(&newChallenge)
		newChallengePacket := dto.GenerateChallengePacket(clientId, &newChallenge)

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

// ==== error handlers ====

// invalid packet handler
func (layer *BasicLayer2) closeSession(sessionId int64, reason enums.SessionCloseReason) {
	if layer.layer1 != nil {
		// we need to send the disconnect packet first
		// do not forget to add it!!!
		layer.layer1.CloseSession(sessionId)
	}

	if layer.layer3 != nil {
		// the same as in layer 1
		layer.layer3.CloseSession(sessionId)
	}

	layer.sessions.Delete(sessionId)
}

// ==== timeout watcher ====

func (layer *BasicLayer2) sessionTimeoutWatcher() {
	defer layer.wg.Done()

	checkPeriod := layer.configuration.GetSessionWatcherPeriod()

	if checkPeriod <= 0 {
		checkPeriod = time.Second
	}

	ticker := time.NewTicker(checkPeriod)
	defer ticker.Stop()

	for range ticker.C {
		if !layer.IsWorking() {
			return
		}

		keys := layer.sessions.Keys()

		for _, key := range keys {
			if session, exists := layer.sessions.GetExists(key); exists {
				if session.TimeoutTracker().Expired() {
					if session.GetState() == enums.CHALLENGE_SENT {
						layer.closeSession(key, enums.CloseReasonChallengeTimeout)
					} else {
						layer.closeSession(key, enums.CloseReasonIdleTimeout)
					}
				}
			}
		}
	}
}

// ==== handling authentication ====

// handles a client that has just finished proving its identity successfully
func (layer *BasicLayer2) authorizeSession(sessionId int64) {
	//
}

// ==== constructor ====
func CreateNewBasicLayer2(conf *config.Configuration) *BasicLayer2 {
	var working atomic.Bool
	working.Store(false)

	return &BasicLayer2{
		layer1:        nil,
		layer3:        nil,
		configuration: conf,
		working:       &working,
		sessions:      structs.CreateNewSessionDictionary[*dto.Layer2Session](),
		wg:            &sync.WaitGroup{},
	}
}
