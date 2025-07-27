package kcp_handler

import (
	"errors"
	"expansion-gateway/config"
	"expansion-gateway/dto"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/commands"
	"expansion-gateway/interfaces/packets"
	"expansion-gateway/interfaces/parsers"
	"fmt"

	kcp "github.com/xtaci/kcp-go/v5"
)

type KcpAsLayer1 struct {
	running       bool                             // If this layer is running or not. This is just a flag
	sessions      map[int64]*kcp.UDPSession        // all the sessions that are currently active
	listener      *kcp.Listener                    // the kcp listener
	outputChannel chan<- packets.Packet            // this is the channel used to communicate with the next layer
	inputChannel  <-chan commands.Command          // the commands this layer will recive from the 2nd layer(remember, that is the decission-making layer)
	configuration *config.Configuration            // this is the configuration module. it contains all the config details
	parser        parsers.ByteStreamToPacketParser // the byte array to packet parser
}

func CreateNewKcpLayer1(configuration *config.Configuration,
	outputChannel chan<- packets.Packet,
	inputChannel <-chan commands.Command,
	parser parsers.ByteStreamToPacketParser) *KcpAsLayer1 {
	return &KcpAsLayer1{
		outputChannel: outputChannel,
		running:       false,
		sessions:      make(map[int64]*kcp.UDPSession),
		listener:      nil,
		configuration: configuration,
		parser:        parser,
		inputChannel:  inputChannel,
	}
}

func (layer *KcpAsLayer1) Start() error {
	if layer.running {
		return nil
	}

	if layer.outputChannel == nil {
		return errors.New("channel already closed")
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
		return err
	}
}

func (layer *KcpAsLayer1) Stop() error {
	layer.running = false

	if layer.outputChannel == nil {
		return errors.New("channel already closed")
	}

	close(layer.outputChannel)

	return nil
}

func (layer *KcpAsLayer1) listenFromInputChannel() {
	// for now, still pending...
}

func (layer *KcpAsLayer1) process() {
	for layer.running {
		if session, err := layer.listener.AcceptKCP(); err == nil {
			connectionId := helpers.GenerateRandomInt64()

			for {
				if _, exists := layer.sessions[connectionId]; !exists {
					break
				}

				connectionId = helpers.GenerateRandomInt64()
			}

			layer.sessions[connectionId] = session

			go layer.handleSession(connectionId)
		}
	}
}

func (layer *KcpAsLayer1) handleSession(connectionId int64) {
	buffer := make([]byte, layer.configuration.GetBufferSize())

	for {
		if layer.outputChannel == nil {
			layer.sessions[connectionId].Close()
			delete(layer.sessions, connectionId)
			return
		}

		if session, sessionExists := layer.sessions[connectionId]; sessionExists {
			if dataLen, err := session.Read(buffer); err == nil {
				rawPacket := buffer[:dataLen]

				if packet, err := layer.parser.ParseByteArrayToPacket(&rawPacket, connectionId); err == nil {
					layer.outputChannel <- packet
				}
			} else {
				fmt.Printf("error in session %d: %s\n", connectionId, err.Error())
				layer.outputChannel <- dto.CreateInvalidPacket(connectionId)
			}
		} else {
			return
		}
	}
}
