// file: /internals/structs/subscription_key.go
package structs

import (
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

	for x := 0; x < size; x++ {
		b := s[x]

		if b == oneLevelWildcard {
			return false
		}
	}

	return true
}

// returns all key levels. ex: core/+/idea -> ["core", "+", "idea"]
func (s SubscriptionKey) GetKeyLevels() []string {
	return strings.Split(string(s), "/")
}
