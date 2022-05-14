package goset

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestSafeSet_Len(t *testing.T) {
	s := NewSafeSet[string]()
	require.Equal(t, s.Len(), 0)
	s.Add("a", "a")
	require.Equal(t, s.Len(), 1)
	s.Add("b", "c")
	require.Equal(t, s.Len(), 3)
}

func TestSafeSet_Len_Concurrency(t *testing.T) {
	s := NewSafeSet[int](1, 2, 3)

	concurrentFor(func(_ int) {
		s.Len()
		s.Add(1)
	})

	require.Equal(t, 3, s.Len())
}

func TestSafeSet_IsEmpty(t *testing.T) {
	s := NewSafeSet[string]()
	require.True(t, s.IsEmpty())
	s.Add("a", "b")
	require.False(t, s.IsEmpty())
	s.Add("b", "c")
	require.False(t, s.IsEmpty())
}

func TestSafeSet_IsEmpty_Concurrency(t *testing.T) {
	s := NewSafeSet[int](1, 2, 3)

	concurrentFor(func(_ int) {
		s.IsEmpty()
		s.Add(1)
	})

	require.Equal(t, 3, s.Len())
}

func BenchmarkSafeSet_Items(b *testing.B) {
	s1 := NewSafeSet[string]("a", "b", "c")

	for i := 0; i < b.N; i++ {
		s1.Items()
	}
}

func TestSafeSet_Items(t *testing.T) {
	s1 := NewSafeSet[string]("a", "b", "c")
	s2 := SafeFromSlice(s1.Items())
	require.True(t, s1.Equal(s2))
}

func TestSafeSet_Items_Concurrency(t *testing.T) {
	s := NewSafeSet[int](1, 2, 3)

	concurrentFor(func(_ int) {
		s.Items()
		s.Add(1)
	})

	require.Equal(t, 3, s.Len())
}

func BenchmarkSafeFromSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SafeFromSlice[string]([]string{"a", "b", "c"})
	}
}

func TestSafeFromSlice(t *testing.T) {
	s1 := NewSafeSet[string]("c", "b", "a")
	s2 := SafeFromSlice[string]([]string{"a", "b", "c"})
	require.True(t, s1.Equal(s2))
}

func TestSafeSet_Union(t *testing.T) {
	s1 := NewSafeSet[string]("a")
	s2 := NewSafeSet[string]("b", "c")
	s3 := NewSafeSet[string]("d", "e", "f")
	union := s1.Union(s2, s3)
	require.Equal(t, s1.Len(), 1)
	require.Equal(t, s2.Len(), 2)
	require.Equal(t, s3.Len(), 3)
	require.Equal(t, union.Len(), 6)
	require.True(t, union.Equal(NewSafeSet[string]("a", "b", "c", "d", "e", "f")))
}

func TestSafeSet_Union_Concurrency(t *testing.T) {
	s := NewSafeSet[int](1, 2, 3)
	s1 := NewSafeSet[int](1, 2, 3)

	concurrentForOthers(func(_ int) {
		_ = s.Union(s1)
	}, func(i int) {
		s1.Add(i)
	})

	require.Equal(t, 3, s.Len())
}

func TestSafeSet_Equal(t *testing.T) {
	s1 := NewSafeSet[string]("a", "b")
	s2 := NewSafeSet[string]("b", "a")
	require.True(t, s1.Equal(s2))
	s2.Discard("a")
	require.False(t, s1.Equal(s2))
	s3 := NewSafeSet[string]()
	s4 := NewSafeSet[string]()
	require.True(t, s3.Equal(s4))
	require.False(t, s3.Equal(s1))
}

func TestSafeSet_Equal_Concurrency(t *testing.T) {
	s := NewSafeSet[int](1, 2, 3)
	s1 := NewSafeSet[int](1, 2, 3)

	concurrentFor(func(_ int) {
		_ = s.Equal(s1)
		s.Add(1)
	})

	require.Equal(t, 3, s.Len())
}

func TestSafeSet_Copy(t *testing.T) {
	s1 := NewSafeSet[string]("a", "b")
	s2 := s1.Copy()
	require.True(t, s1 != s2)
	require.True(t, s1.Equal(s2))
}

func TestSafeSet_Copy_Concurrency(t *testing.T) {
	s := NewSafeSet[int](1, 2, 3)

	concurrentFor(func(_ int) {
		s.Copy()
		s.Add(1)
	})

	require.Equal(t, 3, s.Len())
}

func TestSafeSet_Update(t *testing.T) {
	s := NewSafeSet[int]()
	s.Add(1, 2)
	s2 := NewSafeSet[int]()
	s2.Add(3)
	s2.Update(s)
	require.Equal(t, 2, s.Len())
	require.Equal(t, 3, s2.Len())
	require.True(t, s2.Contains(3))
}

func TestSafeSet_Update_Concurrency(t *testing.T) {
	s := NewSafeSet[int]()
	s1 := SafeFromSlice([]int{1, 2, 3})
	s2 := SafeFromSlice([]int{4, 5, 6})

	concurrentFor(func(_ int) {
		s.Update(s1, s2)
	})

	require.Equal(t, 6, s.Len())
}

func TestSafeSet_Remove(t *testing.T) {
	s := NewSafeSet[int]()
	s.Add(1)
	require.NoError(t, s.Remove(1))
	require.Error(t, s.Remove(1)) // should return error if item not found
}

func TestSafeSet_Remove_Concurrency(t *testing.T) {
	s := NewSafeSet[int](1, 2, 3)

	concurrentFor(func(_ int) {
		_ = s.Remove(2)
		s.Add(1)
	})

	require.Equal(t, 2, s.Len())
}

func TestSafeSet_Pop(t *testing.T) {
	s := NewSafeSet[string]()
	s.Add("a")
	require.Equal(t, s.Len(), 1)
	item, err := s.Pop()
	require.NoError(t, err)
	require.Equal(t, s.Len(), 0)
	require.Equal(t, item, "a")
	_, err = s.Pop()
	require.Error(t, err)
}

func TestSafeSet_Pop_Concurrency(t *testing.T) {
	s := NewSafeSet[int](1, 2, 3)

	concurrentFor(func(_ int) {
		_, _ = s.Pop()
	})

	require.Equal(t, 0, s.Len())
}

func TestSafeSet_Intersection(t *testing.T) {
	s1 := NewSafeSet[string]("a", "b", "c", "d", "e", "f")
	s2 := NewSafeSet[string]("a", "", "c", "d", "e")
	s3 := NewSafeSet[string]("z", "d", "e", "k")
	intersection := s1.Intersection(s2, s3)
	require.True(t, intersection.Equal(NewSafeSet[string]("e", "d")))
}

func TestSafeSet_Intersection_Concurrency(t *testing.T) {
	s := NewSafeSet[int](1, 2, 3)
	s1 := NewSafeSet[int](1, 3)

	concurrentForOthers(func(_ int) {
		_ = s.Intersection(s1)
	}, func(i int) {
		s1.Add(i)
	})

	require.Equal(t, 3, s.Len())
}

func TestSafeSet_Difference(t *testing.T) {
	s1 := NewSafeSet[string]("a", "b", "c", "d", "e", "f")
	s2 := NewSafeSet[string]("a", "", "c", "d", "e")
	s3 := NewSafeSet[string]("z", "d", "e", "k")
	difference := s1.Difference(s2, s3)
	require.True(t, difference.Equal(NewSafeSet[string]("b", "f")))
}

func TestSafeSet_Difference_Concurrency(t *testing.T) {
	s := NewSafeSet[int](1, 2, 3)
	s1 := NewSafeSet[int](1, 3)

	concurrentFor(func(_ int) {
		_ = s.Difference(s1)
		s.Add(1)
	})

	require.Equal(t, 3, s.Len())
}

func TestSafeSet_SymmetricDifference(t *testing.T) {
	s1 := NewSafeSet[string]("a", "b", "c", "d", "e", "f")
	s2 := NewSafeSet[string]("z", "d", "e", "k")
	difference := s1.SymmetricDifference(s2)
	require.True(t, difference.Equal(NewSafeSet[string]("a", "b", "c", "f", "z", "k")))
}

func TestSafeSet_SymmetricDifference_Concurrency(t *testing.T) {
	s := NewSafeSet[int](1, 2, 3)
	s1 := NewSafeSet[int](1, 3)

	concurrentFor(func(_ int) {
		_ = s.SymmetricDifference(s1)
		s.Add(1)
	})

	require.Equal(t, 3, s.Len())
}

func TestSafeSet_IsSubset(t *testing.T) {
	s1 := NewSafeSet[string]("a", "b", "c", "d", "e", "f")
	s2 := NewSafeSet[string]("z", "d", "e", "k")
	require.False(t, s2.IsSubset(s1))
	s3 := NewSafeSet[string]("b", "c")
	require.True(t, s3.IsSubset(s1))
	s4 := NewSafeSet[string]()
	require.True(t, s4.IsSubset(s1))
}

func TestSafeSet_IsSubset_Concurrency(t *testing.T) {
	s := NewSafeSet[int](1, 2, 3)
	s1 := NewSafeSet[int](1, 3)

	concurrentFor(func(_ int) {
		s1.IsSubset(s)
		s.Add(1)
	})

	require.Equal(t, 3, s.Len())
}

func TestSafeSet_IsSuperset(t *testing.T) {
	s1 := NewSafeSet[string]("a", "b", "c", "d", "e", "f")
	s2 := NewSafeSet[string]("z", "d", "e", "k")
	require.False(t, s1.IsSuperset(s2))
	s3 := NewSafeSet[string]("b", "c")
	require.True(t, s1.IsSuperset(s3))
	s4 := NewSafeSet[string]()
	require.True(t, s1.IsSuperset(s4))
}

func TestSafeSet_IsSuperset_Concurrency(t *testing.T) {
	s := NewSafeSet[int](1, 2, 3)
	s1 := NewSafeSet[int](1, 3)

	concurrentFor(func(_ int) {
		s1.IsSuperset(s)
		s.Add(1)
	})

	require.Equal(t, 3, s.Len())
}

func TestSafeSet_IsDisjoint(t *testing.T) {
	s1 := NewSafeSet[string]("a", "b", "c", "d", "e", "f")
	s2 := NewSafeSet[string]("z", "d", "e", "k")
	require.False(t, s1.IsDisjoint(s2))
	s3 := NewSafeSet[string]("g", "h", "i")
	require.True(t, s1.IsDisjoint(s3))
}

func TestSafeSet_IsDisjoint_Concurrency(t *testing.T) {
	s := NewSafeSet[int](1, 2, 3)
	s1 := NewSafeSet[int](1, 3)

	concurrentFor(func(_ int) {
		s1.IsDisjoint(s)
		s.Add(1)
	})

	require.Equal(t, 3, s.Len())
}

func TestSafeSet_String(t *testing.T) {
	s := NewSafeSet[string]("a", "b")
	str := fmt.Sprintf("%v", s)
	possibleOutputs := []string{"Set[string]{a b}", "Set[string]{b a}"}
	require.Contains(t, possibleOutputs, str)
}

func TestSafeSet_String_Concurrency(t *testing.T) {
	s := NewSafeSet[int](1, 2, 3)

	concurrentFor(func(_ int) {
		_ = s.String()
		s.Add(1)
	})

	require.Equal(t, 3, s.Len())
}

func TestSafeSetWithStruct(t *testing.T) {
	type Person struct {
		Name string
	}
	peopleSet := NewSafeSet(Person{Name: "Amit"}, Person{Name: "Amit"})
	require.Equal(t, peopleSet.Len(), 1)
}

func TestSafeSet_Add_Concurrency(t *testing.T) {
	s := NewSafeSet[int]()

	concurrentFor(func(i int) {
		s.Add(i)
	})

	require.Equal(t, 10, s.Len())
}

func TestSafeSet_Discard_Concurrency(t *testing.T) {
	s := NewSafeSet[int](1, 2, 3)

	concurrentFor(func(_ int) {
		s.Discard(2)
		s.Add(1)
	})

	require.Equal(t, 2, s.Len())
}

func TestSafeSet_Contains_Concurrency(t *testing.T) {
	s := NewSafeSet[int](1, 2, 3)

	concurrentFor(func(_ int) {
		_ = s.Contains(2)
		s.Add(1)
	})

	require.Equal(t, 3, s.Len())
}

func concurrentFor(fn func(i int)) {
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			fn(i)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func concurrentForOthers(fn1 func(i int), fn2 func(i int)) {
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			fn1(i)
			wg.Done()
		}(i)
		go func(i int) {
			fn2(i)
		}(i)
	}

	wg.Wait()
}

func TestSafeSet_Equal_Deadlock_Concurrency(t *testing.T) {
	s1 := NewSafeSet[int](1, 2, 3)
	s2 := NewSafeSet[int](1, 2, 3)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func() {
			s1.Equal(s2)
			wg.Done()
		}()
		go func() {
			s2.Equal(s1)
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestSafeSet_Union_Deadlock_Concurrency(t *testing.T) {
	s1 := NewSafeSet[int](1, 2, 3)
	s2 := NewSafeSet[int](1, 2, 3)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func() {
			s1.Union(s2)
			wg.Done()
		}()
		go func() {
			s2.Union(s1)
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestSafeSet_Intersection_Deadlock_Concurrency(t *testing.T) {
	s1 := NewSafeSet[int](1, 2, 3)
	s2 := NewSafeSet[int](1, 2, 3)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func() {
			s1.Intersection(s2)
			wg.Done()
		}()
		go func() {
			s2.Intersection(s1)
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestSafeSet_Difference_Deadlock_Concurrency(t *testing.T) {
	s1 := NewSafeSet[int](1, 2, 3)
	s2 := NewSafeSet[int](1, 2, 3)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func() {
			s1.Difference(s2)
			wg.Done()
		}()
		go func() {
			s2.Difference(s1)
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestSafeSet_SymmetricDifference_Deadlock_Concurrency(t *testing.T) {
	s1 := NewSafeSet[int](1, 2, 3)
	s2 := NewSafeSet[int](1, 2, 3)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func() {
			s1.SymmetricDifference(s2)
			wg.Done()
		}()
		go func() {
			s2.SymmetricDifference(s1)
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestSafeSet_IsDisjoint_Deadlock_Concurrency(t *testing.T) {
	s1 := NewSafeSet[int](1, 2, 3)
	s2 := NewSafeSet[int](1, 2, 3)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func() {
			s1.IsDisjoint(s2)
			wg.Done()
		}()
		go func() {
			s2.IsDisjoint(s1)
			wg.Done()
		}()
	}

	wg.Wait()
}
