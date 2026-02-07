package pocket

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

// AssertNotNil asserts that the given value is not nil.
func AssertNotNil(t *testing.T, got any) {
	t.Helper()
	if isNil(got) {
		t.Errorf("expected non-nil, got nil")
	}
}

// AssertNil asserts that the given value is nil.
func AssertNil(t *testing.T, got any) {
	t.Helper()
	if !isNil(got) {
		t.Errorf("expected nil, got %v", got)
	}
}

// AssertTrue asserts that the given value is true.
func AssertTrue(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Errorf("expected true, got false")
	}
}

// AssertFalse asserts that the given value is false.
func AssertFalse(t *testing.T, got bool) {
	t.Helper()
	if got {
		t.Errorf("expected false, got true")
	}
}

// AssertEqual asserts that the given values are equal.
// It uses reflection to do a deep comparison.
func AssertEqual[T any](t *testing.T, a T, b T) {
	t.Helper()
	if !isEqual(a, b) {
		t.Errorf("expected values to equal, but %v does not equal %v", a, b)
	}
}

// AssertNotEqual asserts that the given values are not equal.
// It uses reflection to do a deep comparison.
func AssertNotEqual[T any](t *testing.T, a T, b T) {
	t.Helper()
	if isEqual(a, b) {
		t.Errorf("expected values not to equal, but got %v and %v", a, b)
	}
}

// AssertErrorIs asserts that the given error is of the given type.
// It uses the errors.Is to do the comparison, checking for wrapped errors.
func AssertErrorIs(t *testing.T, got error, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Errorf("expected error '%v' to be '%v'", got, want)
	}
}

// AssertContains asserts that the given string contains the given substring.
func AssertContains(t *testing.T, got string, substr string) {
	t.Helper()
	if !strings.Contains(got, substr) {
		t.Errorf("%q does not include the substring %q", got, substr)
	}
}

// AssertPanics asserts that the given function panics.
func AssertPanics(t *testing.T, f func()) {
	t.Helper()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but function did not panic")
			return
		}
	}()

	f()
}

func isEqual[T any](got T, want T) bool {
	if isNil(got) && isNil(want) {
		return true
	}
	return reflect.DeepEqual(got, want)
}

func isNil(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.Slice,
		reflect.UnsafePointer:
		return rv.IsNil()
	}

	return false
}
