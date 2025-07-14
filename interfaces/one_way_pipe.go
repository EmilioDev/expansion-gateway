package interfaces

type OneWayPipe interface {
	// Initialices this layer with the cannel
	InitOutputChannel(channel chan<- string)
}
