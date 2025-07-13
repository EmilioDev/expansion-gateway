package layers

type Layer3 interface {
	Layer
	SetLayer2(layer2 *Layer2)
}
