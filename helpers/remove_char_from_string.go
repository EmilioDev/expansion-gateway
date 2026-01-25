package helpers

func RemoveCharFromString(source string, char byte) string {
	b := make([]byte, 0, len(source))

	for i := 0; i < len(source); i++ {
		if source[i] != char {
			b = append(b, source[i])
		}
	}

	return string(b)
}
