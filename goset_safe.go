package goset

import (
	"sync"
)

// SafeSet represents a thread safe set data structure.
// You should not call it directly, use NewSafeSet() or SafeFromSlice()
type SafeSet[T comparable] struct {
	store *Set[T]
	l     sync.Mutex
}

// NewSafeSet returns a new SafeSet of the given items
func NewSafeSet[T comparable](items ...T) *SafeSet[T] {
	return &SafeSet[T]{store: NewSet(items...)}
}

// SafeFromSlice returns a new SafeSet with all the items of the slice.
func SafeFromSlice[T comparable](slice []T) *SafeSet[T] {
	return NewSafeSet[T](slice...)
}

// SafeFromSet returns a new SafeSet with Set values.
func SafeFromSet[T comparable](s *Set[T]) *SafeSet[T] {
	return &SafeSet[T]{store: s}
}

// String returns a string that represents the SafeSet
func (s *SafeSet[T]) String() string {
	s.Lock()
	defer s.Unlock()
	return s.store.String()
}

// Add adds item(s) to the SafeSet
func (s *SafeSet[T]) Add(items ...T) {
	s.Lock()
	defer s.Unlock()
	s.store.Add(items...)
}

// Remove removes a single item from the SafeSet. Returns error if the item is not in the SafeSet
// See also: Discard()
func (s *SafeSet[T]) Remove(item T) error {
	s.Lock()
	defer s.Unlock()
	return s.store.Remove(item)
}

// Discard removes item(s) from the SafeSet if exist
// See also: Remove()
func (s *SafeSet[T]) Discard(items ...T) {
	s.Lock()
	defer s.Unlock()
	s.store.Discard(items...)
}

// Len returns the number of items in the SafeSet
func (s *SafeSet[T]) Len() int {
	s.Lock()
	defer s.Unlock()
	return s.store.Len()
}

// IsEmpty returns true if there are no items in the SafeSet
func (s *SafeSet[T]) IsEmpty() bool {
	return s.Len() == 0
}

// Contains returns whether an item is in the SafeSet
func (s *SafeSet[T]) Contains(item T) bool {
	s.Lock()
	defer s.Unlock()
	return s.store.Contains(item)
}

// Update adds all the items from the other SafeSets to the current SafeSet
func (s *SafeSet[T]) Update(others ...*SafeSet[T]) {
	s.Lock()
	defer s.Unlock()
	defer unlockOthers(others)
	s.store.Update(mapFunc(others, func(other *SafeSet[T]) *Set[T] {
		other.Lock()
		return other.GetSet()
	})...)
}

// Pop removes an arbitrary item from the SafeSet and returns it. Returns error if the SafeSet is empty
func (s *SafeSet[T]) Pop() (T, error) {
	s.Lock()
	defer s.Unlock()
	return s.store.Pop()
}

// Copy returns a new SafeSet with the same items as the current SafeSet
func (s *SafeSet[T]) Copy() *SafeSet[T] {
	s.Lock()
	defer s.Unlock()
	return SafeFromSet[T](s.store.Copy())
}

// Items returns a slice of all the SafeSet items
func (s *SafeSet[T]) Items() []T {
	s.Lock()
	defer s.Unlock()
	return s.store.Items()
}

// Equal returns whether the current SafeSet contains the same items as the other one
func (s *SafeSet[T]) Equal(other *SafeSet[T]) bool {
	other.Lock()
	set := other.GetSet()
	other.Unlock()
	s.Lock()
	defer s.Unlock()
	return s.store.Equal(set)
}

// Union returns a new SafeSet of the items from the current set and all others
func (s *SafeSet[T]) Union(others ...*SafeSet[T]) *SafeSet[T] {

	sets := mapFunc(others, func(other *SafeSet[T]) *Set[T] {
		other.Lock()
		set := other.GetSet().Copy()
		other.Unlock()
		return set
	})

	s.Lock()
	defer s.Unlock()
	return SafeFromSet[T](s.store.Union(sets...))
}

// Intersection returns a new SafeSet with the common items of the current set and all others.
func (s *SafeSet[T]) Intersection(others ...*SafeSet[T]) *SafeSet[T] {
	intersectionSet := NewSafeSet[T]()

	for _, item := range s.Items() {
		inAllOthers := true
		for _, other := range others {
			if !other.Contains(item) {
				inAllOthers = false
				break
			}
		}
		if inAllOthers {
			intersectionSet.Add(item)
		}
	}
	return intersectionSet
}

// Difference returns a new SafeSet of all the items in the current SafeSet that are not in any of the others
func (s *SafeSet[T]) Difference(others ...*SafeSet[T]) *SafeSet[T] {
	sets := mapFunc(others, func(other *SafeSet[T]) *Set[T] {
		other.Lock()
		set := other.GetSet()
		other.Unlock()
		return set
	})
	s.Lock()
	defer s.Unlock()
	return SafeFromSet[T](s.store.Difference(sets...))
}

// SymmetricDifference returns all the items that exist in only one of the SafeSets
func (s *SafeSet[T]) SymmetricDifference(other *SafeSet[T]) *SafeSet[T] {
	other.Lock()
	set := other.GetSet()
	other.Unlock()
	s.Lock()
	defer s.Unlock()
	return SafeFromSet[T](s.store.SymmetricDifference(set))
}

// IsDisjoint returns whether the two SafeSets have no item in common
func (s *SafeSet[T]) IsDisjoint(other *SafeSet[T]) bool {
	return s.Intersection(other).IsEmpty()
}

// IsSubset returns whether all the items of the current set exist in the other one
func (s *SafeSet[T]) IsSubset(other *SafeSet[T]) bool {
	return s.Intersection(other).Len() == s.Len()
}

// IsSuperset returns whether all the items of the other SafeSet exist in the current one
func (s *SafeSet[T]) IsSuperset(other *SafeSet[T]) bool {
	return other.IsSubset(s)
}

// Lock is proxy method for mutex
func (s *SafeSet[T]) Lock() {
	s.l.Lock()
}

// Unlock is proxy method for mutes
func (s *SafeSet[T]) Unlock() {
	s.l.Unlock()
}

// GetSet returns internal set
func (s *SafeSet[T]) GetSet() *Set[T] {
	return s.store
}

func unlockOthers[T comparable](others []*SafeSet[T]) {
	iterFunc(others, func(other *SafeSet[T]) {
		other.Unlock()
	})
}

// mapFunc is generic map over list
func mapFunc[T1, T2 any](input []T1, f func(T1) T2) (output []T2) {
	output = make([]T2, 0, len(input))
	for _, v := range input {
		output = append(output, f(v))
	}
	return output
}

// iterFunc is generic iterator over list
func iterFunc[T comparable](input []T, f func(T)) {
	for _, v := range input {
		f(v)
	}
}
