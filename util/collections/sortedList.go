/*
 * Copyright © 2024 Piotr Kozak <piotrkrzysztofkozak@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

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
