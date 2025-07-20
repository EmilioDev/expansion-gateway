package helpers

func EncodeAsVariableByteInteger(number int) []byte {
	var result []byte = []byte{}

	for number > 0 {
		encodedByte := number % 128
		number /= 128

		if number > 0 {
			encodedByte = encodedByte | 0x80
		}

		result = append(result, byte(encodedByte))
	}

	return result
}
