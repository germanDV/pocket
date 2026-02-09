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

func TestTrySafeMul(t *testing.T) {
	type testCase[N Int] struct {
		name      string
		a         N
		b         N
		want      N
		wantError bool
	}

	t.Run("int", func(t *testing.T) {
		intTests := []testCase[int]{
			{
				name: "positive * positive",
				a:    3,
				b:    4,
				want: 12,
			},
			{
				name: "positive * negative",
				a:    3,
				b:    -4,
				want: -12,
			},
			{
				name: "negative * positive",
				a:    -3,
				b:    4,
				want: -12,
			},
			{
				name: "negative * negative",
				a:    -3,
				b:    -4,
				want: 12,
			},
			{
				name:      "int overflow: positive * positive",
				a:         math.MaxInt,
				b:         2,
				wantError: true,
			},
			{
				name:      "int overflow: negative * negative",
				a:         math.MinInt,
				b:         -1,
				wantError: true,
			},
			{
				name:      "int overflow: mixed signs",
				a:         math.MaxInt,
				b:         -2,
				wantError: true,
			},
			{
				name:      "int overflow: MinInt * -1",
				a:         math.MinInt,
				b:         -1,
				wantError: true,
			},
			{
				name: "multiply by zero",
				a:    math.MaxInt,
				b:    0,
				want: 0,
			},
			{
				name: "multiply by one",
				a:    42,
				b:    1,
				want: 42,
			},
		}
		for _, tt := range intTests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result, err := TrySafeMul(tt.a, tt.b)
				if tt.wantError {
					AssertNotNil(t, err)
				} else {
					AssertNil(t, err)
					AssertEqual(t, result, tt.want)
				}
			})
		}
	})

	t.Run("int8", func(t *testing.T) {
		int8Tests := []testCase[int8]{
			{
				name: "int8 * int8",
				a:    int8(3),
				b:    int8(4),
				want: int8(12),
			},
			{
				name: "int8 * -int8",
				a:    int8(3),
				b:    int8(-4),
				want: int8(-12),
			},
			{
				name: "-int8 * -int8",
				a:    int8(-3),
				b:    int8(-4),
				want: int8(12),
			},
			{
				name:      "int8 overflow: positive * positive",
				a:         int8(math.MaxInt8),
				b:         int8(2),
				wantError: true,
			},
			{
				name:      "int8 overflow: negative * negative",
				a:         int8(math.MinInt8),
				b:         int8(-1),
				wantError: true,
			},
			{
				name:      "int8 overflow: mixed signs",
				a:         int8(math.MaxInt8),
				b:         int8(-2),
				wantError: true,
			},
		}
		for _, tt := range int8Tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result, err := TrySafeMul(tt.a, tt.b)
				if tt.wantError {
					AssertNotNil(t, err)
				} else {
					AssertNil(t, err)
					AssertEqual(t, result, tt.want)
				}
			})
		}
	})

	t.Run("uint", func(t *testing.T) {
		uintTests := []testCase[uint]{
			{
				name: "uint * uint",
				a:    uint(3),
				b:    uint(4),
				want: uint(12),
			},
			{
				name:      "uint overflow",
				a:         math.MaxUint,
				b:         uint(2),
				wantError: true,
			},
			{
				name: "multiply by zero",
				a:    uint(42),
				b:    uint(0),
				want: 0,
			},
			{
				name: "multiply by one",
				a:    uint(42),
				b:    uint(1),
				want: uint(42),
			},
		}
		for _, tt := range uintTests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result, err := TrySafeMul(tt.a, tt.b)
				if tt.wantError {
					AssertNotNil(t, err)
				} else {
					AssertNil(t, err)
					AssertEqual(t, result, tt.want)
				}
			})
		}
	})
}

func TestTrySafeDiv(t *testing.T) {
	type testCase[T Int] struct {
		name      string
		a, b      T
		want      T
		wantError bool
	}

	t.Run("int", func(t *testing.T) {
		intTests := []testCase[int]{
			{
				name: "10 / 2",
				a:    10,
				b:    2,
				want: 5,
			},
			{
				name: "10 / -2",
				a:    10,
				b:    -2,
				want: -5,
			},
			{
				name:      "division by zero",
				a:         10,
				b:         0,
				wantError: true,
			},
			{
				name:      "int overflow (MinInt / -1)",
				a:         math.MinInt,
				b:         -1,
				wantError: true,
			},
		}
		for _, tt := range intTests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result, err := TrySafeDiv(tt.a, tt.b)
				if tt.wantError {
					AssertNotNil(t, err)
				} else {
					AssertNil(t, err)
					AssertEqual(t, result, tt.want)
				}
			})
		}
	})

	t.Run("uint", func(t *testing.T) {
		uintTests := []testCase[uint]{
			{
				name: "10 / 2",
				a:    uint(10),
				b:    uint(2),
				want: uint(5),
			},
			{
				name:      "division by zero",
				a:         uint(10),
				b:         uint(0),
				wantError: true,
			},
		}
		for _, tt := range uintTests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result, err := TrySafeDiv(tt.a, tt.b)
				if tt.wantError {
					AssertNotNil(t, err)
				} else {
					AssertNil(t, err)
					AssertEqual(t, result, tt.want)
				}
			})
		}
	})

	t.Run("int8", func(t *testing.T) {
		int8Tests := []testCase[int8]{
			{
				name:      "int8 overflow (MinInt8 / -1)",
				a:         math.MinInt8,
				b:         -1,
				wantError: true,
			},
		}
		for _, tt := range int8Tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				result, err := TrySafeDiv(tt.a, tt.b)
				if tt.wantError {
					AssertNotNil(t, err)
				} else {
					AssertNil(t, err)
					AssertEqual(t, result, tt.want)
				}
			})
		}
	})
}
