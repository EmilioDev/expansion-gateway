package helpers

func ConvertByteArrayIntoInt32Array(input []byte) []int32 {
	inputSize := len(input)

	result := make([]int32, inputSize)

	for index, value := range input {
		result[index] = int32(value)
	}

	return result
}
