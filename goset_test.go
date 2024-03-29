package goset

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSet_Len(t *testing.T) {
	s := NewSet[string]()
	require.Equal(t, s.Len(), 0)
	s.Add("a", "a")
	require.Equal(t, s.Len(), 1)
	s.Add("b", "c")
	require.Equal(t, s.Len(), 3)
}

func TestSet_IsEmpty(t *testing.T) {
	s := NewSet[string]()
	require.True(t, s.IsEmpty())
	s.Add("a", "b")
	require.False(t, s.IsEmpty())
	s.Add("b", "c")
	require.False(t, s.IsEmpty())
}

func BenchmarkSet_Items(b *testing.B) {
	s1 := NewSet[string]("a", "b", "c")

	for i := 0; i < b.N; i++ {
		s1.Items()
	}
}

func TestSet_Items(t *testing.T) {
	s1 := NewSet[string]("a", "b", "c")
	s2 := FromSlice(s1.Items())
	require.True(t, s1.Equal(s2))
}

func TestSet_For(t *testing.T) {
	s1 := NewSet[string]("a", "b", "c")
	s2 := NewSet[string]("a", "b", "c")
	counter := 0
	s1.For(func(item string) {
		s2.Discard(item)
		counter++
	})
	require.Equal(t, counter, s1.Len())
	require.True(t, s2.IsEmpty())
}

func TestSet_ForWithBreak(t *testing.T) {
	s1 := NewSet[string]("a", "b", "c")
	s2 := NewSet[string]("a", "b", "c")
	counter := 0
	s1.ForWithBreak(func(item string) bool {
		if counter == 2 {
			return false
		}
		counter++
		s2.Discard(item)
		return true
	})
	require.Equal(t, 2, counter)
	require.Equal(t, 1, s2.Len())
}

func BenchmarkFromSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FromSlice[string]([]string{"a", "b", "c"})
	}
}

func TestFromSlice(t *testing.T) {
	s1 := NewSet[string]("c", "b", "a")
	s2 := FromSlice[string]([]string{"a", "b", "c"})
	require.True(t, s1.Equal(s2))
}

func TestSet_Union(t *testing.T) {
	s1 := NewSet[string]("a")
	s2 := NewSet[string]("b", "c")
	s3 := NewSet[string]("d", "e", "f")
	union := s1.Union(s2, s3)
	require.Equal(t, s1.Len(), 1)
	require.Equal(t, s2.Len(), 2)
	require.Equal(t, s3.Len(), 3)
	require.Equal(t, union.Len(), 6)
	require.True(t, union.Equal(NewSet[string]("a", "b", "c", "d", "e", "f")))
}

func TestSet_Equal(t *testing.T) {
	s1 := NewSet[string]("a", "b")
	s2 := NewSet[string]("b", "a")
	require.True(t, s1.Equal(s2))
	s2.Discard("a")
	require.False(t, s1.Equal(s2))
	s3 := NewSet[string]()
	s4 := NewSet[string]()
	require.True(t, s3.Equal(s4))
	require.False(t, s3.Equal(s1))
}

func TestSet_Copy(t *testing.T) {
	s1 := NewSet[string]("a", "b")
	s2 := s1.Copy()
	require.True(t, s1 != s2)
	require.True(t, s1.Equal(s2))
}

func TestSet_Update(t *testing.T) {
	s := NewSet[int]()
	s.Add(1, 2)
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
	s1 := NewSet[string]("a", "b", "c", "d", "e", "f")
	s2 := NewSet[string]("a", "", "c", "d", "e")
	s3 := NewSet[string]("z", "d", "e", "k")
	intersection := s1.Intersection(s2, s3)
	require.True(t, intersection.Equal(NewSet[string]("e", "d")))
}

func TestSet_Difference(t *testing.T) {
	s1 := NewSet[string]("a", "b", "c", "d", "e", "f")
	s2 := NewSet[string]("a", "", "c", "d", "e")
	s3 := NewSet[string]("z", "d", "e", "k")
	difference := s1.Difference(s2, s3)
	require.True(t, difference.Equal(NewSet[string]("b", "f")))
}

func TestSet_SymmetricDifference(t *testing.T) {
	s1 := NewSet[string]("a", "b", "c", "d", "e", "f")
	s2 := NewSet[string]("z", "d", "e", "k")
	difference := s1.SymmetricDifference(s2)
	require.True(t, difference.Equal(NewSet[string]("a", "b", "c", "f", "z", "k")))
}

func TestSet_IsSubset(t *testing.T) {
	s1 := NewSet[string]("a", "b", "c", "d", "e", "f")
	s2 := NewSet[string]("z", "d", "e", "k")
	require.False(t, s2.IsSubset(s1))
	s3 := NewSet[string]("b", "c")
	require.True(t, s3.IsSubset(s1))
	s4 := NewSet[string]()
	require.True(t, s4.IsSubset(s1))
}

func TestSet_IsSuperset(t *testing.T) {
	s1 := NewSet[string]("a", "b", "c", "d", "e", "f")
	s2 := NewSet[string]("z", "d", "e", "k")
	require.False(t, s1.IsSuperset(s2))
	s3 := NewSet[string]("b", "c")
	require.True(t, s1.IsSuperset(s3))
	s4 := NewSet[string]()
	require.True(t, s1.IsSuperset(s4))
}

func TestSet_IsDisjoint(t *testing.T) {
	s1 := NewSet[string]("a", "b", "c", "d", "e", "f")
	s2 := NewSet[string]("z", "d", "e", "k")
	require.False(t, s1.IsDisjoint(s2))
	s3 := NewSet[string]("g", "h", "i")
	require.True(t, s1.IsDisjoint(s3))
}

func TestSet_String(t *testing.T) {
	s := NewSet[string]("a", "b")
	str := fmt.Sprintf("%v", s)
	possibleOutputs := []string{"Set[string]{a b}", "Set[string]{b a}"}
	require.Contains(t, possibleOutputs, str)
}

func TestSetWithStruct(t *testing.T) {
	type Person struct {
		Name string
	}
	peopleSet := NewSet(Person{Name: "Amit"}, Person{Name: "Amit"})
	require.Equal(t, peopleSet.Len(), 1)
}

type DummyStructWithMap struct {
	A int                     `json:"a"`
	B string                  `json:"b"`
	S *Set[string]            `json:"s"`
	M map[string]*Set[string] `json:"m"`
}

func TestSet_MarshalJSON(t *testing.T) {
	s1 := NewSet[string]("a", "b", "c", "d", "e", "f")
	bytes, err := json.Marshal(s1)
	require.NoError(t, err)
	s2 := NewSet[string]()
	err = json.Unmarshal(bytes, &s2)
	require.NoError(t, err)
	require.True(t, s1.Equal(s2))

	d := DummyStructWithMap{A: 123, B: "test string", S: s1, M: map[string]*Set[string]{"bla": s1}}
	bytes, err = json.Marshal(d)
	require.NoError(t, err)
	d2 := DummyStructWithMap{}
	err = json.Unmarshal(bytes, &d2)
	require.NoError(t, err)
	require.True(t, d.M["bla"].Equal(d2.M["bla"]))
}
