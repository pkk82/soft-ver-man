package collections

import (
	"sort"
)

type SortedList[T comparable] struct {
	Values []T
	sort   func(i, j T) bool
}

func NewEmptySortedList[T comparable](sort func(i, j T) bool) SortedList[T] {
	return SortedList[T]{
		Values: make([]T, 0),
		sort:   sort,
	}
}

func NewSortedList[T comparable](values []T, sort func(i, j T) bool) SortedList[T] {
	result := NewEmptySortedList[T](sort)
	for _, v := range values {
		result.Insert(v)
	}
	return result
}

func (sl *SortedList[T]) Insert(value T) {
	values := sl.Values
	values = append(values, value)
	sort.Slice(values, func(i, j int) bool {
		return (*sl).sort(values[i], values[j])
	})
	sl.Values = values
}

func (sl *SortedList[T]) Delete(value T) bool {
	values := sl.Values
	for i, v := range values {
		if v == value {
			values = append(values[:i], values[i+1:]...)
			sl.Values = values
			return true
		}
	}
	return false
}

func (sl *SortedList[T]) Get(index int) (T, bool) {
	values := sl.Values
	if index < 0 || index >= len(values) {
		var t T
		return t, false
	}
	return values[index], true
}
