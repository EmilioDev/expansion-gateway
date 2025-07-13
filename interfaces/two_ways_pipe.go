package interfaces

type TwoWaysPipe interface {
	SetChannelPreviousLayer(channel chan<- string)
	SetChannelToNextLayer(channel chan<- string)
}
