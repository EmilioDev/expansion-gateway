// file: /controllers/layer_2_core.go
package controllers

import (
	"expansion-gateway/config"
	dto "expansion-gateway/dto/sessions"
	"expansion-gateway/enums"
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

type Layer2Core struct {
	layer1                    layers.Layer1
	layer3                    layers.Layer3
	working                   *atomic.Bool
	configuration             *config.Configuration
	layer1Reciver             disp.Reciver
	sessions                  *structs.SessionsDictionary[*dto.Layer2Session]
	wg                        *sync.WaitGroup
	layer1PacketHandler       func(packets.Packet) errorinfo.GatewayError
	layer3PacketHandler       func(packets.Packet) errorinfo.GatewayError
	initializeClusterCallback func() errorinfo.GatewayError
	stopClusterCallback       func() errorinfo.GatewayError
}

// starts the server
func (layer *Layer2Core) Start() errorinfo.GatewayError {
	if layer.working.Load() {
		return nil
	}

	// if layer.layer1 == nil || layer.layer3 == nil {
	// 	return layererrors.CreateDumbLayersNotConfigured_LayerError(
	// 		"/controllers/layer_2_core.go",
	// 		40,
	// 		enums.LAYER_2,
	// 		layer.layer1,
	// 		layer.layer3)
	// }

	// remove those ifs for checking if the layer exists, and uncomment the previous one
	// Start Layer 1
	if layer.layer1 != nil {
		if err := layer.layer1.Start(); err != nil {
			return err
		}
	}

	// Start Layer 3 (if applicable)
	if layer.layer3 != nil {
		if err := layer.layer3.Start(); err != nil {
			return err
		}
	}

	layer.working.Store(true)

	layer.initializeLayer1Listeners()
	layer.initializeLayer3Listeners()

	// start session timeout manager
	layer.wg.Add(1)
	go layer.sessionTimeoutWatcher()

	return layer.initializeClusterCallback()
}

// stops the server
func (layer *Layer2Core) Stop() errorinfo.GatewayError {
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

	if err := layer.stopClusterCallback(); err != nil {
		return err
	}

	layer.wg.Wait()

	return nil
}

func (layer *Layer2Core) ConfigureFirstLayer(target layers.Layer1) errorinfo.GatewayError {
	layer.layer1 = target

	dispatcher, reciver := others.NewShardedDispatcher(layer.configuration)

	layer.layer1Reciver = reciver

	return layer.layer1.ConfigureDumbLayer(dispatcher)
}

func (layer *Layer2Core) ConfigureThirdLayer(target layers.Layer3) errorinfo.GatewayError {
	layer.layer3 = target
	return nil
}

func (layer *Layer2Core) IsWorking() bool {
	return layer.working.Load()
}

func (layer *Layer2Core) GetActiveSessions() int32 {
	return int32(layer.sessions.Len())
}

func (layer *Layer2Core) HasSession(sessionID int64) bool {
	return layer.sessions.Exists(sessionID)
}

func (layer *Layer2Core) Wait() {
	layer.wg.Wait()
}

// ==== layer 1 handlers

// layer 1 listener
func (layer *Layer2Core) initializeLayer1Listeners() {
	shards := layer.layer1Reciver.ShardCount()

	for x := 0; x < shards; x++ {
		layer.wg.Add(1)
		go layer.listenLayer1(x)
	}
}

// layer 1 packet listener
func (layer *Layer2Core) listenLayer1(shardIndex int) {
	channel := layer.layer1Reciver.GetShard(shardIndex)
	defer layer.wg.Done()

	for layer.IsWorking() {
		select {
		case packet, ok := <-channel:
			if !ok {
				return
			}

			if err := layer.layer1PacketHandler(packet); err != nil {
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

// ==== layer 3 listener

func (layer *Layer2Core) initializeLayer3Listeners() {
	// Reserved for later
}

// ==== close ====

// invalid packet handler
func (layer *Layer2Core) closeSession(sessionId int64, reason enums.SessionCloseReason) {
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

func (layer *Layer2Core) sessionTimeoutWatcher() {
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

// ==== constructor ====

func CreateNewLayer2Core(
	conf *config.Configuration,
	layer1PacketHandler func(packets.Packet) errorinfo.GatewayError,
	layer3PacketHandler func(packets.Packet) errorinfo.GatewayError,
	initializeClusterCallback func() errorinfo.GatewayError,
	stopClusterCallback func() errorinfo.GatewayError,
) *Layer2Core {
	var working atomic.Bool
	working.Store(false)
	wg := &sync.WaitGroup{}

	return &Layer2Core{
		layer1:                    nil,
		layer3:                    nil,
		configuration:             conf,
		working:                   &working,
		sessions:                  structs.CreateNewSessionDictionary[*dto.Layer2Session](),
		wg:                        wg,
		layer1PacketHandler:       layer1PacketHandler,
		layer3PacketHandler:       layer3PacketHandler,
		initializeClusterCallback: initializeClusterCallback,
		stopClusterCallback:       stopClusterCallback,
	}
}
