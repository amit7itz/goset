package store

import (
	"encoding/json"
	"errors"
	"fmt"
)

type SetStore[T comparable] interface {
	Add(items ...T)
	Remove(item T) error
	Discard(items ...T)
	Len() int
	IsEmpty() bool
	Contains(item T) bool
	Pop() (T, error)
	Items() []T
	For(func(item T))
	ForWithBreak(func(item T) bool)
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(b []byte) error
}

type SimpleSetStore[T comparable] struct {
	store map[T]struct{}
}

func NewSimpleStore[T comparable]() *SimpleSetStore[T] {
	return &SimpleSetStore[T]{
		store: make(map[T]struct{}),
	}
}

// Add adds item(s) to the store
func (s *SimpleSetStore[T]) Add(items ...T) {
	for _, item := range items {
		s.store[item] = struct{}{}
	}
}

// Remove removes a single item from the store. Returns error if the item is not in the Set
// See also: Discard()
func (s *SimpleSetStore[T]) Remove(item T) error {
	if s.Contains(item) {
		delete(s.store, item)
		return nil
	}
	return fmt.Errorf("item not found: %v ", item)
}

// Discard removes item(s) from the store if exist
// See also: Remove()
func (s *SimpleSetStore[T]) Discard(items ...T) {
	for _, item := range items {
		delete(s.store, item)
	}
}

// Len returns the number of items in the store
func (s *SimpleSetStore[T]) Len() int {
	return len(s.store)
}

// IsEmpty returns true if there are no items in the store
func (s *SimpleSetStore[T]) IsEmpty() bool {
	return len(s.store) == 0
}

// Contains returns whether an item is in the store
func (s *SimpleSetStore[T]) Contains(item T) bool {
	_, ok := s.store[item]
	return ok
}

// Pop removes an arbitrary item from the store and returns it. Returns error if the store is empty
func (s *SimpleSetStore[T]) Pop() (T, error) {
	var item T
	if s.IsEmpty() {
		return item, errors.New("set is empty")
	}
	for item = range s.store {
		break
	}
	s.Discard(item)
	return item, nil
}

// Items returns a slice of all the Set items
func (s *SimpleSetStore[T]) Items() []T {
	items := make([]T, 0, s.Len())
	for item := range s.store {
		items = append(items, item)
	}
	return items
}

// For runs a function on all the items in the store
func (s *SimpleSetStore[T]) For(f func(item T)) {
	for item := range s.store {
		f(item)
	}
}

// ForWithBreak runs a function on all the items in the store
// if f returns false, the iteration stops
func (s *SimpleSetStore[T]) ForWithBreak(f func(item T) bool) {
	for item := range s.store {
		if f(item) == false {
			break
		}
	}
}

func (s *SimpleSetStore[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Items())
}

func (s *SimpleSetStore[T]) UnmarshalJSON(b []byte) error {
	var items []T
	err := json.Unmarshal(b, &items)
	if err != nil {
		return err
	}
	s.Add(items...)
	return nil
}
