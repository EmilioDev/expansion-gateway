package layers

type SmartLayer interface {
	Layer
	// method for configuring layer 1
	ConfigureFirstLayer(layer DumbLayer)

	// ...later we add for layer 2 and 3
}
