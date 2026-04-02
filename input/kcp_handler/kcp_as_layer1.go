package kcp_handler

import (
	"expansion-gateway/config"
	"expansion-gateway/dto"
	"expansion-gateway/enums"
	"expansion-gateway/errors"
	"expansion-gateway/errors/layererrors"
	sessionsErrors "expansion-gateway/errors/sessions"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/dispatchers"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
	"expansion-gateway/interfaces/parsers"
	structs "expansion-gateway/internal/structs/dictionaries"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	kcp "github.com/xtaci/kcp-go/v5"
)

type KcpAsLayer1 struct {
	sessions         *structs.SessionsDictionary[*kcp.UDPSession] // all the sessions that are currently active
	listener         *kcp.Listener                                // the kcp listener
	outputDispatcher dispatchers.Dispatcher[packets.Packet]       // this is the channel used to communicate with the next layer
	inputReceiver    dispatchers.Reciver[packets.Packet]          // receiver from layer 1
	configuration    *config.Configuration                        // this is the configuration module. it contains all the config details
	parser           parsers.ByteStreamToPacketParser             // the byte array to packet parser
	working          *atomic.Bool                                 // tells you if this layer is working or not
	shutdownOnce     *sync.Once                                   // executes the shutdown only once
	startOnce        *sync.Once                                   // executes the start only once
	wg               *sync.WaitGroup                              // the wait group of layer 1
}

func CreateNewKcpLayer1(configuration *config.Configuration,
	parser parsers.ByteStreamToPacketParser) *KcpAsLayer1 {
	result := &KcpAsLayer1{
		sessions:      structs.CreateNewSessionDictionary[*kcp.UDPSession](),
		listener:      nil,
		configuration: configuration,
		parser:        parser,
		working:       &atomic.Bool{},
		shutdownOnce:  &sync.Once{},
		startOnce:     &sync.Once{},
		wg:            &sync.WaitGroup{},
	}

	result.working.Store(false)

	return result
}

func (layer *KcpAsLayer1) Start() errorinfo.GatewayError {
	const filePath string = "/input/kcp_handler/kcp_as_layer1.go"
	if layer.IsWorking() {
		return nil
	}

	if layer.outputDispatcher == nil {
		return layererrors.CreateChannelClosed_LayerError(filePath, 50, enums.LAYER_1, enums.OUTPUT_CHANNEL)
	}

	var result errorinfo.GatewayError = nil

	layer.startOnce.Do(func() {
		var serverPath string = layer.configuration.GetKcpPathToThisGateway()
		var universalPath string = layer.configuration.GetUniversalKcpPathToThisGateway()

		if listener, err := kcp.ListenWithOptions(universalPath, nil, 10, 3); err == nil {
			layer.working.Store(true)

			layer.listener = listener

			go layer.process()
			go layer.initializeLayer2Listeners()

			if layer.configuration.AreWeClusterLeaders() {
				fmt.Printf("server running on %s\n", serverPath)
			}
		} else {
			result = helpers.WithStackTrace(errors.CreateErrorWrapper(filePath, 82, err), 2)
		}
	})

	return result
}

// IsWorking reports whether the layer is still accepting connections.
func (layer *KcpAsLayer1) IsWorking() bool {
	return layer.working.Load()
}

func (layer *KcpAsLayer1) Stop() errorinfo.GatewayError {
	if !layer.IsWorking() {
		return nil
	}

	layer.shutdownOnce.Do(func() {
		layer.working.Store(false)

		// closing the listener
		if layer.listener != nil {
			_ = layer.listener.Close()
		}

		// closing the sessions
		layer.sessions.Clear()

		fmt.Println("Gateway exited")
	})

	return nil
}

func (layer *KcpAsLayer1) ConfigureDumbLayer(outputChannel dispatchers.Dispatcher[packets.Packet], inputChannel dispatchers.Reciver[packets.Packet]) errorinfo.GatewayError {
	layer.outputDispatcher = outputChannel
	layer.inputReceiver = inputChannel

	return nil
}

func (layer *KcpAsLayer1) SendPacket(packet packets.Packet) errorinfo.GatewayError {
	const filePath string = "/input/kcp_handler/kcp_as_layer1.go" // this is for error handling

	if !layer.IsWorking() { // this if is for checking if we're still up
		return layererrors.CreateLayerClosed_LayerError(filePath, 126, enums.LAYER_1)
	}

	sessionId := packet.GetSender() // this will give you the id you need to find the *kcp.UDPSession

	if session, exists := layer.sessions.GetExists(sessionId); exists && session != nil { // here you have the session of type *kcp.UDPSession
		if byteArray, err := packet.Marshal(); err == nil { // here you have the byte array you will send
			if _, err := session.Write(byteArray); err != nil { // and then you send the data here
				return helpers.WithStackTrace(errors.CreateErrorWrapper(filePath, 134, err), 0) // if error on sending, this happen
			}
		} else {
			return err
		}
	} else {
		return sessionsErrors.CreateInvalidSessionError(filePath, 142, sessionId)
	}

	return nil
}

func (layer *KcpAsLayer1) CloseSession(sessionId int64) errorinfo.GatewayError {
	if connection, exists := layer.sessions.GetExists(sessionId); exists {
		layer.sessions.Delete(sessionId)

		if connection != nil {
			connection.Close()
		}
	}

	return nil
}

func (layer *KcpAsLayer1) MoveClientTo(origin, destiny int64) {
	layer.sessions.MoveTo(origin, destiny)
}

func (layer *KcpAsLayer1) DisableSession(sessionId int64) {
	if session, exists := layer.sessions.GetExists(sessionId); exists && session != nil {
		session.Close()
		layer.sessions.Store(nil, sessionId)
	}
}

func (layer *KcpAsLayer1) process() {
	for layer.IsWorking() {
		if session, err := layer.listener.AcceptKCP(); err == nil {
			session.SetACKNoDelay(true)
			// session.SetStreamMode(true)
			session.SetWindowSize(256, 256)
			session.SetNoDelay(1, 10, 2, 1)
			session.SetMtu(1200)

			connectionId := layer.sessions.Add(session)

			layer.wg.Add(1)
			go layer.handleSession(connectionId)
		}
	}
}

func (layer *KcpAsLayer1) handleSession(connectionId int64) {
	defer layer.wg.Done()
	buffer := make([]byte, layer.configuration.GetBufferSize())

	for {
		if session, sessionExists := layer.sessions.GetExists(connectionId); sessionExists {
			if session == nil {
				return
			}

			if !layer.IsWorking() {
				session.Close()
				return
			}

			if dataLen, err := session.Read(buffer); err == nil {
				rawPacket := buffer[:dataLen]

				if packet, err := layer.parser.ParseByteArrayToPacket(&rawPacket, connectionId); err == nil {
					layer.outputDispatcher.Dispatch(packet)
				}
			} else {
				// check if the error is a timeout
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// it is just a timeout, ignore it
					continue
				} else if err.Error() == "io: read/write on closed pipe" {
					// the connection is closed
					layer.outputDispatcher.Dispatch(dto.CreateDisconnectPacket(
						connectionId,
						enums.CloseReasonConnectionLost,
					))

					return
				}

				fmt.Printf("kcp layer error in session %d: %s\n", connectionId, err.Error())
				layer.outputDispatcher.Dispatch(dto.CreateInvalidPacket(connectionId))
			}
		} else {
			return
		}
	}
}

func (layer *KcpAsLayer1) initializeLayer2Listeners() {
	shards := layer.inputReceiver.ShardCount()

	for x := 0; x < shards; x++ {
		layer.wg.Add(1)
		go layer.handlePacketsFromLayer2(x)
	}
}

func (layer *KcpAsLayer1) handlePacketsFromLayer2(shardIndex int) {
	defer layer.wg.Done()
	channel := layer.inputReceiver.GetShard(shardIndex)

	for layer.IsWorking() {
		select {
		case packet, ok := <-channel:
			if !ok {
				return
			}

			if packet == nil {
				continue
			} else {
				layer.SendPacket(packet)
			}

		default:
			time.Sleep(time.Millisecond) // prevents tight loop
		}
	}
}
