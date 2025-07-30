package kcp_handler

import (
	"expansion-gateway/config"
	"expansion-gateway/dto"
	"expansion-gateway/enums"
	"expansion-gateway/errors"
	"expansion-gateway/errors/layererrors"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/commands"
	"expansion-gateway/interfaces/dispatchers"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/parsers"
	"fmt"
	"sync"

	kcp "github.com/xtaci/kcp-go/v5"
)

type KcpAsLayer1 struct {
	running          bool                             // If this layer is running or not. This is just a flag
	sessions         map[int64]*kcp.UDPSession        // all the sessions that are currently active
	listener         *kcp.Listener                    // the kcp listener
	outputDispatcher dispatchers.Dispatcher           // this is the channel used to communicate with the next layer
	inputChannel     <-chan commands.Command          // the commands this layer will recive from the 2nd layer(remember, that is the decission-making layer)
	configuration    *config.Configuration            // this is the configuration module. it contains all the config details
	parser           parsers.ByteStreamToPacketParser // the byte array to packet parser
	sessionsMutex    *sync.RWMutex
}

func CreateNewKcpLayer1(configuration *config.Configuration,
	parser parsers.ByteStreamToPacketParser) *KcpAsLayer1 {
	return &KcpAsLayer1{
		// outputChannel: outputChannel,
		running:       false,
		sessions:      make(map[int64]*kcp.UDPSession),
		listener:      nil,
		configuration: configuration,
		parser:        parser,
		sessionsMutex: &sync.RWMutex{},
		// inputChannel:  inputChannel,

		// the channels are not going to be assigned here, this layer will be passed to layer 2
		// and layer 2 will use the "ConfigureDumbLayer" method to configure the channels
		// with the ones created there
	}
}

func (layer KcpAsLayer1) Start() errorinfo.GatewayError {
	const filePath string = "/kcp_handler/kcp_as_layer1.go"
	if layer.running {
		return nil
	}

	if layer.outputDispatcher == nil {
		return layererrors.CreateChannelClosed_LayerError(filePath, 50, enums.LAYER_1, enums.OUTPUT_CHANNEL)
	}

	var serverPath string = layer.configuration.GetServerAddress()

	if listener, err := kcp.ListenWithOptions(serverPath, nil, 10, 3); err == nil {
		layer.running = true
		fmt.Printf("server running on %s\n", serverPath)

		layer.listener = listener

		go layer.process()
		go layer.listenFromInputChannel()

		return nil
	} else {
		return helpers.WithStackTrace(errors.CreateErrorWrapper(filePath, 66, err), 2)
	}
}

func (layer KcpAsLayer1) Stop() errorinfo.GatewayError {
	layer.running = false
	const filePath string = "/kcp_handler/kcp_as_layer1.go"

	if layer.outputDispatcher == nil {
		return layererrors.CreateChannelClosed_LayerError(filePath, 75, enums.LAYER_1, enums.OUTPUT_CHANNEL)
	}

	// close(layer.outputDispatcher)

	return nil
}

func (layer KcpAsLayer1) ConfigureDumbLayer(outputDispatcher dispatchers.Dispatcher, inputChannel <-chan commands.Command) errorinfo.GatewayError {
	layer.outputDispatcher = outputDispatcher
	layer.inputChannel = inputChannel
	return nil
}

func (layer *KcpAsLayer1) listenFromInputChannel() {
	// for now, still pending...
}

func (layer *KcpAsLayer1) process() {
	for layer.running {
		if session, err := layer.listener.AcceptKCP(); err == nil {
			connectionId := helpers.GenerateRandomInt64()

			layer.sessionsMutex.Lock()
			for {
				if _, exists := layer.sessions[connectionId]; !exists {
					break
				}

				connectionId = helpers.GenerateRandomInt64()
			}

			layer.sessions[connectionId] = session

			layer.sessionsMutex.Unlock()

			go layer.handleSession(connectionId)
		}
	}
}

func (layer *KcpAsLayer1) handleSession(connectionId int64) {
	buffer := make([]byte, layer.configuration.GetBufferSize())

	var session *kcp.UDPSession = nil
	sessionExists := true

	for {
		if layer.outputDispatcher == nil {
			layer.sessions[connectionId].Close()
			delete(layer.sessions, connectionId)
			return
		}

		layer.sessionsMutex.RLock()
		session, sessionExists = layer.sessions[connectionId]
		layer.sessionsMutex.RUnlock()

		if sessionExists {
			if dataLen, err := session.Read(buffer); err == nil {
				rawPacket := buffer[:dataLen]

				if packet, err := layer.parser.ParseByteArrayToPacket(&rawPacket, connectionId); err == nil {
					layer.outputDispatcher.Dispatch(packet)
				}
			} else {
				fmt.Printf("error in session %d: %s\n", connectionId, err.Error())
				layer.outputDispatcher.Dispatch(dto.CreateInvalidPacket(connectionId))
			}
		} else {
			return
		}
	}
}
