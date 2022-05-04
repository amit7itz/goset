package goset

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSet_Len(t *testing.T) {
	s := NewSet[string]()
	require.Equal(t, s.Len(), 0)
	s.Add("a")
	s.Add("a")
	require.Equal(t, s.Len(), 1)
	s.Add("b")
	s.Add("c")
	require.Equal(t, s.Len(), 3)
}

func TestSet_IsEmpty(t *testing.T) {
	s := NewSet[string]()
	require.True(t, s.IsEmpty())
	s.Add("a")
	s.Add("a")
	require.False(t, s.IsEmpty())
	s.Add("b")
	s.Add("c")
	require.False(t, s.IsEmpty())
}

func TestSet_ToSlice(t *testing.T) {
	s1 := NewSet[string]()
	s1.Add("a")
	s1.Add("b")
	s1.Add("c")
	s2 := FromSlice(s1.ToSlice())
	require.True(t, s1.Eq(s2))
}

func TestFromSlice(t *testing.T) {
	s1 := NewSet[string]()
	s1.Add("c")
	s1.Add("b")
	s1.Add("a")
	s2 := FromSlice[string]([]string{"a", "b", "c"})
	require.True(t, s1.Eq(s2))
}

func TestSet_Union(t *testing.T) {
	s1 := FromSlice[string]([]string{"a"})
	s2 := FromSlice[string]([]string{"b", "c"})
	s3 := FromSlice[string]([]string{"d", "e", "f"})
	union := s1.Union(s2, s3)
	require.Equal(t, s1.Len(), 1)
	require.Equal(t, s2.Len(), 2)
	require.Equal(t, s3.Len(), 3)
	require.Equal(t, union.Len(), 6)
	require.True(t, union.Eq(FromSlice[string]([]string{"a", "b", "c", "d", "e", "f"})))
}

func TestSet_Eq(t *testing.T) {
	s1 := FromSlice[string]([]string{"a", "b"})
	s2 := FromSlice[string]([]string{"b", "a"})
	require.True(t, s1.Eq(s2))
	s2.Discard("a")
	require.False(t, s1.Eq(s2))
}

func TestSet_Copy(t *testing.T) {
	s1 := FromSlice[string]([]string{"a", "b"})
	s2 := s1.Copy()
	require.True(t, s1 != s2)
	require.True(t, s1.Eq(s2))
}

func TestSet_Update(t *testing.T) {
	s := NewSet[int]()
	s.Add(1)
	s.Add(2)
	s2 := NewSet[int]()
	s2.Add(3)
	s2.Update(s)
	require.Equal(t, 2, s.Len())
	require.Equal(t, 3, s2.Len())
	require.True(t, s2.Contains(3))
}

func TestSet_Remove(t *testing.T) {
	s := NewSet[int]()
	s.Add(1)
	require.NoError(t, s.Remove(1))
	require.Error(t, s.Remove(1)) // should return error if item not found
}

func TestSet_Pop(t *testing.T) {
	s := NewSet[string]()
	s.Add("a")
	require.Equal(t, s.Len(), 1)
	item, err := s.Pop()
	require.NoError(t, err)
	require.Equal(t, s.Len(), 0)
	require.Equal(t, item, "a")
	_, err = s.Pop()
	require.Error(t, err)
}

func TestSet_Intersection(t *testing.T) {
	s1 := FromSlice[string]([]string{"a", "b", "c", "d", "e", "f"})
	s2 := FromSlice[string]([]string{"a", "", "c", "d", "e"})
	s3 := FromSlice[string]([]string{"z", "d", "e", "k"})
	intersection := s1.Intersection(s2, s3)
	require.True(t, intersection.Eq(FromSlice[string]([]string{"e", "d"})))
}

func TestSet_Difference(t *testing.T) {
	s1 := FromSlice[string]([]string{"a", "b", "c", "d", "e", "f"})
	s2 := FromSlice[string]([]string{"a", "", "c", "d", "e"})
	s3 := FromSlice[string]([]string{"z", "d", "e", "k"})
	difference := s1.Difference(s2, s3)
	require.True(t, difference.Eq(FromSlice[string]([]string{"b", "f"})))
}

func TestSet_SymmetricDifference(t *testing.T) {
	s1 := FromSlice[string]([]string{"a", "b", "c", "d", "e", "f"})
	s2 := FromSlice[string]([]string{"z", "d", "e", "k"})
	difference := s1.SymmetricDifference(s2)
	require.True(t, difference.Eq(FromSlice[string]([]string{"a", "b", "c", "f", "z", "k"})))
}

func TestSet_IsSubset(t *testing.T) {
	s1 := FromSlice[string]([]string{"a", "b", "c", "d", "e", "f"})
	s2 := FromSlice[string]([]string{"z", "d", "e", "k"})
	require.False(t, s2.IsSubset(s1))
	s3 := FromSlice[string]([]string{"b", "c"})
	require.True(t, s3.IsSubset(s1))
	s4 := NewSet[string]()
	require.True(t, s4.IsSubset(s1))
}

func TestSet_IsSuperset(t *testing.T) {
	s1 := FromSlice[string]([]string{"a", "b", "c", "d", "e", "f"})
	s2 := FromSlice[string]([]string{"z", "d", "e", "k"})
	require.False(t, s1.IsSuperset(s2))
	s3 := FromSlice[string]([]string{"b", "c"})
	require.True(t, s1.IsSuperset(s3))
	s4 := NewSet[string]()
	require.True(t, s1.IsSuperset(s4))
}

func TestSet_IsDisjoint(t *testing.T) {
	s1 := FromSlice[string]([]string{"a", "b", "c", "d", "e", "f"})
	s2 := FromSlice[string]([]string{"z", "d", "e", "k"})
	require.False(t, s1.IsDisjoint(s2))
	s3 := FromSlice[string]([]string{"g", "h", "i"})
	require.True(t, s1.IsDisjoint(s3))
}

func TestSet_String(t *testing.T) {
	s := FromSlice[string]([]string{"a", "b"})
	str := fmt.Sprintf("%v", s)
	possibleOutputs := []string{"Set{a b}", "Set{b a}"}
	require.Contains(t, possibleOutputs, str)
}
