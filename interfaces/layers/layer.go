package layers

// Base type of all the three layers of the gateway
type Layer interface {
	// Sends data to be processed by this layer
	Start() error
	Stop() error
}
