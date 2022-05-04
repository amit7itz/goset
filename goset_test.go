package goset

import (
	"github.com/stretchr/testify/require"
	"testing"
)

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
