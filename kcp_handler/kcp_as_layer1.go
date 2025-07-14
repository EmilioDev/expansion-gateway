package kcp_handler

import (
	"errors"
	"expansion-gateway/config"
	"fmt"

	kcp "github.com/xtaci/kcp-go/v5"
)

type KcpAsLayer1 struct {
	conf          *config.Configuration
	outputChannel chan<- string
	running       bool
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
		layer.process(listener)
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

func (layer *KcpAsLayer1) process(listener *kcp.Listener) {
	for layer.running {
		if session, err := listener.AcceptKCP(); err == nil {
			go handleSession(session, layer.outputChannel)
		}
	}
}

func handleSession(session *kcp.UDPSession, outputChannel chan<- string) {
	const bufferSize int = 4096
	buffer := make([]byte, bufferSize)

	for {
		if outputChannel == nil {
			break
		}

		if dataLen, err := session.Read(buffer); err == nil {
			//packet := buffer[:dataLen]

			fmt.Println("packet length is", dataLen)
		}
	}
}
