package interfaces

type OneWayPipe interface {
	// Initialices this layer with the cannel
	SendToNextLayer(channel chan<- string)
}
