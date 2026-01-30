// file: /internals/structs/subscription_key.go
package structs

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

// returns all key levels. ex: core/+/idea -> ["core", "+", "idea"]
func (s SubscriptionKey) GetKeyLevels() []string {
	return strings.Split(string(s), "/")
}

// takes a string, and if it can be used as subscription, it creates one from it, or if it is invalid as
// subscription, then returns an error
func ConvertStringToSubscriptionKey(subscription string) (SubscriptionKey, errorinfo.GatewayError) {
	levels := strings.Split(subscription, "/")
	lastLevel := len(levels) - 1

	const oneLevelWildcard byte = byte('+')
	const multiLevelWildcard byte = byte('#')
	const filePath string = "/internals/structs/subscription_key.go"

	for index, level := range levels {
		if len(level) > 1 && helpers.IsAnyOfTwoCharactersInString(level, oneLevelWildcard, multiLevelWildcard) {
			return SubscriptionKey(""), subsErrors.CreateInvalidUseOfWildcardsError(filePath, subscription, 44)
		} else if index != lastLevel && level == "#" {
			return SubscriptionKey(""), subsErrors.CreateInvalidUseOfWildcardsError(filePath, subscription, 46)
		}
	}

	return SubscriptionKey(subscription), nil
}
