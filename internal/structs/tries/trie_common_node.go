package tries

import (
	"expansion-gateway/internal/structs/dictionaries"
)

type TrieCommonNode struct {
	children *dictionaries.MutexedDictionary[SubscriptionKey, TrieNode] // exact children
	subs     *dictionaries.SessionsDictionary[struct{}]                 // subscriptions to this exact node
}

// creates a new node
func createTrieCommonNode() *TrieCommonNode {
	children := dictionaries.CreateMutexedDictionary[SubscriptionKey, TrieNode]()

	return &TrieCommonNode{
		children: children,
		subs:     dictionaries.CreateNewSessionDictionary[struct{}](),
	}
}

func (node *TrieCommonNode) GetSubscribers() []int64 {
	result := dictionaries.CreateNewSessionDictionary[struct{}]()

	result.Import(node.subs)

	return result.Keys()
}

func (node *TrieCommonNode) GetSubscribersAsMap() *dictionaries.SessionsDictionary[struct{}] {
	result := dictionaries.CreateNewSessionDictionary[struct{}]()

	result.Import(node.subs)

	return result
}

func (node *TrieCommonNode) GetChildren() []TrieNode {
	return node.children.Values()
}

func (node *TrieCommonNode) GetChild(key SubscriptionKey) TrieNode {
	if res, exist := node.children.GetExists(key); exist {
		return res
	}

	return nil
}

func (node *TrieCommonNode) HasChild(key SubscriptionKey) bool {
	return node.children.Exists(key)
}

func (node *TrieCommonNode) GetExistChild(key SubscriptionKey) (TrieNode, bool) {
	return node.children.GetExists(key)
}

func (node *TrieCommonNode) CreateSubscriptionChild(key SubscriptionKey) {
	if key == OneLevelWildcard && !node.children.Exists(OneLevelWildcard) {
		node.children.Store(createTriePlusNode(), key)
	} else if key == MultiLevelWildcard && !node.children.Exists(MultiLevelWildcard) {
		node.children.Store(createTrieHashNode(), key)
	} else if !node.children.Exists(key) {
		node.children.Store(createTrieCommonNode(), key)
	}
}

func (node *TrieCommonNode) Subscribe(identifier int64) {
	node.subs.Store(struct{}{}, identifier)
}

func (node *TrieCommonNode) Unsubscribe(identifier int64) {
	node.subs.Delete(identifier)
}
