package kcp_handler

import (
	"expansion-gateway/config"
	"expansion-gateway/dto"
	"expansion-gateway/enums"
	"expansion-gateway/errors"
	"expansion-gateway/errors/layererrors"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/dispatchers"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
	"expansion-gateway/interfaces/parsers"
	"expansion-gateway/internal/structs"
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
	outputDispatcher dispatchers.Dispatcher                       // this is the channel used to communicate with the next layer
	configuration    *config.Configuration                        // this is the configuration module. it contains all the config details
	parser           parsers.ByteStreamToPacketParser             // the byte array to packet parser
	working          *atomic.Bool                                 // tells you if this layer is working or not
	shutdownOnce     *sync.Once                                   // executes the shutdown only once
	startOnce        *sync.Once                                   //executes the start only once
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

		// the channels are not going to be assigned here, this layer will be passed to layer 2
		// and layer 2 will use the "ConfigureDumbLayer" method to configure the channels
		// with the ones created there
	}

	result.working.Store(false)

	return result
}

func (layer *KcpAsLayer1) Start() errorinfo.GatewayError {
	const filePath string = "/kcp_handler/kcp_as_layer1.go"
	if layer.IsWorking() {
		return nil
	}

	if layer.outputDispatcher == nil {
		return layererrors.CreateChannelClosed_LayerError(filePath, 50, enums.LAYER_1, enums.OUTPUT_CHANNEL)
	}

	var result errorinfo.GatewayError = nil

	layer.startOnce.Do(func() {
		var serverPath string = layer.configuration.GetServerAddress()

		if listener, err := kcp.ListenWithOptions(serverPath, nil, 10, 3); err == nil {
			layer.working.Store(true)
			fmt.Printf("server running on %s\n", serverPath)

			layer.listener = listener

			go layer.process()
			go layer.listenFromInputChannel()
		} else {
			result = helpers.WithStackTrace(errors.CreateErrorWrapper(filePath, 66, err), 2)
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

func (layer *KcpAsLayer1) ConfigureDumbLayer(outputDispatcher dispatchers.Dispatcher) errorinfo.GatewayError {
	layer.outputDispatcher = outputDispatcher
	return nil
}

func (layer *KcpAsLayer1) SendPacket(packet packets.Packet) errorinfo.GatewayError {
	const filePath string = "/kcp_handler/kcp_as_layer1.go"

	if !layer.IsWorking() {
		return layererrors.CreateLayerClosed_LayerError(filePath, 126, enums.LAYER_1)
	}

	sessionId := packet.GetSender()

	session, exists := layer.sessions.GetExists(sessionId)

	if !exists {
		return layererrors.CreateSessionNotRegistered_LayerError(filePath, 136, enums.LAYER_1, sessionId)
	}

	if byteArray, err := packet.Marshal(); err == nil {
		session.SetWriteDeadline(time.Now().Add(2 * time.Second)) // timeout

		if _, err := session.Write(byteArray); err != nil {
			return helpers.WithStackTrace(errors.CreateErrorWrapper(filePath, 143, err), 0)
		}
	} else {
		return err
	}

	return nil
}

func (layer *KcpAsLayer1) CloseSession(sessionId int64) errorinfo.GatewayError {
	if connection, exists := layer.sessions.GetExists(sessionId); exists {
		layer.sessions.Delete(sessionId)
		connection.Close()
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

func (layer *KcpAsLayer1) listenFromInputChannel() {
	// for now, still pending...
}

func (layer *KcpAsLayer1) process() {
	for layer.IsWorking() {
		if session, err := layer.listener.AcceptKCP(); err == nil {
			connectionId := layer.sessions.Add(session)

			layer.wg.Add(1)
			go layer.handleSession(connectionId)
		}
	}
}

func (layer *KcpAsLayer1) handleSession(connectionId int64) {
	defer layer.wg.Done()

	buffer := make([]byte, layer.configuration.GetBufferSize())
	timeoutDuration := time.Duration(layer.configuration.GetConnectionTimeout()) * time.Second

	for {
		if session, sessionExists := layer.sessions.GetExists(connectionId); sessionExists {
			if session == nil {
				return
			}

			if !layer.IsWorking() {
				session.Close()
				return
			}

			if err := session.SetReadDeadline(time.Now().Add(timeoutDuration)); err != nil {
				continue
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
				}

				fmt.Printf("error in session %d: %s\n", connectionId, err.Error())
				layer.outputDispatcher.Dispatch(dto.CreateInvalidPacket(connectionId))
			}
		} else {
			return
		}
	}
}
