package collections

import (
	"golang.org/x/exp/constraints"
	"sort"
)

type SortedList[T comparable] struct {
	values []T
	sort   func(i, j T) bool
}

func NewEmptySortedList[T comparable](sort func(i, j T) bool) SortedList[T] {
	return SortedList[T]{
		values: make([]T, 0),
		sort:   sort,
	}
}

func NewSortedList[T constraints.Ordered](values []T, sort func(i, j T) bool) SortedList[T] {
	result := NewEmptySortedList[T](sort)
	for _, v := range values {
		result.Insert(v)
	}
	return result
}

func (sl *SortedList[T]) Insert(value T) {
	values := sl.values
	values = append(values, value)
	sort.Slice(values, func(i, j int) bool {
		return (*sl).sort(values[i], values[j])
	})
	sl.values = values
}

func (sl *SortedList[T]) Delete(value T) bool {
	values := sl.values
	for i, v := range values {
		if v == value {
			values = append(values[:i], values[i+1:]...)
			sl.values = values
			return true
		}
	}
	return false
}

func (sl *SortedList[T]) Get(index int) (T, bool) {
	values := sl.values
	println(len(values))
	if index < 0 || index >= len(values) {
		var t T
		return t, false
	}
	return values[index], true
}
