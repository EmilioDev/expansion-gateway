package enums

type LayerChannelType byte

const (
	INPUT_CHANNEL LayerChannelType = iota
	OUTPUT_CHANNEL
	BOTH_CHANNELS
)

func LayerChannelTypeToString(kind LayerChannelType) string {
	switch kind {
	case INPUT_CHANNEL:
		return "Input Channel"

	case OUTPUT_CHANNEL:
		return "Output Channel"

	case BOTH_CHANNELS:
		return "Both Channels"
	}

	return ""
}
