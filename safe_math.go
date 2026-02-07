package pocket

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
	switch any(a).(type) {

	case int, int8, int16, int32, int64:
		result := a + b

		// If two positives add to a negative, overflow occurred.
		if a > 0 && b > 0 && result < 0 {
			panic("integer overflow")
		}

		// If two negatives add to a positive, underflow occurred.
		if a < 0 && b < 0 && result > 0 {
			panic("integer underflow")
		}

		// Adding a negative to a positive can never overflow.
		return result

	case uint, uint8, uint16, uint32, uint64, uintptr:
		// Max value of T - all bits set to 1.
		maxUint := ^T(0)

		if a > maxUint-b {
			panic("unsigned integer overflow")
		}

		return a + b

	default:
		panic("invalid type for SafeAdd")
	}
}

// SafeSub returns the difference of two integers, panicking if the result overflows.
func SafeSub[T Int](a T, b T) T {
	switch any(a).(type) {

	case int, int8, int16, int32, int64:
		result := a - b

		// If a positive minus a negative gives a negative, overflow occurred.
		if a > 0 && b < 0 && result < 0 {
			panic("integer overflow")
		}

		// If a negative minus a positive gives a positive, underflow occurred.
		if a < 0 && b > 0 && result > 0 {
			panic("integer underflow")
		}

		return result

	case uint, uint8, uint16, uint32, uint64, uintptr:
		if a < b {
			panic("unsigned integer underflow")
		}

		return a - b

	default:
		panic("invalid type for SafeSub")
	}
}
