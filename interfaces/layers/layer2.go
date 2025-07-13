package layers

type Layer2 interface {
	Layer
	SetLayer3(layer3 *Layer3)
	SetLayer1(layer1 *Layer1)
}
