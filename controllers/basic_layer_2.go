// file: /controllers/basic_layer_2.go
package controllers

import (
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

	return nil
}

func (layer *BasicLayer2) initializeLayer1Listeners() {
	shards := layer.layer1Reciver.ShardCount()

	for x := 0; x < shards; x++ {
		go layer.listenLayer1(x)
	}
}

func (layer *BasicLayer2) initializeLayer3Listeners() {
	// Reserved for later
}

func (layer *BasicLayer2) listenLayer1(shardIndex int) {
	channel := layer.layer1Reciver.GetShard(shardIndex)

	for layer.IsWorking() {
		select {
		case packet, ok := <-channel:
			if !ok {
				return
			}

			layer.handlePacketFromLayer1(packet)

		default:
			time.Sleep(time.Millisecond * 10) // Yield CPU, prevent tight loop
		}
	}
}

func (layer *BasicLayer2) handlePacketFromLayer1(packet packets.Packet) errorinfo.GatewayError {
	const filePath string = "/controllers/basic_layer_2.go"

	switch packet.GetPacketType() {
	case enums.HELLO:
		helloPacket, ok := packet.(*dto.HelloPacket)

		if ok {
			return layer.handleHelloPacket(helloPacket)
		}
		// handle later a not hello packet, but that should not occur
		// if the architecture is respected

	case enums.CHALLENGE: // clients should never send a challenge
		return layererrors.CreateProtocolFlowViolation_LayerError(
			filePath,
			141,
			enums.LAYER_2,
			enums.CLIENT_SENT_CHALLENGE)

	case enums.NONE: //the layer 1 received an invalid packet
		return layererrors.CreateProtocolFlowViolation_LayerError(
			filePath,
			148,
			enums.LAYER_2,
			enums.INVALID_PACKET)
	}

	return nil
}

func (layer *BasicLayer2) handleHelloPacket(packet *dto.HelloPacket) errorinfo.GatewayError {
	clientId := packet.GetSender()
	const filePath string = "/controllers/basic_layer_2.go"
	var newChallenge []byte
	var err errorinfo.GatewayError = nil

	// if the session exists
	if sessionStored, sessionExist := layer.sessions.GetExists(clientId); sessionExist {
		// if the state of the session is not CHALLENGE_SENT, then this is an invalid packet
		if sessionStored.GetState() != enums.CHALLENGE_SENT {
			// apply the corresponding measure for sending invalid packet
			return layererrors.CreateProtocolFlowViolation_LayerError(filePath, 150, enums.LAYER_2, enums.INVALID_HELLO)
		}

		// we update the session from the hello packet
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

		return layer.layer1.SendPacket(newChallengePacket)
	}

	// the session does not exist
	// then we generate a new one
	newSession := dto.GenerateNewLayer2Session()

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

// constructor
func CreateNewBasicLayer2(conf *config.Configuration) *BasicLayer2 {
	var working atomic.Bool
	working.Store(false)

	return &BasicLayer2{
		layer1:        nil,
		layer3:        nil,
		configuration: conf,
		working:       &working,
		sessions:      structs.CreateNewSessionDictionary[*dto.Layer2Session](),
	}
}
