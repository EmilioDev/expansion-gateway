package helpers

import "math"

func ConvertUint64IntoUint32(input uint64) uint32 {
	if input >= math.MaxUint32 {
		return math.MaxUint32
	}

	return uint32(input)
}
