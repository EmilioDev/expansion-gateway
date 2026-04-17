package helpers

import "fmt"

func ByteArrayToString(data []byte) string {
	output := "["

	for index, value := range data {
		if index > 0 {
			output = output + fmt.Sprintf(" %d", value)
		} else {
			output = output + fmt.Sprintf("%d", value)
		}
	}

	return output + "]"
}
