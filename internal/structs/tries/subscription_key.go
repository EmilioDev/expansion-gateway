// file: /internals/structs/subscription_key.go
package tries

import (
	subsErrors "expansion-gateway/errors/subscriptions"
	"expansion-gateway/helpers"
	"expansion-gateway/interfaces/errorinfo"
	"strings"
)

type SubscriptionKey string

// true if this subscription key has no wildcards
func (s SubscriptionKey) IsFixedKey() bool {
	size := len(s)
	const oneLevelWildcard byte = byte('+')
	const multiLevelWildcard byte = byte('#')

	if size == 0 {
		return true
	}

	if s[size-1] == multiLevelWildcard {
		return false
	}

	return !helpers.IsCharacterInString(string(s), oneLevelWildcard)
}

// checks if this key completely is a wildcard. Is not the same as IsFixedKey that checks if contains
// wildcards among other things
func (s SubscriptionKey) IsWildCard() bool {
	key := string(s)

	return key == "#" || key == "+"
}

// checks if this key is exactly equal to +
func (s SubscriptionKey) IsPlus() bool {
	key := string(s)

	return key == "+"
}

// checks if this key is exactly equal to #
func (s SubscriptionKey) IsHash() bool {
	key := string(s)

	return key == "#"
}

// returns all key levels. ex: core/+/idea -> ["core", "+", "idea"]
func (s SubscriptionKey) GetKeyLevels() []SubscriptionKey {
	keysAsString := strings.Split(string(s), "/")
	answer := make([]SubscriptionKey, len(keysAsString))

	for index, key := range keysAsString {
		answer[index] = SubscriptionKey(key)
	}

	return answer
}

// returns the number of levels this subscription key has
func (s SubscriptionKey) GetCantOfKeyLevels() int {
	return len(strings.Split(string(s), "/"))
}

// converts this key into a simple string
func (s SubscriptionKey) ToString() string {
	return string(s)
}

// returns the number of bytes in this subscription key
func (s SubscriptionKey) KeyLength() int {
	return len([]byte(s))
}

// converts this subscription into a byte array
func (s SubscriptionKey) ToByteArray() []byte {
	return []byte(s)
}

func (s SubscriptionKey) ToNATSkey() string {
	return strings.ReplaceAll(string(s), "/", ".")
}

// takes a string, and if it can be used as subscription, it creates one from it, or if it is invalid as
// subscription, then returns an error
func ConvertStringToSubscriptionKey(subscription string) (SubscriptionKey, errorinfo.GatewayError) {
	const filePath string = "/internals/structs/subscription_key.go"

	if helpers.IsAnyOfTheseCharactersInString(subscription, []byte{'@', '.'}) {
		return SubscriptionKey(""), subsErrors.CreateUseOfInvalidCharactersInSubscriptionError(filePath, subscription, 93)
	}

	levels := strings.Split(subscription, "/")
	lastLevel := len(levels) - 1

	const oneLevelWildcard byte = byte('+')
	const multiLevelWildcard byte = byte('#')

	for index, level := range levels {
		if len(level) > 1 && helpers.IsAnyOfTwoCharactersInString(level, oneLevelWildcard, multiLevelWildcard) {
			return SubscriptionKey(""), subsErrors.CreateInvalidUseOfWildcardsError(filePath, subscription, 104)
		} else if index != lastLevel && level == "#" {
			return SubscriptionKey(""), subsErrors.CreateInvalidUseOfWildcardsError(filePath, subscription, 106)
		}
	}

	return SubscriptionKey(subscription), nil
}
