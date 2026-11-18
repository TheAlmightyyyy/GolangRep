package storage

import (
	// "financialExchange/enum"
	// "github.com/google/uuid"
	"sync"
)

type Store[K comparable, V any] struct {
	mutex sync.RWMutex
	items map[K]V
}

func NewStore[K comparable, V any](capaity int) *Store[K, V] {
	return &Store[K, V]{
		items: make(map[K]V, capaity),
	}
}

func (store *Store[K, V]) Put(key K, value V) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.items[key] = value
}

func (store *Store[K, V]) Get(key K) (V, bool) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	value, ok := store.items[key]
	return value, ok
}

func (store *Store[K, V]) Delete(key K) bool {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if _, ok := store.items[key]; !ok {
		return false
	}
	delete(store.items, key)
	return true
}

func (store *Store[K, V]) Update(key K, modif func(V) V) (V, bool) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	value, ok := store.items[key]
	if !ok {
		var zero V
		return zero, false
	}

	UpdatedValue := modif(value)
	return UpdatedValue, true
}

type IndexedStore[K comparable, V any] struct {
	mutex sync.RWMutex
	order []K
	byKey map[K]int
	data  map[K]V
}

func (indexStore *IndexedStore[K, V]) Get(key K) (V, bool) {
	indexStore.mutex.RLock()
	defer indexStore.mutex.RUnlock()

	value, ok := indexStore.data[key]
	return value, ok
}

func NewIndexedStore[K comparable, V any](capacity int) *IndexedStore[K, V] {
	return &IndexedStore[K, V]{
		order: make([]K, 0, capacity),
		byKey: make(map[K]int, capacity),
		data:  make(map[K]V, capacity),
	}
}

func (storage *IndexedStore[K, V]) Put(key K, newValue V) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	if _, ok := storage.byKey[key]; ok {
		storage.data[key] = newValue
		return
	}
	storage.byKey[key] = len(storage.order)
	storage.order = append(storage.order, key)
	storage.data[key] = newValue
	// //если ключа нет, то сохраняем текущую длину слайса как индекс для нового ключа
	// storage.byKey[key] = len(storage.order)
	// //добавляем ключ в конец среза, чтобы сохранить порядок добавления
	// storage.order = append(storage.order, key)
	// //сохраняем значение в хранилище
	// storage.data[key] = newValue
}

func (storage *IndexedStore[K, V]) Delete(key K) bool {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	indexInSlice, ok := storage.byKey[key]

	if !ok {
		return false
	}

	lastIndex := len(storage.order) - 1
	if indexInSlice != lastIndex {
		lastValue := storage.order[lastIndex]
		storage.order[indexInSlice] = lastValue
		storage.byKey[lastValue] = indexInSlice
	}

	storage.order = storage.order[:lastIndex]

	delete(storage.byKey, key)
	delete(storage.data, key)

	return true
}

func (indexStore *IndexedStore[K, V]) Ordered() []V {
	indexStore.mutex.RLock()
	defer indexStore.mutex.RUnlock()

	out := make([]V, len(indexStore.order))
	for index, key := range indexStore.order {
		out[index] = indexStore.data[key]
	}
	return out
}

type Storage[T comparable] struct {
	storage map[T]struct{}
}

func NewStorage[T comparable](capacity int) *Storage[T] {
	return &Storage[T]{
		storage: make(map[T]struct{}, capacity),
	}
}

func (storage *Storage[T]) Add(value T) {
	storage.storage[value] = struct{}{}
}

func (storage *Storage[T]) Remove(value T) {
	delete(storage.storage, value)
}

func (storage *Storage[T]) Has(value T) bool {
	_, ok := storage.storage[value]
	return ok
}

func (storage *Storage[T]) Lenght() int {
	return len(storage.storage)
}

func (storage *Storage[T]) Items() []T {
	outSlice := make([]T, 0, len(storage.storage))
	for value := range storage.storage {
		outSlice = append(outSlice, value)
	}
	return outSlice
}
