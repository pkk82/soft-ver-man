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
	"reflect"
	"testing"
)

func TestGroupBy(t *testing.T) {
	classifier := func(s string) int { return len(s) }
	type args[T any, K comparable] struct {
		values     []T
		classifier func(T) K
	}
	type testCase[T any, K comparable] struct {
		name string
		args args[T, K]
		want map[K][]T
	}
	tests := []testCase[string, int]{
		{
			name: "empty",
			args: args[string, int]{values: []string{}, classifier: classifier},
			want: map[int][]string{},
		},
		{
			name: "not empty",
			args: args[string, int]{values: []string{"Hello", "World", "Goodbye", "Cruel", "World"}, classifier: classifier},
			want: map[int][]string{5: {"Hello", "World", "Cruel", "World"}, 7: {"Goodbye"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GroupBy(tt.args.values, tt.args.classifier); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroupByAndCollect(t *testing.T) {
	classifier := func(s string) int { return len(s) }
	collector := func(ss []string) int { return len(ss) }
	type args[T any, K comparable, S any] struct {
		values     []T
		classifier func(T) K
		collector  func([]T) S
	}
	type testCase[T any, K comparable, S any] struct {
		name string
		args args[T, K, S]
		want map[K]S
	}
	tests := []testCase[string, int, int]{
		{
			name: "empty",
			args: args[string, int, int]{values: []string{}, classifier: classifier, collector: collector},
			want: map[int]int{},
		},
		{
			name: "not empty",
			args: args[string, int, int]{values: []string{"Hello", "World", "Goodbye", "Cruel", "World"}, classifier: classifier, collector: collector},
			want: map[int]int{5: 4, 7: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GroupByAndCollect(tt.args.values, tt.args.classifier, tt.args.collector); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupByAndCollect() = %v, want %v", got, tt.want)
			}
		})
	}
}
