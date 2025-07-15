package kcp_handler

import (
	"errors"
	"expansion-gateway/config"
	"expansion-gateway/dto"
	"expansion-gateway/helpers"
	"expansion-gateway/parsers"
	"fmt"

	kcp "github.com/xtaci/kcp-go/v5"
)

type KcpAsLayer1 struct {
	conf          *config.Configuration
	outputChannel chan<- dto.Packet
	running       bool
	sessions      map[int64]*kcp.UDPSession
	listener      *kcp.Listener
}

func (layer *KcpAsLayer1) Start() error {
	if layer.running {
		return nil
	}

	if layer.outputChannel == nil {
		return errors.New("channel already closed")
	}

	var serverPath string = layer.conf.GetServerAddress()

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
	buffer := make([]byte, layer.conf.GetBufferSize())

	for {
		if layer.outputChannel == nil {
			layer.sessions[connectionId].Close()
			delete(layer.sessions, connectionId)
			return
		}

		if _, sessionExists := layer.sessions[connectionId]; sessionExists {
			if dataLen, err := layer.sessions[connectionId].Read(buffer); err == nil {
				rawPacket := buffer[:dataLen]

				if packet, err := parsers.ParseByteArrayToPacket(&rawPacket, connectionId); err == nil {
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
