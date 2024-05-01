package collections

import (
	"golang.org/x/exp/constraints"
	"sort"
)

type SortedList[T constraints.Ordered] []T

func NewSortedList[T constraints.Ordered]() SortedList[T] {
	return make(SortedList[T], 0)
}

func (sl *SortedList[T]) Insert(value T) {
	*sl = append(*sl, value)
	sort.Slice(*sl, func(i, j int) bool {
		return (*sl)[i] < (*sl)[j]
	})
}

func (sl *SortedList[T]) Delete(value T) bool {
	for i, v := range *sl {
		if v == value {
			*sl = append((*sl)[:i], (*sl)[i+1:]...)
			return true
		}
	}
	return false
}

func (sl *SortedList[T]) Get(index int) (T, bool) {
	if index < 0 || index >= len(*sl) {
		var t T
		return t, false
	}
	return (*sl)[index], true
}
