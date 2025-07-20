package helpers

// Decodes a variable byte integer into a single int, and also in the second int returns
// the number of bytes read
func DecodeVariableByteInteger(data []byte) (int, int) {
	variableByteInteger := 0
	multiplier := 1
	bytesRead := 0

	for i, b := range data {
		variableByteInteger += int(b&127) * multiplier
		multiplier *= 128
		bytesRead++

		if b&128 == 0 {
			break
		}

		if i >= 3 { // Maximum of 4 bytes allowed
			return 0, 0
		}
	}

	return variableByteInteger, bytesRead
}
