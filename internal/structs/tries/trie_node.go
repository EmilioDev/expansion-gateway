// file: /internals/structs/tries/trie_node.go
package tries

import (
	"expansion-gateway/internal/structs/dictionaries"
)

type TrieNode interface {
	GetSubscribers() []int64                                         // gets all the subscribers of this node
	GetSubscribersAsMap() *dictionaries.SessionsDictionary[struct{}] // gets all the subscribers of this node as a map
	GetChildren() []TrieNode                                         // gets all the childrens of this node
	Subscribe(identifier int64)                                      // subscribes a user to this member
	Unsubscribe(identifier int64)                                    // unsubscribe a member from this subscription
	CreateSubscriptionChild(key SubscriptionKey)                     // adds a child to this node
	GetChild(key SubscriptionKey) TrieNode                           // gets the child (or nil) asosciated to the key
	HasChild(key SubscriptionKey) bool                               // returns true if there is a child asosciated to that key
	GetExistChild(key SubscriptionKey) (TrieNode, bool)              // returns the child and true if there is a child node asosciated to that key, or nil and false if none is asosciated
}
