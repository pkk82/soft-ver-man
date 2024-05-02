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
