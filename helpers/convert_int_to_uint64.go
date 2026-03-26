package helpers

func ConvertIntToUint64(source int) uint64 {
	if source < 0 {
		return 0
	}

	return uint64(source)
}
