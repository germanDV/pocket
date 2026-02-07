package pocket

import (
	"strings"
	"testing"
)

func TestMap(t *testing.T) {
	type testCase[T any, U any] struct {
		name  string
		slice []T
		f     func(T) U
		want  []U
	}

	intTests := []testCase[int, int]{
		{
			name:  "int -> int",
			slice: []int{1, 2, 3},
			f: func(i int) int {
				return i * 2
			},
			want: []int{2, 4, 6},
		},
	}
	for _, tt := range intTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			AssertEqual(t, Map(tt.slice, tt.f), tt.want)
		})
	}

	stringTests := []testCase[string, string]{
		{
			name:  "string -> string",
			slice: []string{"abc", "def", "", "xyz"},
			f: func(s string) string {
				return strings.ToUpper(s)
			},
			want: []string{"ABC", "DEF", "", "XYZ"},
		},
	}
	for _, tt := range stringTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			AssertEqual(t, Map(tt.slice, tt.f), tt.want)
		})
	}
}

func TestFilter(t *testing.T) {
	type testCase[T any] struct {
		name  string
		slice []T
		f     func(T) bool
		want  []T
	}

	intTests := []testCase[int]{
		{
			name:  "int -> bool",
			slice: []int{1, 2, 3},
			f: func(i int) bool {
				return i%2 == 0
			},
			want: []int{2},
		},
	}
	for _, tt := range intTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			AssertEqual(t, Filter(tt.slice, tt.f), tt.want)
		})
	}
}
