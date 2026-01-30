package helpers

// checks if a string contains a specific character
func IsCharacterInString(stringToCheck string, character byte) bool {
	size := len(stringToCheck)

	for x := 0; x < size; x++ {
		currentCharacter := stringToCheck[x]

		if currentCharacter == character {
			return true
		}
	}

	return false
}

// checks if a string contains any of two specific characters
func IsAnyOfTwoCharactersInString(stringToCheck string, character1, character2 byte) bool {
	size := len(stringToCheck)

	for x := 0; x < size; x++ {
		currentCharacter := stringToCheck[x]

		if currentCharacter == character1 || currentCharacter == character2 {
			return true
		}
	}

	return false
}
