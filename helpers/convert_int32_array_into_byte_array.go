package helpers

func ConvertInt32ArrayIntoByteArray(input []int32) []byte {
	result := make([]byte, len(input))

	for index, value := range input {
		result[index] = byte(value)
	}

	return result
}
