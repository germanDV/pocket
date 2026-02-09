package pocket

import (
	"fmt"
	"math"
)

// Signed is a type constraint that matches all signed integer types.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a type constraint that matches all unsigned integer types.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Int is a type constraint that matches all integer types (signed and unsigned).
type Int interface {
	Signed | Unsigned
}

// SafeAdd returns the sum of two integers, panicking if the result overflows.
func SafeAdd[T Int](a T, b T) T {
	result, err := TrySafeAdd(a, b)
	if err != nil {
		panic(err)
	}
	return result
}

// TrySafeAdd returns the sum of two integers.
// Returns an error if the result would overflow or underflow.
func TrySafeAdd[T Int](a T, b T) (T, error) {
	var zero T

	switch any(a).(type) {
	case int, int8, int16, int32, int64:
		result := a + b

		// If two positives add to a negative, overflow occurred.
		if a > 0 && b > 0 && result < 0 {
			return zero, fmt.Errorf("integer overflow: %v + %v", a, b)
		}

		// If two negatives add to a positive, underflow occurred.
		if a < 0 && b < 0 && result > 0 {
			return zero, fmt.Errorf("integer underflow: %v + %v", a, b)
		}

		// Adding a negative to a positive can never overflow.
		return result, nil

	case uint, uint8, uint16, uint32, uint64, uintptr:
		// Max value of T - all bits set to 1.
		maxUint := ^T(0)

		if a > maxUint-b {
			return zero, fmt.Errorf("unsigned integer overflow: %v + %v", a, b)
		}

		return a + b, nil

	default:
		return zero, fmt.Errorf("invalid type for TrySafeAdd: %T", a)
	}
}

// SafeSub returns the difference of two integers, panicking if the result overflows.
func SafeSub[T Int](a T, b T) T {
	result, err := TrySafeSub(a, b)
	if err != nil {
		panic(err)
	}
	return result
}

// TrySafeSub returns the difference of two integers.
// Returns an error if the result would overflow or underflow.
func TrySafeSub[T Int](a T, b T) (T, error) {
	var zero T

	switch any(a).(type) {
	case int, int8, int16, int32, int64:
		result := a - b

		// If a positive minus a negative gives a negative, overflow occurred.
		if a > 0 && b < 0 && result < 0 {
			return zero, fmt.Errorf("integer overflow: %v - %v", a, b)
		}

		// If a negative minus a positive gives a positive, underflow occurred.
		if a < 0 && b > 0 && result > 0 {
			return zero, fmt.Errorf("integer underflow: %v - %v", a, b)
		}

		return result, nil

	case uint, uint8, uint16, uint32, uint64, uintptr:
		if a < b {
			return zero, fmt.Errorf("unsigned integer underflow: %v - %v", a, b)
		}

		return a - b, nil

	default:
		return zero, fmt.Errorf("invalid type for TrySafeSub: %T", a)
	}
}

// SafeMul returns the product of two integers, panicking if the result overflows.
func SafeMul[T Int](a T, b T) T {
	result, err := TrySafeMul(a, b)
	if err != nil {
		panic(err)
	}
	return result
}

// TrySafeMul returns the product of two integers.
// Returns an error if the result would overflow.
func TrySafeMul[T Int](a T, b T) (T, error) {
	var zero T

	switch any(a).(type) {
	case int, int8, int16, int32, int64:
		result := a * b

		// If two positives multiply to a negative, overflow occurred.
		if a > 0 && b > 0 && result < 0 {
			return zero, fmt.Errorf("integer overflow: %v * %v", a, b)
		}

		// If two negatives multiply to a negative, overflow occurred.
		if a < 0 && b < 0 && result < 0 {
			return zero, fmt.Errorf("integer overflow: %v * %v", a, b)
		}

		// Positive × negative or negative × positive can overflow if result
		// has the wrong sign. The result should be negative.
		if (a > 0 && b < 0 && result > 0) || (a < 0 && b > 0 && result > 0) {
			return zero, fmt.Errorf("integer overflow: %v * %v", a, b)
		}

		return result, nil

	case uint, uint8, uint16, uint32, uint64, uintptr:
		if a == 0 || b == 0 {
			return 0, nil
		}

		maxUint := ^T(0)
		if a > maxUint/b {
			return zero, fmt.Errorf("unsigned integer overflow: %v * %v", a, b)
		}

		return a * b, nil

	default:
		return zero, fmt.Errorf("invalid type for TrySafeMul: %T", a)
	}
}

// SafeDiv return the division of two integers, panicking if the result overflows.
func SafeDiv[T Int](a T, b T) T {
	result, err := TrySafeDiv(a, b)
	if err != nil {
		panic(err)
	}
	return result
}

// TrySafeDiv returns the division of two integers.
// Returns an error if the result would overflow or if dividing by zero.
func TrySafeDiv[T Int](a T, b T) (T, error) {
	var zero T

	if b == 0 {
		return zero, fmt.Errorf("division by zero")
	}

	switch any(a).(type) {
	case int, int8, int16, int32, int64:
		// Signed integer division can overflow only if a is MinInt and b is -1.
		// In this case, the result would be MaxInt + 1, which is not representable.
		isOverflow := false
		switch v := any(a).(type) {
		case int:
			isOverflow = v == math.MinInt && any(b).(int) == -1
		case int8:
			isOverflow = v == math.MinInt8 && any(b).(int8) == -1
		case int16:
			isOverflow = v == math.MinInt16 && any(b).(int16) == -1
		case int32:
			isOverflow = v == math.MinInt32 && any(b).(int32) == -1
		case int64:
			isOverflow = v == math.MinInt64 && any(b).(int64) == -1
		}
		if isOverflow {
			return zero, fmt.Errorf("integer overflow: %v / %v", a, b)
		}
		return a / b, nil

	case uint, uint8, uint16, uint32, uint64, uintptr:
		return a / b, nil

	default:
		return zero, fmt.Errorf("invalid type for TrySafeDiv: %T", a)
	}
}
