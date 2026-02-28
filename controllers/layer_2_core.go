// file: /controllers/layer_2_core.go
package controllers

import (
	"expansion-gateway/config"
	"expansion-gateway/dto"
	sessionsDTO "expansion-gateway/dto/sessions"
	"expansion-gateway/enums"
	authErrors "expansion-gateway/errors/auth"
	"expansion-gateway/errors/layererrors"
	disp "expansion-gateway/interfaces/dispatchers"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/layers"
	"expansion-gateway/interfaces/packets"
	"expansion-gateway/internal/others"
	structs "expansion-gateway/internal/structs/dictionaries"
	"expansion-gateway/internal/structs/tries"
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
	layer1Dispatcher          disp.Dispatcher
	layer3Reciver             disp.Reciver
	layer3Dispatcher          disp.Dispatcher
	sessions                  *structs.SessionsDictionary[*sessionsDTO.Layer2Session]
	wg                        *sync.WaitGroup
	layer1PacketHandler       func(packets.Packet) errorinfo.GatewayError
	layer3PacketHandler       func(packets.Packet) errorinfo.GatewayError
	initializeClusterCallback func() errorinfo.GatewayError
	stopClusterCallback       func() errorinfo.GatewayError
	startOnce                 *sync.Once
	stopOnce                  *sync.Once
	subscriptions             *tries.Trie
}

// starts the server
func (layer *Layer2Core) Start() errorinfo.GatewayError {
	if layer.working.Load() {
		return nil
	}

	if layer.layer1 == nil || layer.layer3 == nil {
		return layererrors.CreateDumbLayersNotConfigured_LayerError(
			"/controllers/layer_2_core.go",
			40,
			enums.LAYER_2,
			layer.layer1,
			layer.layer3)
	}

	var result errorinfo.GatewayError = nil

	layer.startOnce.Do(func() {
		// Start Layer 1
		if err := layer.layer1.Start(); err != nil {
			result = err
			return
		}

		// Start Layer 3 (if applicable)
		if err := layer.layer3.Start(); err != nil {
			result = err
			layer.layer1.Stop()

			return
		}

		layer.working.Store(true)

		layer.initializeLayer1Listeners()
		layer.initializeLayer3Listeners()

		// start session timeout manager
		layer.wg.Add(1)
		go layer.sessionTimeoutWatcher()

		result = layer.initializeClusterCallback()
	})

	return result
}

// stops the server
func (layer *Layer2Core) Stop() errorinfo.GatewayError {
	var result errorinfo.GatewayError = nil

	layer.stopOnce.Do(func() {
		layer.working.Store(false)

		if layer.layer1 != nil {
			if err := layer.layer1.Stop(); err != nil {
				result = err
				return
			}
		}

		if layer.layer3 != nil {
			if err := layer.layer3.Stop(); err != nil {
				result = err
				return
			}
		}

		if err := layer.stopClusterCallback(); err != nil {
			result = err
			return
		}

		layer.wg.Wait()
	})

	return result
}

func (layer *Layer2Core) ConfigureFirstLayer(target layers.Layer1) errorinfo.GatewayError {
	layer.layer1 = target

	layer1InDispatcher, layer1InReceiver := others.NewShardedDispatcher(layer.configuration)
	layer1OutDispatcher, layer1OutReceiver := others.NewShardedDispatcher(layer.configuration)

	layer.layer1Reciver = layer1InReceiver
	layer.layer1Dispatcher = layer1OutDispatcher

	return layer.layer1.ConfigureDumbLayer(layer1InDispatcher, layer1OutReceiver)
}

func (layer *Layer2Core) ConfigureThirdLayer(target layers.Layer3) errorinfo.GatewayError {
	layer.layer3 = target

	layer3InDispatcher, layer3InReceiver := others.NewShardedDispatcher(layer.configuration)
	layer3OutDispatcher, layer3OutReceiver := others.NewShardedDispatcher(layer.configuration)

	layer.layer3Reciver = layer3InReceiver
	layer.layer3Dispatcher = layer3OutDispatcher

	return layer.layer3.ConfigureDumbLayer(layer3InDispatcher, layer3OutReceiver)
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

			if packet == nil {
				continue
			} else {
				senderId := packet.GetSender()

				switch packet.GetPacketType() {
				case enums.DISCONNECT:
					if disconnectPacket, isDisconnect := packet.(*dto.DisconnectPacket); isDisconnect {
						layer.closeSession(senderId, disconnectPacket.GetDisconnectReason())
					} else {
						layer.closeSession(senderId, enums.CloseReasonClosedByGateway)
					}

				case enums.PING:
					if session, exists := layer.sessions.GetExists(senderId); exists {
						if session.GetState() == enums.SESSION_CONNECTED {
							session.RefreshActivity()
							layer.sendPacketToLayer1(dto.CreatePingACKpacket(senderId))
						} else {
							layer.closeSession(senderId, enums.CloseReasonProtocolViolation)
						}
					} else {
						layer.closeSession(senderId, enums.CloseReasonConnectionUnauthorized)
					}

				case enums.PINGACK:
					if session, exists := layer.sessions.GetExists(senderId); exists {
						if session.GetPingHasBeenRequested() {
							session.RefreshActivity()
							session.SetPingHasBeenRequested(false)
						} else {
							layer.closeSession(senderId, enums.CloseReasonProtocolViolation)
						}
					} else {
						layer.closeSession(senderId, enums.CloseReasonConnectionUnauthorized)
					}

				case enums.SUBSCRIBE:
					if pkt, ok := packet.(*dto.SubscribePacket); ok {
						if err := layer.handleSubscribePacket(pkt); err != nil {
							layer.closeSession(senderId, enums.CloseReasonProtocolViolation)
						}
					} else {
						layer.closeSession(senderId, enums.CloseReasonGatewayInternalError)
					}

				case enums.UNSUBSCRIBE:
					if pkt, ok := packet.(*dto.UnsubscribePacket); ok {
						if err := layer.handleUnsubscribePacket(pkt); err != nil {
							layer.closeSession(senderId, enums.CloseReasonProtocolViolation)
						}
					} else {
						layer.closeSession(senderId, enums.CloseReasonGatewayInternalError)
					}

				case enums.PUBLISH:
					if pkt, ok := packet.(*dto.PublishPacket); ok {
						if err := layer.handlePublishPacket(pkt); err != nil {
							layer.closeSession(senderId, enums.CloseReasonProtocolViolation)
						}
					} else {
						layer.closeSession(senderId, enums.CloseReasonGatewayInternalError)
					}

				default:
					if err := layer.layer1PacketHandler(packet); err != nil {
						layer.closeSession(senderId, enums.ByteReasonToDisconnectReason(err.GetErrorCode()))
					}
				}
			}

		default:
			time.Sleep(time.Millisecond * 10) // Yield CPU, prevent tight loop
		}
	}
}

// ==== layer 3 listener

func (layer *Layer2Core) initializeLayer3Listeners() {
	shards := layer.layer3Reciver.ShardCount()

	for x := 0; x < shards; x++ {
		layer.wg.Add(1)
		go layer.listenLayer3(x)
	}
}

func (layer *Layer2Core) listenLayer3(shardIndex int) {
	defer layer.wg.Done()

	channel := layer.layer3Reciver.GetShard(shardIndex)

	for layer.IsWorking() {
		select {
		case packet, ok := <-channel:
			if !ok {
				return
			}

			if packet == nil {
				continue
			}

			subscription := tries.SubscriptionKey(packet.GetIdentifier())
			subscribers := layer.subscriptions.GetSubscribers(subscription)

			for _, senderId := range subscribers {
				packet.SetNewOwner(senderId)
				layer.sendPacketToLayer1(packet)
			}

		default:
			time.Sleep(time.Millisecond * 10) // Yield CPU, prevent tight loop
		}
	}
}

// ==== send packet ====

// sends a packet to layer 1
func (layer *Layer2Core) sendPacketToLayer1(packet packets.Packet) {
	layer.layer1Dispatcher.Dispatch(packet)
}

// sends a packet to layer 3
func (layer *Layer2Core) sendPacketToLayer3(packet packets.Packet) {
	layer.layer3Dispatcher.Dispatch(packet)
}

// ==== close ====

// close session handler
func (layer *Layer2Core) closeSession(sessionId int64, reason enums.DisconnectReason) {
	if layer.sessions.WasDeleted(sessionId) {
		layer.subscriptions.RemoveSubscriberFromAllSubscriptions(sessionId)

		layer.closeSessionInLayer1(sessionId, reason)
		layer.closeSessionInLayer3(sessionId, reason)
	}
}

func (layer *Layer2Core) closeSessionInLayer1(sessionId int64, reason enums.DisconnectReason) {
	if layer.IsWorking() {
		layer.layer1.SendPacket(dto.CreateDisconnectPacket(sessionId, reason))
		layer.layer1.CloseSession(sessionId)
	}
}

func (layer *Layer2Core) closeSessionInLayer3(sessionId int64, reason enums.DisconnectReason) {
	if layer.layer3 != nil {
		// the same as in layer 1
		layer.layer3.CloseSession(sessionId)
	}
}

// ===== aprove session =====

// updates a session to the connected state and sends a connect packet to that client
func (layer *Layer2Core) approveSession(sessionId int64) {
	if sessionToApprove, sessionExist := layer.sessions.GetExists(sessionId); sessionExist {
		// we first check if the user has requested a session id
		if requestedSessionId := sessionToApprove.GetRequestedSessionId(); requestedSessionId != 0 {
			if layer.sessions.Exists(requestedSessionId) { // we check if there is another session on that connection
				layer.closeSession(requestedSessionId, enums.CloseReasonSessionIdTakenByOtherConnection)
			}

			// and then we replace
			layer.sessions.MoveTo(sessionId, requestedSessionId)

			if layer.layer1 != nil {
				layer.layer1.MoveClientTo(sessionId, requestedSessionId)
			}

			if layer.layer3 != nil {
				layer.layer3.MoveClientTo(sessionId, requestedSessionId)
			}

			// and we update the session id
			sessionId = requestedSessionId
		}

		packet := dto.CreateNewConnectedPacket(sessionId, sessionToApprove)
		layer.sendPacketToLayer1(packet)

		sessionToApprove.SetState(enums.SESSION_CONNECTED)
	}
}

// ===== redirecting =====

// if this session is marked as redirecting it re-sends the redirect packet to the client
// and returns true, if not, it only returns false and does nothing else.
func (layer *Layer2Core) isRedirecting(sessionId int64) bool {
	if session, sessionExists := layer.sessions.GetExists(sessionId); sessionExists {
		if session.GetState() == enums.REDIRECTING {
			layer.sendPacketToLayer1(session.GetRedirectPacket())
			return true
		}
	}

	return false
}

// ===== timeout watcher =====

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

// ===== pub/sub handlers =====

// subscribes a user in a subscription
func (layer *Layer2Core) handleSubscribePacket(packet *dto.SubscribePacket) errorinfo.GatewayError {
	sender := packet.GetSender()

	if session, exists := layer.sessions.GetExists(sender); exists {
		if session.GetState() != enums.SESSION_CONNECTED {
			return authErrors.CreateConnectionUnauthorizedError("controllers/layer_2_core.go", 351, sender)
		}

		session.RefreshActivity()
		subId := packet.GetSubscriptionID()

		if subId != session.LastSubscribeId.Load() {
			session.LastSubscribeId.Store(subId)
			layer.subscriptions.SubscribeTo(packet.Key, sender)

			layer.sendPacketToLayer1(dto.CreateSubackPacket(sender, subId))
		} else {
			layer.sendPacketToLayer1(dto.CreateSubackPacket(sender, subId))
		}
	}

	return nil
}

// unsubscribes a user to a subscription
func (layer *Layer2Core) handleUnsubscribePacket(packet *dto.UnsubscribePacket) errorinfo.GatewayError {
	sender := packet.GetSender()

	if session, exists := layer.sessions.GetExists(sender); exists {
		if session.GetState() != enums.SESSION_CONNECTED {
			return authErrors.CreateConnectionUnauthorizedError("controllers/layer_2_core.go", 397, sender)
		}

		session.RefreshActivity()
		unsubId := packet.GetUnsubscriptionID()

		if unsubId != session.LastSubscribeId.Load() {
			session.LastSubscribeId.Store(unsubId)
			layer.subscriptions.UnsubscribeTo(packet.Key, sender)

			layer.sendPacketToLayer1(dto.CreateUnsubackPacket(sender, unsubId))
		} else {
			layer.sendPacketToLayer1(dto.CreateUnsubackPacket(sender, unsubId))
		}
	}

	return nil
}

func (layer *Layer2Core) handlePublishPacket(packet *dto.PublishPacket) errorinfo.GatewayError {
	sender := packet.GetSender()
	const filePath string = "controllers/layer_2_core.go"

	if session, exists := layer.sessions.GetExists(sender); exists {
		if session.GetState() != enums.SESSION_CONNECTED {
			return authErrors.CreateConnectionUnauthorizedError(filePath, 481, sender)
		}

		session.RefreshActivity()

		if packet.Key.IsFixedKey() {
			if layer.subscriptions.SubscriptionHasSubscriber(packet.Key, sender) {
				if packet.NeedsAcknowledgement() {
					layer.sendPacketToLayer1(dto.CreatePubackPacket(
						sender,
						packet.GetPublishPacketID(),
						enums.SUCCEED,
					))
				}

				layer.sendPacketToLayer3(packet) // we send with this
			} else if packet.NeedsAcknowledgement() {
				layer.sendPacketToLayer1(dto.CreatePubackPacket(
					sender,
					packet.GetPublishPacketID(),
					enums.USER_NOT_REGISTERED_IN_SUBSCRIPTION,
				))
			}
		} else if packet.NeedsAcknowledgement() {
			layer.sendPacketToLayer1(dto.CreatePubackPacket(
				sender,
				packet.GetPublishPacketID(),
				enums.INVALID_KEY,
			))
		}
	} else {
		layer.closeSession(sender, enums.CloseReasonConnectionUnauthorized)
	}

	return nil
}

// ===== constructor =====

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
		sessions:                  structs.CreateNewSessionDictionary[*sessionsDTO.Layer2Session](),
		wg:                        wg,
		layer1PacketHandler:       layer1PacketHandler,
		layer3PacketHandler:       layer3PacketHandler,
		initializeClusterCallback: initializeClusterCallback,
		stopClusterCallback:       stopClusterCallback,
		startOnce:                 &sync.Once{},
		stopOnce:                  &sync.Once{},
		subscriptions:             tries.CreateTrie(),
	}
}
