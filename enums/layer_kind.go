package enums

type LayerKind byte

const (
	LAYER_1 LayerKind = iota
	LAYER_2
	LAYER_3
)

func LayerKindToString(kind LayerKind) string {
	switch kind {
	case LAYER_1:
		return "Layer 1"

	case LAYER_2:
		return "Layer 2"

	case LAYER_3:
		return "Layer 3"

	default:
		return ""
	}
}
