// file: /internals/structs/sessions_dictionary.go
package structs

import (
	"sync"
)

type SessionsDictionary[T any] struct {
	sessionsMutex sync.RWMutex
	sessions      map[int64]T
}

// adds an element to the collection
func (store *SessionsDictionary[T]) Store(data T, index int64) {
	store.sessionsMutex.Lock()
	store.sessions[index] = data
	store.sessionsMutex.Unlock()
}

// checks if an index has elements in the collection or not
func (store *SessionsDictionary[T]) Exists(index int64) bool {
	store.sessionsMutex.RLock()
	_, exists := store.sessions[index]
	store.sessionsMutex.RUnlock()

	return exists
}

// gets an element from the collection. does not return the exists
func (store *SessionsDictionary[T]) Get(index int64) T {
	store.sessionsMutex.RLock()
	result := store.sessions[index]
	store.sessionsMutex.RUnlock()

	return result
}

// gets an element from the collection and if it exists or not inside it
func (store *SessionsDictionary[T]) GetExists(index int64) (T, bool) {
	store.sessionsMutex.RLock()
	result, exists := store.sessions[index]
	store.sessionsMutex.RUnlock()

	return result, exists
}

// replaces an element in the collection. it also does the same as store
func (store *SessionsDictionary[T]) Replace(data T, index int64) {
	store.sessionsMutex.Lock()
	store.sessions[index] = data
	store.sessionsMutex.Unlock()
}

// swaps the index of two elements inside the collection
func (store *SessionsDictionary[T]) Swap(index1, index2 int64) {
	store.sessionsMutex.Lock()

	temp := store.sessions[index1]
	store.sessions[index1] = store.sessions[index2]
	store.sessions[index2] = temp

	store.sessionsMutex.Unlock()
}

// moves one element to another index, completely replacingthe value of the new index,
// and the old index gets empty
func (store *SessionsDictionary[T]) MoveTo(index1, index2 int64) {
	store.sessionsMutex.Lock()

	store.sessions[index2] = store.sessions[index1]
	delete(store.sessions, index1)

	store.sessionsMutex.Unlock()
}

// deletes the value at the selected index
func (store *SessionsDictionary[T]) Delete(index int64) {
	store.sessionsMutex.Lock()

	delete(store.sessions, index)

	store.sessionsMutex.Unlock()
}

// Iterate applies the given function to each session in the dictionary.
// The function receives the index and the data as arguments.
func (store *SessionsDictionary[T]) Iterate(fn func(index int64, data T)) {
	store.sessionsMutex.RLock()
	defer store.sessionsMutex.RUnlock()

	for index, data := range store.sessions {
		fn(index, data)
	}
}

// Clear removes all sessions from the dictionary.
func (store *SessionsDictionary[T]) Clear() {
	store.sessionsMutex.Lock()
	defer store.sessionsMutex.Unlock()

	store.sessions = make(map[int64]T)
}

// returns the number of elements in the collection
func (store *SessionsDictionary[T]) Len() int {
	store.sessionsMutex.RLock()
	defer store.sessionsMutex.RUnlock()

	return len(store.sessions)
}

// returns the keys stored in the collection
func (store *SessionsDictionary[T]) Keys() []int64 {
	store.sessionsMutex.RLock()
	defer store.sessionsMutex.RUnlock()

	answer := make([]int64, 0, store.Len())

	for index := range store.sessions {
		answer = append(answer, index)
	}

	return answer
}

// returns the values stored in the collection
func (store *SessionsDictionary[T]) Values() []T {
	store.sessionsMutex.RLock()
	defer store.sessionsMutex.RUnlock()

	answer := make([]T, 0, store.Len())

	for _, value := range store.sessions {
		answer = append(answer, value)
	}

	return answer
}

// creates a new session dictionary
func CreateNewSessionDictionary[T any]() *SessionsDictionary[T] {
	return &SessionsDictionary[T]{
		sessionsMutex: sync.RWMutex{},
		sessions:      make(map[int64]T),
	}
}
