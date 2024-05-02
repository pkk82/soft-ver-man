package collections

import (
	"golang.org/x/exp/constraints"
	"reflect"
	"testing"
)

func TestSortedList_Get(t *testing.T) {
	type args struct {
		index int
	}
	type Want[T constraints.Ordered] struct {
		value T
		found bool
	}
	type testCase[T constraints.Ordered] struct {
		name string
		sl   SortedList[T]
		args args
		want Want[T]
	}
	tests := []testCase[string]{
		{
			name: "empty list - not found",
			sl:   newList(),
			args: args{index: 0},
			want: Want[string]{
				value: "",
				found: false,
			},
		},
		{
			name: "list - first element",
			sl:   newList("a", "b", "c"),
			args: args{index: 0},
			want: Want[string]{
				value: "a",
				found: true,
			},
		},
		{
			name: "list - last element",
			sl:   newList("a", "b", "c"),
			args: args{index: 2},
			want: Want[string]{
				value: "c",
				found: true,
			},
		},
		{
			name: "list - not found",
			sl:   newList("a", "b", "c"),
			args: args{index: 3},
			want: Want[string]{
				value: "",
				found: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := tt.sl.Get(tt.args.index)
			if !reflect.DeepEqual(got, tt.want.value) {
				t.Errorf("Get().value got = %v, want.value %v", got, tt.want.value)
			}
			if found != tt.want.found {
				t.Errorf("Get().found got = %v, want.found %v", found, tt.want.found)
			}
		})
	}
}

func TestSortedList_Delete(t *testing.T) {
	type args[T constraints.Ordered] struct {
		value T
	}
	type Want[T constraints.Ordered] struct {
		value SortedList[T]
		found bool
	}
	type testCase[T constraints.Ordered] struct {
		name string
		init SortedList[T]
		args args[T]
		want Want[T]
	}
	tests := []testCase[string]{
		{
			name: "empty list",
			init: newList(),
			args: args[string]{value: "a"},
			want: Want[string]{
				value: newList(),
				found: false,
			},
		}, {
			name: " list - element found",
			init: newList("a"),
			args: args[string]{value: "a"},
			want: Want[string]{
				value: newList(),
				found: true,
			},
		}, {
			name: " list - element not found",
			init: newList("a", "b"),
			args: args[string]{value: "c"},
			want: Want[string]{
				value: newList("a", "b"),
				found: false,
			},
		}, {
			name: " list - element duplicated",
			init: newList("a", "b", "c", "a"),
			args: args[string]{value: "a"},
			want: Want[string]{
				value: newList("b", "c", "a"),
				found: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := tt.init.Delete(tt.args.value)
			if !reflect.DeepEqual(tt.init.values, tt.want.value.values) {
				t.Errorf("Delete() got = %v, want %v", tt.init, tt.want.value)
			}
			if found != tt.want.found {
				t.Errorf("Delete().found got = %v, want.found %v", found, tt.want.found)
			}
		})
	}
}

func TestSortedList_Insert(t *testing.T) {
	type args[T constraints.Ordered] struct {
		value T
	}
	type testCase[T constraints.Ordered] struct {
		name string
		init SortedList[T]
		args args[T]
		want SortedList[T]
	}
	tests := []testCase[string]{
		{
			name: "empty list",
			init: newList(),
			args: args[string]{value: "a"},
			want: newList("a"),
		}, {
			name: " list - same element",
			init: newList("a"),
			args: args[string]{value: "a"},
			want: newList("a", "a"),
		}, {
			name: " list - in the middle",
			init: newList("a", "c"),
			args: args[string]{value: "b"},
			want: newList("a", "b", "c"),
		},

		{
			name: " list - at the beginning",
			init: newList("b", "c"),
			args: args[string]{value: "a"},
			want: newList("a", "b", "c"),
		}, {
			name: " list - at the end",
			init: newList("b", "c"),
			args: args[string]{value: "d"},
			want: newList("b", "c", "d"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.init.Insert(tt.args.value)
			if !reflect.DeepEqual(tt.init.values, tt.want.values) {
				t.Errorf("Get() got = %v, want %v", tt.init, tt.want)
			}
		})
	}
}

func TestNewSortedList(t *testing.T) {
	type args[T constraints.Ordered] struct {
		values []T
	}
	type testCase[T constraints.Ordered] struct {
		name string
		args args[T]
		want SortedList[T]
	}
	tests := []testCase[string]{
		{
			name: "empty list",
			args: args[string]{values: []string{}},
			want: newList(),
		},
		{
			name: "list",
			args: args[string]{values: []string{"b", "a", "c"}},
			want: newList("a", "b", "c"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSortedList(tt.args.values, func(a, b string) bool { return a < b }); !reflect.DeepEqual(got.values, tt.want.values) {
				t.Errorf("NewSortedList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func newList(values ...string) SortedList[string] {
	sl := NewEmptySortedList[string](func(a, b string) bool { return a < b })
	for _, value := range values {
		sl.Insert(value)
	}
	return sl
}
