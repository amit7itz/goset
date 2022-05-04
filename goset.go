package goset

import (
	"fmt"
	"reflect"
)

type Set[T comparable] struct {
	store map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{store: make(map[T]struct{}, 0)}
}

func FromSlice[T comparable](slice []T) *Set[T] {
	set := NewSet[T]()
	for _, item := range slice {
		set.Add(item)
	}
	return set
}

func (s *Set[T]) Add(item T) {
	s.store[item] = struct{}{}
}

func (s *Set[T]) Remove(item T) error {
	if s.Contains(item) {
		delete(s.store, item)
		return nil
	}
	return fmt.Errorf("KeyError: %v", item)
}

func (s *Set[T]) Discard(item T) {
	delete(s.store, item)
}

func (s *Set[T]) Len() int {
	return len(s.store)
}

func (s *Set[T]) Contains(item T) bool {
	_, ok := s.store[item]
	return ok
}

func (s *Set[T]) Update(sets ...*Set[T]) {
	for _, set := range sets {
		for item := range set.store {
			s.Add(item)
		}
	}
}

func (s *Set[T]) Copy() *Set[T] {
	set := NewSet[T]()
	for item := range s.store {
		set.Add(item)
	}
	return set
}

func (s *Set[T]) ToSlice() []T {
	slice := make([]T, 0)
	for item := range s.store {
		slice = append(slice, item)
	}
	return slice
}

func (s *Set[T]) Eq(set *Set[T]) bool {
	return reflect.DeepEqual(s, set)
}

func (s *Set[T]) Union(sets ...*Set[T]) *Set[T] {
	unionSet := s.Copy()
	unionSet.Update(sets...)
	return unionSet
}

func (s *Set[T]) Intersection(sets ...*Set[T]) *Set[T] {
	intersectionSet := NewSet[T]()
	for item := range s.store {
		inAllOthers := true
		for _, set := range sets {
			if !set.Contains(item) {
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

func (s *Set[T]) Difference(sets ...*Set[T]) *Set[T] {
	differenceSet := NewSet[T]()
	for item := range s.store {
		inAnyOther := false
		for _, set := range sets {
			if set.Contains(item) {
				inAnyOther = true
				break
			}
		}
		if !inAnyOther {
			differenceSet.Add(item)
		}
	}
	return differenceSet
}

func (s *Set[T]) SymmetricDifference(set *Set[T]) *Set[T] {
	symmetricDifferenceSet := NewSet[T]()
	for item := range s.store {
		if !set.Contains(item) {
			symmetricDifferenceSet.Add(item)
		}
	}
	for item := range set.store {
		if !s.Contains(item) {
			symmetricDifferenceSet.Add(item)
		}
	}
	return symmetricDifferenceSet
}
