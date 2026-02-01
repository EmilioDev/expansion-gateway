// file: /internals/structs/trie.go
package tries

import (
	"expansion-gateway/internal/structs/dictionaries"
)

type Trie struct {
	root TrieNode
}

// creates a new subscription trie
func CreateTrie() *Trie {
	return &Trie{
		root: createTrieCommonNode(),
	}
}

// gets all the subscribers asosciated to a subscription
func (trie *Trie) GetSubscribers(subscription SubscriptionKey) []int64 {
	levels := subscription.GetKeyLevels()
	return getAllSubscriptions(levels, 0, trie.root).Keys()
}

// joins a client to a subscription
func (trie *Trie) SubscribeTo(subscription SubscriptionKey, subscriber int64) {
	levels := subscription.GetKeyLevels()
	destiny := trie.root

	for _, key := range levels {
		if currentLevel, exists := destiny.GetExistChild(key); exists {
			destiny = currentLevel
		} else {
			destiny.CreateSubscriptionChild(key)
			destiny = destiny.GetChild(key)
		}
	}

	destiny.Subscribe(subscriber)
}

// removes a client from a subscription
func (trie *Trie) UnsubscribeTo(subscription SubscriptionKey, subscriber int64) {
	levels := subscription.GetKeyLevels()
	destiny := trie.root

	for _, key := range levels {
		if currentLevel, exists := destiny.GetExistChild(key); exists {
			destiny = currentLevel
		} else {
			return
		}
	}

	destiny.Unsubscribe(subscriber)
}

func getAllSubscriptions(subscriptionLevels []SubscriptionKey, currentIndex int, root TrieNode) *dictionaries.SessionsDictionary[struct{}] {
	// if this is the final level of the subscription, then finish here
	if currentIndex == len(subscriptionLevels) {
		return root.GetSubscribersAsMap()
	}

	result := dictionaries.CreateNewSessionDictionary[struct{}]()

	// we take all the hash subscriptions
	if hashSubs, exists := root.GetExistChild(MultiLevelWildcard); exists {
		result.Import(hashSubs.GetSubscribersAsMap())
	}

	// exact subscriptions
	if subs, exists := root.GetExistChild(subscriptionLevels[currentIndex]); exists {
		result.Import(getAllSubscriptions(subscriptionLevels, currentIndex+1, subs))
	}

	// plus subscriptions
	if plusSubs, exists := root.GetExistChild(OneLevelWildcard); exists {
		result.Import(getAllSubscriptions(subscriptionLevels, currentIndex+1, plusSubs))
	}

	return result
}
