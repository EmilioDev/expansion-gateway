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
	result := dictionaries.CreateNewSessionDictionary[struct{}]()

	result.Import(node.subs)

	return result.Keys()
}

func (node *TriePlusNode) GetSubscribersAsMap() *dictionaries.SessionsDictionary[struct{}] {
	result := dictionaries.CreateNewSessionDictionary[struct{}]()

	result.Import(node.subs)

	return result
}

func (node *TriePlusNode) GetChildren() []TrieNode {
	return node.children.Values()
}

func (node *TriePlusNode) GetChild(key SubscriptionKey) TrieNode {
	if res, exist := node.children.GetExists(key); exist {
		return res
	}

	return nil
}

func (node *TriePlusNode) HasChild(key SubscriptionKey) bool {
	return node.children.Exists(key)
}

func (node *TriePlusNode) GetExistChild(key SubscriptionKey) (TrieNode, bool) {
	return node.children.GetExists(key)
}

func (node *TriePlusNode) CreateSubscriptionChild(key SubscriptionKey) {
	if key == OneLevelWildcard && !node.children.Exists(OneLevelWildcard) {
		node.children.Store(createTriePlusNode(), key)
	} else if key == MultiLevelWildcard && !node.children.Exists(MultiLevelWildcard) {
		node.children.Store(createTrieHashNode(), key)
	} else if !node.children.Exists(key) {
		node.children.Store(createTrieCommonNode(), key)
	}
}

func (node *TriePlusNode) Subscribe(identifier int64) {
	node.subs.Store(struct{}{}, identifier)
}

func (node *TriePlusNode) Unsubscribe(identifier int64) {
	node.subs.Delete(identifier)
}
