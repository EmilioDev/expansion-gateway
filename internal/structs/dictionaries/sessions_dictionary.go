// file: /internals/structs/dictionaries/sessions_dictionary.go
package dictionaries

import (
	"crypto/rand"
	"encoding/binary"
	"math"
)

type SessionsDictionary[T any] struct {
	*MutexedDictionary[int64, T]
}

// stores a value and returns the index where it was stored
func (store *SessionsDictionary[T]) Add(data T) int64 {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	var buf [8]byte
	var index int64 = 0

	for {
		if _, err := rand.Read(buf[:]); err == nil {
			raw := binary.LittleEndian.Uint64(buf[:])
			index = int64(raw)

			if index == 0 {
				continue
			}

			if _, exists := store.registers[index]; !exists {
				break
			}
		} else {
			var x int64 = math.MinInt64
			var limit int64 = math.MaxInt64

			for ; x <= limit; x++ {
				if x == 0 {
					continue
				}

				if _, exists := store.registers[x]; !exists {
					index = x
					break
				}
			}

			break
		}
	}

	store.registers[index] = data

	return index
}

// import all registers from another collection
func (store *SessionsDictionary[T]) Import(source *SessionsDictionary[T]) {
	store.MutexedDictionary.Import(source.MutexedDictionary)
}

// creates a new session dictionary
func CreateNewSessionDictionary[T any]() *SessionsDictionary[T] {
	return &SessionsDictionary[T]{
		MutexedDictionary: CreateMutexedDictionary[int64, T](),
	}
}
