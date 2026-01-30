// file: /internals/structs/trie.go
package structs

type Trie struct {
	patternSubscriptions *TrieNode
	exactSubscriptions   *MutexedDictionary[SubscriptionKey, SessionsDictionary[struct{}]]
}

func CreateTrie() *Trie {
	return &Trie{
		patternSubscriptions: createTrieNode(),
		exactSubscriptions:   CreateMutexedDictionary[SubscriptionKey, SessionsDictionary[struct{}]](),
	}
}

func (trie *Trie) GetSubscribers(subscription SubscriptionKey) []int64 {
	return []int64{}
}

func (trie *Trie) SubscribeTo(subscription SubscriptionKey, subscriber int64) {
	//
}
