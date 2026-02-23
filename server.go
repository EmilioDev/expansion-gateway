package main

import (
	"expansion-gateway/config"
	"expansion-gateway/controllers"
	"expansion-gateway/input/kcp_handler"
	"expansion-gateway/interfaces/errorinfo"
	"expansion-gateway/interfaces/layers"
	"expansion-gateway/output/nats"
	"expansion-gateway/parsers"
	"log"
)

type GatewayServer struct {
	inputLayer  layers.Layer1         // the layer responsible of receiving the packets from the clients
	mainLayer   layers.Layer2         // the layer responsible of analyzing the packets
	outputLayer layers.Layer3         // the layer responsible of forward the packets to another services
	config      *config.Configuration // configuration object
}

func (gateway *GatewayServer) initialize() {
	gateway.config = &config.Configuration{}
	gateway.config.Initialize()

	gateway.inputLayer = gateway.getLayer1()
	gateway.mainLayer = gateway.getLayer2()
	gateway.outputLayer = gateway.getLayer3()

	if err := gateway.mainLayer.ConfigureFirstLayer(gateway.inputLayer); err != nil {
		log.Fatalln(err)
	}

	if err := gateway.mainLayer.ConfigureThirdLayer(gateway.outputLayer); err != nil {
		log.Fatalln(err)
	}
}

func (gateway *GatewayServer) Start() errorinfo.GatewayError {
	if err := gateway.mainLayer.Start(); err == nil {
		gateway.mainLayer.Wait()
	} else {
		log.Fatalln(err)
		return err
	}

	return nil
}

func (gateway *GatewayServer) Stop() errorinfo.GatewayError {
	return gateway.mainLayer.Stop()
}

func (gateway *GatewayServer) IsLayer1Working() bool {
	return gateway.inputLayer.IsWorking()
}

func (gateway *GatewayServer) IsLayer2Working() bool {
	return gateway.mainLayer.IsWorking()
}

func (gateway *GatewayServer) IsLayer3Working() bool {
	return gateway.outputLayer.IsWorking()
}

func (gateway *GatewayServer) getLayer1() layers.Layer1 {
	return kcp_handler.CreateNewKcpLayer1(
		gateway.config,
		&parsers.BasicByteArrayToPacketParser{})
}

func (gateway *GatewayServer) getLayer2() layers.Layer2 {
	if gateway.config.AreWeClusterLeaders() {
		return controllers.CreateNewLayer2Leader(gateway.config)
	}

	return controllers.CreateNewLayer2Follower(gateway.config)
}

func (gateway *GatewayServer) getLayer3() layers.Layer3 {
	return nats.GenerateNewNatsLayer3Output(gateway.config)
}

func GetGateway() *GatewayServer {
	answer := GatewayServer{}

	answer.initialize()

	return &answer
}
