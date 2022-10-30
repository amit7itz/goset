package goset

import (
	"fmt"
	"github.com/amit7itz/goset/store"
	"reflect"
	"strings"
)

// Set represents a set data structure.
// You should not call it directly, use NewSet() or FromSlice()
type Set[T comparable] struct {
	store store.SetStore[T]
}

// NewSet returns a new Set of the given items
func NewSet[T comparable](items ...T) *Set[T] {
	set := &Set[T]{store: store.NewSimpleStore[T]()}
	set.Add(items...)
	return set
}

// FromSlice returns a new Set with all the items of the slice.
func FromSlice[T comparable](slice []T) *Set[T] {
	set := NewSet[T]()
	set.Add(slice...)
	return set
}

// Add adds item(s) to the Set
func (s *Set[T]) Add(items ...T) {
	s.store.Add(items...)
}

// Remove removes a single item from the Set. Returns error if the item is not in the Set
// See also: Discard()
func (s *Set[T]) Remove(item T) error {
	return s.store.Remove(item)
}

// Discard removes item(s) from the Set if exist
// See also: Remove()
func (s *Set[T]) Discard(items ...T) {
	s.store.Discard(items...)
}

// Len returns the number of items in the Set
func (s *Set[T]) Len() int {
	return s.store.Len()
}

// IsEmpty returns true if there are no items in the Set
func (s *Set[T]) IsEmpty() bool {
	return s.store.IsEmpty()
}

// Contains returns whether an item is in the Set
func (s *Set[T]) Contains(item T) bool {
	return s.store.Contains(item)
}

// Pop removes an arbitrary item from the Set and returns it. Returns error if the Set is empty
func (s *Set[T]) Pop() (T, error) {
	return s.store.Pop()
}

// Items returns a slice of all the Set items
func (s *Set[T]) Items() []T {
	return s.store.Items()
}

// For runs a function on all the items in the Set
func (s *Set[T]) For(f func(item T)) {
	s.store.For(f)
}

// ForWithBreak runs a function on all the items in the store
// if f returns false, the iteration stops
func (s *Set[T]) ForWithBreak(f func(item T) bool) {
	s.store.ForWithBreak(f)
}

// String returns a string that represents the Set
func (s *Set[T]) String() string {
	var t T
	str := fmt.Sprintf("Set[%s]{", reflect.TypeOf(t).String())
	itemsStr := make([]string, 0, s.Len())
	s.store.For(func(item T) {
		itemsStr = append(itemsStr, fmt.Sprintf("%v", item))
	})
	str += strings.Join(itemsStr, " ")
	str += "}"
	return str
}

// Update adds all the items from the other Sets to the current Set
func (s *Set[T]) Update(others ...*Set[T]) {
	for _, other := range others {
		other.store.For(func(item T) {
			s.Add(item)
		})
	}
}

// Copy returns a new Set with the same items as the current Set
func (s *Set[T]) Copy() *Set[T] {
	set := NewSet[T]()
	s.store.For(func(item T) {
		set.Add(item)
	})
	return set
}

// Equal returns whether the current Set contains the same items as the other one
func (s *Set[T]) Equal(other *Set[T]) bool {
	if s.Len() != other.Len() {
		return false
	}
	equal := true
	s.store.ForWithBreak(func(item T) bool {
		if !other.Contains(item) {
			equal = false
			return false // stop iteration
		}
		return true
	})
	return equal
}

// Union returns a new Set of the items from the current set and all others
func (s *Set[T]) Union(others ...*Set[T]) *Set[T] {
	unionSet := s.Copy()
	unionSet.Update(others...)
	return unionSet
}

// Intersection returns a new Set with the common items of the current set and all others.
func (s *Set[T]) Intersection(others ...*Set[T]) *Set[T] {
	intersectionSet := NewSet[T]()
	s.store.For(func(item T) {
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
	})
	return intersectionSet
}

// Difference returns a new Set of all the items in the current Set that are not in any of the others
func (s *Set[T]) Difference(others ...*Set[T]) *Set[T] {
	differenceSet := NewSet[T]()
	s.store.For(func(item T) {
		inAnyOther := false
		for _, other := range others {
			if other.Contains(item) {
				inAnyOther = true
				break
			}
		}
		if !inAnyOther {
			differenceSet.Add(item)
		}
	})
	return differenceSet
}

// SymmetricDifference returns all the items that exist in only one of the Sets
func (s *Set[T]) SymmetricDifference(other *Set[T]) *Set[T] {
	symmetricDifferenceSet := NewSet[T]()
	s.store.For(func(item T) {
		if !other.Contains(item) {
			symmetricDifferenceSet.Add(item)
		}
	})
	other.store.For(func(item T) {
		if !s.Contains(item) {
			symmetricDifferenceSet.Add(item)
		}
	})
	return symmetricDifferenceSet
}

// IsDisjoint returns whether the two Sets have no item in common
func (s *Set[T]) IsDisjoint(other *Set[T]) bool {
	intersection := s.Intersection(other)
	return intersection.IsEmpty()
}

// IsSubset returns whether all the items of the current set exist in the other one
func (s *Set[T]) IsSubset(other *Set[T]) bool {
	intersection := s.Intersection(other)
	return intersection.Len() == s.Len()
}

// IsSuperset returns whether all the items of the other set exist in the current one
func (s *Set[T]) IsSuperset(other *Set[T]) bool {
	return other.IsSubset(s)
}

func (s *Set[T]) MarshalJSON() ([]byte, error) {
	return s.store.MarshalJSON()
}

func (s *Set[T]) UnmarshalJSON(b []byte) error {
	if s.store == nil {
		s.store = store.NewSimpleStore[T]()
	}
	return s.store.UnmarshalJSON(b)
}
