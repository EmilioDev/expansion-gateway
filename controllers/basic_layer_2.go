package controllers

import (
	"expansion-gateway/config"
	"expansion-gateway/enums"
	"expansion-gateway/errors/layererrors"
	disp "expansion-gateway/interfaces/dispatchers"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/layers"
	"expansion-gateway/interfaces/packets"
	"expansion-gateway/internal/others"
	"sync/atomic"
	"time"
)

type BasicLayer2 struct {
	layer1        layers.Layer1
	layer3        layers.Layer3
	working       *atomic.Bool
	configuration *config.Configuration
	layer1Reciver disp.Reciver
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
			47,
			enums.LAYER_2,
			layer.layer1,
			layer.layer3)
	}

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

	layer.layer1.Stop()
	layer.layer3.Stop()

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

			layer.handlePacket(packet)

		default:
			time.Sleep(time.Millisecond * 10) // Yield CPU, prevent tight loop
		}
	}
}

func (layer *BasicLayer2) handlePacket(packet packets.Packet) errorinfo.GatewayError {
	// empty for now
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
	}
}
