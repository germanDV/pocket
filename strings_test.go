package pocket

import (
	"strings"
	"testing"
)

func TestSafeCompare(t *testing.T) {
	type testCase struct {
		name   string
		a      string
		b      string
		expect bool
	}

	tests := []testCase{
		{
			name:   "same strings",
			a:      "token123",
			b:      "token123",
			expect: true,
		},
		{
			name:   "different strings",
			a:      "token123",
			b:      "token456",
			expect: false,
		},
		{
			name:   "different lengths",
			a:      "short",
			b:      "longer string",
			expect: false,
		},
		{
			name:   "empty strings",
			a:      "",
			b:      "",
			expect: true,
		},
		{
			name:   "one empty",
			a:      "token",
			b:      "",
			expect: false,
		},
		{
			name:   "unicode strings",
			a:      "héllo wörld",
			b:      "héllo wörld",
			expect: true,
		},
		{
			name:   "different unicode",
			a:      "héllo",
			b:      "hello",
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			AssertEqual(t, SafeCompare(tt.a, tt.b), tt.expect)
		})
	}
}

func TestGenerateString(t *testing.T) {
	t.Run("generates string of expected length", func(t *testing.T) {
		t.Parallel()
		s := GenerateString(32)
		AssertEqual(t, len(s) > 0, true)
		// base64 URL encoding expands the length
		AssertEqual(t, len(s) > 32, true)
	})

	t.Run("generates different strings", func(t *testing.T) {
		t.Parallel()
		s1 := GenerateString(32)
		s2 := GenerateString(32)
		AssertEqual(t, s1 == s2, false)
	})

	t.Run("generates valid base64 string", func(t *testing.T) {
		t.Parallel()
		s := GenerateString(16)
		// Should not contain standard base64 characters not in URL-safe set
		AssertEqual(t, strings.Contains(s, "+"), false)
		AssertEqual(t, strings.Contains(s, "/"), false)
	})

	t.Run("generates strings with different lengths", func(t *testing.T) {
		t.Parallel()
		s1 := GenerateString(8)
		s2 := GenerateString(64)
		AssertEqual(t, len(s1) < len(s2), true)
	})
}
