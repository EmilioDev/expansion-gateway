package tries

import (
	"expansion-gateway/internal/structs/dictionaries"
)

type TriePlusNode struct {
	children *dictionaries.MutexedDictionary[SubscriptionKey, TrieNode] // direct children
	subs     *dictionaries.SessionsDictionary[struct{}]                 // direct subscription to this node
}

func createTriePlusNode() *TriePlusNode {
	return &TriePlusNode{
		children: dictionaries.CreateMutexedDictionary[SubscriptionKey, TrieNode](),
		subs:     dictionaries.CreateNewSessionDictionary[struct{}](),
	}
}

func (node *TriePlusNode) GetSubscribers() []int64 {
	return node.subs.Keys()
}

func (node *TriePlusNode) GetSubscribersAsMap() *dictionaries.SessionsDictionary[struct{}] {
	return node.subs
}

func (node *TriePlusNode) GetChildren() []TrieNode {
	return node.children.Values()
}

func (node *TriePlusNode) GetChild(key SubscriptionKey) TrieNode {
	return node.children.Get(key)
}

func (node *TriePlusNode) HasChild(key SubscriptionKey) bool {
	return node.children.Exists(key)
}

func (node *TriePlusNode) GetExistChild(key SubscriptionKey) (TrieNode, bool) {
	return node.children.GetExists(key)
}

func (node *TriePlusNode) CreateSubscriptionChild(key SubscriptionKey) {
	switch key {
	case OneLevelWildcard:
		node.children.StoreIfIndexEmpty(createTriePlusNode(), OneLevelWildcard)
	case MultiLevelWildcard:
		node.children.StoreIfIndexEmpty(createTrieHashNode(), MultiLevelWildcard)
	default:
		node.children.StoreIfIndexEmpty(createTrieCommonNode(), key)
	}
}

func (node *TriePlusNode) Subscribe(identifier int64) {
	node.subs.Store(struct{}{}, identifier)
}

func (node *TriePlusNode) Unsubscribe(identifier int64) {
	node.subs.Delete(identifier)
}
