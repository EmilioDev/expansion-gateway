package interfaces

type Initializable interface {
	// Initialices this layer with the cannel
	Init(channel chan<- string)
}
