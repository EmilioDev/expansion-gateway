package nats

import (
	"expansion-gateway/config"
	natsDto "expansion-gateway/dto/nats"
	"expansion-gateway/enums"
	"expansion-gateway/errors/layererrors"
	natsErrors "expansion-gateway/errors/nats"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/dispatchers"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/packets"
	"expansion-gateway/internal/structs/tries"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"

	"os"
	"sync"
	"sync/atomic"
	"time"
)

type NatsBasicHandlerLayer3 struct {
	layer2Dispatcher  dispatchers.Dispatcher[packets.OutputPacket] // packet dispatcher to layer 2
	layer2Receiver    dispatchers.Reciver[packets.OutputPacket]    // packet receiver from layer 2
	natsServerPath    string                                       // path to nats
	connection        *nats.Conn                                   // connection to nats
	working           *atomic.Bool                                 // tells you if this layer is working or not
	shutdownOnce      *sync.Once                                   // executes the shutdown only once
	startOnce         *sync.Once                                   // executes the start only once
	wg                *sync.WaitGroup                              // the wait group of layer 3
	messageCounter    atomic.Uint64                                // used in the identifier of each message
	gatewayNameInNats string                                       // this is going to be used in the identifier of each message too
	ecoPath           tries.SubscriptionKey                        // the path that should be answered back with the same payload inmediatly, without forwarding
}

func (layer *NatsBasicHandlerLayer3) Publish(data packets.OutputPacket) errorinfo.GatewayError {
	return nil
}

func (layer *NatsBasicHandlerLayer3) SubscribeTo(topic string) errorinfo.GatewayError {
	return nil
}

func (layer *NatsBasicHandlerLayer3) UnsubscribeTo(topic string) errorinfo.GatewayError {
	return nil
}

func (layer *NatsBasicHandlerLayer3) ConfigureDumbLayer(outputChannel dispatchers.Dispatcher[packets.OutputPacket], inputChannel dispatchers.Reciver[packets.OutputPacket]) errorinfo.GatewayError {
	layer.layer2Dispatcher = outputChannel
	layer.layer2Receiver = inputChannel

	return nil
}

func (layer *NatsBasicHandlerLayer3) CloseSession(sessionId int64) errorinfo.GatewayError {
	return nil
}

func (layer *NatsBasicHandlerLayer3) MoveClientTo(origin, destiny int64) {
	//
}

func (layer *NatsBasicHandlerLayer3) Start() errorinfo.GatewayError {
	const filePath string = "/output/nats/nats_basic_layer_3_handler.go"

	if layer.IsWorking() {
		return nil
	}

	if layer.layer2Dispatcher == nil {
		return layererrors.CreateChannelClosed_LayerError(filePath, 66, enums.LAYER_3, enums.OUTPUT_CHANNEL)
	}

	var result errorinfo.GatewayError = nil

	layer.startOnce.Do(func() {
		if err := layer.connect(); err != nil {
			result = err
		} else {
			layer.working.Store(true)

			go layer.process()
			layer.initializeLayer2Listeners()
		}
	})

	return result
}

func (layer *NatsBasicHandlerLayer3) Stop() errorinfo.GatewayError {
	if !layer.IsWorking() {
		return nil
	}

	layer.shutdownOnce.Do(func() {
		layer.working.Store(false)
		layer.connection.Close()
		layer.wg.Wait()
	})

	return nil
}

func (layer *NatsBasicHandlerLayer3) IsWorking() bool {
	return layer.working.Load()
}

func (layer *NatsBasicHandlerLayer3) RenameGateway(newName string) {
	layer.gatewayNameInNats = newName
}

// ===== privates =====

// connect
func (layer *NatsBasicHandlerLayer3) connect() errorinfo.GatewayError {
	const certificatePath string = "./certificates/expansion-user.creds"
	const filePath string = "/output/nats/nats_basic_layer_3_handler.go"

	if _, errNotExist := os.Stat(certificatePath); errNotExist == nil {
		if nc, err := nats.Connect(layer.natsServerPath, nats.UserCredentials(certificatePath)); err == nil {
			layer.connection = nc
		} else {
			return natsErrors.CreateConnectionToNatsFailedError(filePath, 73, err)
		}
	} else {
		if nc, err := nats.Connect(layer.natsServerPath); err == nil {
			layer.connection = nc
		} else {
			return natsErrors.CreateConnectionToNatsFailedError(filePath, 79, err)
		}
	}

	return nil
}

// process packets from NATS
func (layer *NatsBasicHandlerLayer3) process() {
	defer layer.wg.Done()

	// *: partial match (one level)
	// >: full match (all levels from the wildcard onward)
	layer.connection.Subscribe("output.v1.@.*", func(message *nats.Msg) {
		splittedSubject := strings.SplitN(message.Subject, "@", 2)

		if len(splittedSubject) == 2 {
			if key, err := tries.ConvertStringToSubscriptionKey(splittedSubject[1]); err == nil {
				packet := natsDto.CreateNewNatsDataBasicTransferRecipe(key, message.Data)
				layer.layer2Dispatcher.Dispatch(packet)
			}
		}
	})
}

// initialize layer 2 listeners
func (layer *NatsBasicHandlerLayer3) initializeLayer2Listeners() {
	shards := layer.layer2Receiver.ShardCount()

	for x := 0; x < shards; x++ {
		layer.wg.Add(1)
		go layer.handleShardFromLayer2(x)
	}
}

// handle shard from layer 2
func (layer *NatsBasicHandlerLayer3) handleShardFromLayer2(shardIndex int) {
	defer layer.wg.Done()
	channel := layer.layer2Receiver.GetShard(shardIndex)

	for layer.IsWorking() {
		select {
		case packet, ok := <-channel:
			if !ok {
				return
			}

			if packet == nil {
				continue
			}

			// if this is a eco
			if packet.GetKey() == layer.ecoPath {
				layer.layer2Dispatcher.Dispatch(natsDto.CreateNewNatsDataBasicTransferRecipe(packet.GetKey(), packet.GetPayload()))
				continue
			}

			identifier := helpers.ConvertInt64Into8Bytes(packet.GetSender())
			payload := append([]byte{}, identifier[:]...)
			payload = append(payload, packet.GetPayload()...)

			layer.connection.Publish(
				fmt.Sprintf("input.v1.@.%s", packet.GetKey().ToString()),
				payload,
			)

		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func GenerateNewNatsLayer3Output(configuration *config.Configuration) *NatsBasicHandlerLayer3 {
	return &NatsBasicHandlerLayer3{
		natsServerPath:    configuration.GetNATSserverPath(),
		layer2Dispatcher:  nil,
		layer2Receiver:    nil,
		connection:        nil,
		working:           &atomic.Bool{},
		shutdownOnce:      &sync.Once{},
		startOnce:         &sync.Once{},
		wg:                &sync.WaitGroup{},
		messageCounter:    atomic.Uint64{},
		gatewayNameInNats: "gateway",
		ecoPath:           configuration.GetEcoPath(),
	}
}
