// file: /internals/structs/mutexed_dictionary.go
package structs

import (
	"sync"
)

type MutexedDictionary[K comparable, V any] struct {
	mutex     sync.RWMutex // semaphore used for avoiding collisions
	registers map[K]V      // collection of children elements
}

// adds an element to the collection in the specified index
func (store *MutexedDictionary[T, V]) Store(data V, index T) {
	store.mutex.Lock()
	store.registers[index] = data
	store.mutex.Unlock()
}

// checks if an index has elements in the collection or not
func (store *MutexedDictionary[T, V]) Exists(index T) bool {
	store.mutex.RLock()
	_, exists := store.registers[index]
	store.mutex.RUnlock()

	return exists
}

// gets an element from the collection. does not return the exists
func (store *MutexedDictionary[T, V]) Get(index T) V {
	store.mutex.RLock()
	result := store.registers[index]
	store.mutex.RUnlock()

	return result
}

// gets an element from the collection and if it exists or not inside it
func (store *MutexedDictionary[T, V]) GetExists(index T) (V, bool) {
	store.mutex.RLock()
	result, exists := store.registers[index]
	store.mutex.RUnlock()

	return result, exists
}

// replaces an element in the collection. it also does the same as store
func (store *MutexedDictionary[T, V]) Replace(data V, index T) {
	store.mutex.Lock()
	store.registers[index] = data
	store.mutex.Unlock()
}

// swaps the index of two elements inside the collection
func (store *MutexedDictionary[T, V]) Swap(index1, index2 T) {
	store.mutex.Lock()

	temp := store.registers[index1]
	store.registers[index1] = store.registers[index2]
	store.registers[index2] = temp

	store.mutex.Unlock()
}

// moves one element to another index, completely replacingthe value of the new index,
// and the old index gets empty
func (store *MutexedDictionary[T, V]) MoveTo(index1, index2 T) {
	store.mutex.Lock()

	store.registers[index2] = store.registers[index1]
	delete(store.registers, index1)

	store.mutex.Unlock()
}

// deletes the value at the selected index
func (store *MutexedDictionary[T, V]) Delete(index T) {
	store.mutex.Lock()

	delete(store.registers, index)

	store.mutex.Unlock()
}

// if an element exists at the specified index, it will be deleted and true will be returned; but if
// there is nothing at that index, it will just return false, and nothing else happens
func (store *MutexedDictionary[T, V]) WasDeleted(index T) bool {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if _, exists := store.registers[index]; exists {
		delete(store.registers, index)
		return true
	}

	return false
}

// Iterate applies the given function to each session in the dictionary.
// The function receives the index and the data as arguments.
func (store *MutexedDictionary[T, V]) Iterate(fn func(index T, data V)) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for index, data := range store.registers {
		fn(index, data)
	}
}

// Clear removes all sessions from the dictionary.
func (store *MutexedDictionary[T, V]) Clear() {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.registers = make(map[T]V)
}

// returns the number of elements in the collection
func (store *MutexedDictionary[T, V]) Len() int {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	return len(store.registers)
}

// returns the keys stored in the collection
func (store *MutexedDictionary[T, V]) Keys() []T {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	answer := make([]T, 0, len(store.registers))

	for index := range store.registers {
		answer = append(answer, index)
	}

	return answer
}

// returns the values stored in the collection
func (store *MutexedDictionary[T, V]) Values() []V {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	answer := make([]V, 0, len(store.registers))

	for _, value := range store.registers {
		answer = append(answer, value)
	}

	return answer
}

func CreateMutexedDictionary[K comparable, V any]() *MutexedDictionary[K, V] {
	return &MutexedDictionary[K, V]{
		mutex:     sync.RWMutex{},
		registers: make(map[K]V),
	}
}
