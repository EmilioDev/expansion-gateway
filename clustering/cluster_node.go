package clustering

import (
	"expansion-gateway/config"
	"expansion-gateway/interfaces/errorinfo"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
)

type ClusterNode struct {
	grpcCurrentServerPath   string          // server path
	temporalMessagesCounter atomic.Int64    // the temporal messages counter
	server                  *grpc.Server    // the grpc server
	MessagesCounter         *atomic.Int64   // final messages counter
	grpcPort                uint16          // the port the server will use
	isWorking               *atomic.Bool    // wether is this component working or not
	startOnce               *sync.Once      // starts the service only once
	stopOnce                *sync.Once      // stops the service only once
	sessionsTimeout         int64           // the timespan in seconds used to check the client sessions, or the timeout this session has before being declared unhealthy
	epoch                   *atomic.Int64   // the current epoch of the server
	startTime               time.Time       // the time this server was started
	timeMutex               sync.RWMutex    // the mutex used in the timing operations
	wg                      *sync.WaitGroup // the wait group used to indicate when this server has finished it's task
	initCallback            func()          // the function that should be called on the start method, after everything is running
	stopCallback            func()          // the function to be called when the server is going to stop
}

// creates a base cluster node
func CreateBaseClusterNode(
	conf *config.Configuration,
	wg *sync.WaitGroup,
	initCallback func(),
	stopCallback func()) ClusterNode {
	return ClusterNode{
		grpcCurrentServerPath:   conf.GetGrpcCurrentServerPath(),
		grpcPort:                conf.GetClusterGrpcPort(),
		temporalMessagesCounter: atomic.Int64{},
		isWorking:               &atomic.Bool{},
		startOnce:               &sync.Once{},
		stopOnce:                &sync.Once{},
		server:                  grpc.NewServer(),
		MessagesCounter:         &atomic.Int64{},
		sessionsTimeout:         120,
		epoch:                   &atomic.Int64{},
		startTime:               time.Now(),
		timeMutex:               sync.RWMutex{},
		wg:                      wg,
		initCallback:            initCallback,
		stopCallback:            stopCallback,
	}
}

// increases the temporal messages counter by one
func (cluster *ClusterNode) NewMessage() {
	cluster.temporalMessagesCounter.Add(1)
}

// copies the current value of the temporal messages counter to the final one and then resets it
func (cluster *ClusterNode) ResetMessageCounter() {
	cluster.MessagesCounter.Store(cluster.temporalMessagesCounter.Load())
	cluster.temporalMessagesCounter.Store(0)
}

// starts the grpc server of this cluster member
func (cluster *ClusterNode) Start() errorinfo.GatewayError {
	cluster.startOnce.Do(func() {
		go cluster.runServer()
		cluster.wg.Add(1)
	})

	return nil
}

// stops the grpc server of this cluster member
func (cluster *ClusterNode) Stop() errorinfo.GatewayError {
	if !cluster.isWorking.Load() {
		return nil
	}

	cluster.stopOnce.Do(func() {
		cluster.isWorking.Store(false)
		cluster.server.GracefulStop()

		if cluster.stopCallback != nil {
			cluster.stopCallback()
		}

		cluster.wg.Done()
		fmt.Println("gateway cluster member is down now!")
	})

	return nil
}

// gets the current epoch of the server
func (cluster *ClusterNode) GetEpoch() int64 {
	return cluster.epoch.Load()
}

// moves the server to the next epoch
func (cluster *ClusterNode) NextEpoch() {
	cluster.epoch.Add(1)
}

// returns the nanoseconds elapsed since the server started
func (cluster *ClusterNode) GetElapsedTime() int64 {
	cluster.timeMutex.RLock()
	defer cluster.timeMutex.RUnlock()

	return time.Since(cluster.startTime).Nanoseconds()
}

// changes the start time of the server
func (cluster *ClusterNode) SetStartTime(time *time.Time) {
	if time != nil {
		cluster.timeMutex.Lock()
		cluster.startTime = *time
		cluster.timeMutex.Unlock()
	}
}

// ==== privates ====

// runs the server
func (cluster *ClusterNode) runServer() {
	address, err := net.Listen("tcp", fmt.Sprintf(":%d", cluster.grpcPort))

	if err != nil {
		log.Fatalf("failed to create path for cluster leader grpc server. port: %d, error: %v", cluster.grpcPort, err)
	}

	cluster.isWorking.Store(true) // yep, it's running, let's hope

	if cluster.initCallback != nil {
		go cluster.initCallback()
	}

	if err := cluster.server.Serve(address); err != nil {
		cluster.isWorking.Store(false) // this line is not needed, the next one will close the app, but...
		log.Fatalf("failed to start cluster leader grpc server. port: %d, error: %v", cluster.grpcPort, err)
	}
}
