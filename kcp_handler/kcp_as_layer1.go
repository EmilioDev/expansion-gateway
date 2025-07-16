package kcp_handler

import (
	"errors"
	"expansion-gateway/config"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/packets"
	"expansion-gateway/interfaces/parsers"
	"fmt"

	kcp "github.com/xtaci/kcp-go/v5"
)

type KcpAsLayer1 struct {
	running  bool
	sessions map[int64]*kcp.UDPSession
	listener *kcp.Listener

	// requires initialization
	outputChannel chan<- packets.Packet
	configuration *config.Configuration
	parser        parsers.ByteStreamToPacketParser
}

func CreateNewKcpLayer1(configuration *config.Configuration,
	outputChannel chan<- packets.Packet,
	parser parsers.ByteStreamToPacketParser) *KcpAsLayer1 {
	return &KcpAsLayer1{
		outputChannel: outputChannel,
		running:       false,
		sessions:      make(map[int64]*kcp.UDPSession),
		listener:      nil,
		configuration: configuration,
		parser:        parser,
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

		if _, sessionExists := layer.sessions[connectionId]; sessionExists {
			if dataLen, err := layer.sessions[connectionId].Read(buffer); err == nil {
				rawPacket := buffer[:dataLen]

				if packet, err := layer.parser.ParseByteArrayToPacket(&rawPacket, connectionId); err == nil {
					layer.outputChannel <- *packet
				}
			} else {
				fmt.Printf("error in session %d: %s\n", connectionId, err.Error())
				continue
			}
		} else {
			return
		}
	}
}
