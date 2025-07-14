package helpers

import "math"

func ConvertIntToUint16(source int) uint16 {
	if source < 0 {
		return 0
	}

	if source >= math.MaxUint16 {
		return math.MaxUint16
	}

	return uint16(source)
}
