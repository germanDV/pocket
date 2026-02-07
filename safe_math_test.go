package pocket

import (
	"math"
	"testing"
)

func TestSafeAdd(t *testing.T) {
	type testCase[N Int] struct {
		name        string
		a           N
		b           N
		want        N
		shouldPanic bool
	}

	intTests := []testCase[int]{
		{
			name: "positive + positive",
			a:    1,
			b:    2,
			want: 3,
		},
		{
			name: "positive + negative",
			a:    1,
			b:    -2,
			want: -1,
		},
		{
			name: "negative + positive",
			a:    -1,
			b:    2,
			want: 1,
		},
		{
			name: "negative + negative",
			a:    -1,
			b:    -2,
			want: -3,
		},
		{
			name:        "int overflow",
			a:           math.MaxInt,
			b:           1,
			shouldPanic: true,
		},
	}
	for _, tt := range intTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.shouldPanic {
				AssertPanics(t, func() { SafeAdd(tt.a, tt.b) })
			} else {
				AssertEqual(t, SafeAdd(tt.a, tt.b), tt.want)
			}
		})
	}

	int8Tests := []testCase[int8]{
		{
			name: "int8 + int8",
			a:    int8(23),
			b:    int8(45),
			want: int8(68),
		},
		{
			name: "int8 + -int8",
			a:    int8(23),
			b:    int8(-45),
			want: int8(-22),
		},
		{
			name:        "int8 overflow",
			a:           int8(math.MaxInt8),
			b:           int8(1),
			shouldPanic: true,
		},
	}
	for _, tt := range int8Tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.shouldPanic {
				AssertPanics(t, func() { SafeAdd(tt.a, tt.b) })
			} else {
				AssertEqual(t, SafeAdd(tt.a, tt.b), tt.want)
			}
		})
	}

	uintTests := []testCase[uint]{
		{
			name: "uint + uint",
			a:    uint(23),
			b:    uint(45),
			want: uint(68),
		},
		{
			name:        "uint overflow",
			a:           math.MaxUint,
			b:           uint(1),
			shouldPanic: true,
		},
	}
	for _, tt := range uintTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.shouldPanic {
				AssertPanics(t, func() { SafeAdd(tt.a, tt.b) })
			} else {
				AssertEqual(t, SafeAdd(tt.a, tt.b), tt.want)
			}
		})
	}
}

func TestSafeSub(t *testing.T) {
	type testCase[N Int] struct {
		name        string
		a           N
		b           N
		want        N
		shouldPanic bool
	}

	intTests := []testCase[int]{
		{
			name: "positive - positive",
			a:    1,
			b:    2,
			want: -1,
		},
		{
			name: "positive - negative",
			a:    1,
			b:    -2,
			want: 3,
		},
		{
			name: "negative - positive",
			a:    -1,
			b:    2,
			want: -3,
		},
		{
			name: "negative - negative",
			a:    -1,
			b:    -2,
			want: 1,
		},
		{
			name:        "int underflow",
			a:           math.MinInt,
			b:           1,
			shouldPanic: true,
		},
	}
	for _, tt := range intTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.shouldPanic {
				AssertPanics(t, func() { SafeSub(tt.a, tt.b) })
			} else {
				AssertEqual(t, SafeSub(tt.a, tt.b), tt.want)
			}
		})
	}

	int8Tests := []testCase[int8]{
		{
			name: "int8 - int8",
			a:    int8(23),
			b:    int8(45),
			want: int8(-22),
		},
		{
			name: "int8 - -int8",
			a:    int8(23),
			b:    int8(-45),
			want: int8(68),
		},
		{
			name:        "int8 underflow",
			a:           int8(math.MinInt8),
			b:           int8(1),
			shouldPanic: true,
		},
	}
	for _, tt := range int8Tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.shouldPanic {
				AssertPanics(t, func() { SafeSub(tt.a, tt.b) })
			} else {
				AssertEqual(t, SafeSub(tt.a, tt.b), tt.want)
			}
		})
	}

	uintTests := []testCase[uint]{
		{
			name: "uint - uint",
			a:    uint(45),
			b:    uint(23),
			want: uint(22),
		},
		{
			name:        "uint underflow",
			a:           uint(0),
			b:           uint(1),
			shouldPanic: true,
		},
	}
	for _, tt := range uintTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.shouldPanic {
				AssertPanics(t, func() { SafeSub(tt.a, tt.b) })
			} else {
				AssertEqual(t, SafeSub(tt.a, tt.b), tt.want)
			}
		})
	}
}
