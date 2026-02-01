package tries

import (
	"expansion-gateway/internal/structs/dictionaries"
)

type TrieHashNode struct {
	subs *dictionaries.SessionsDictionary[struct{}]
}

func createTrieHashNode() *TrieHashNode {
	return &TrieHashNode{
		subs: dictionaries.CreateNewSessionDictionary[struct{}](),
	}
}

func (node *TrieHashNode) Subscribe(identifier int64) {
	node.subs.Store(struct{}{}, identifier)
}

func (node *TrieHashNode) GetSubscribers() []int64 {
	return node.subs.Keys()
}

func (node *TrieHashNode) GetSubscribersAsMap() *dictionaries.SessionsDictionary[struct{}] {
	return node.subs
}

func (node *TrieHashNode) GetChildren() []TrieNode {
	return []TrieNode{}
}

func (node *TrieHashNode) HasChild(key SubscriptionKey) bool {
	return false
}

func (node *TrieHashNode) GetChild(key SubscriptionKey) TrieNode {
	return nil
}

func (node *TrieHashNode) GetExistChild(key SubscriptionKey) (TrieNode, bool) {
	return nil, false
}

func (node *TrieHashNode) Unsubscribe(identifier int64) {
	node.subs.Delete(identifier)
}

func (node *TrieHashNode) CreateSubscriptionChild(key SubscriptionKey) {
	// no childs here
}
