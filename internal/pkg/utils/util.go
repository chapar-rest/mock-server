package utils

import "time"

// Now returns the current time in UTC, truncated to microseconds so that
// values survive a JSON or protobuf round trip unchanged.
func Now() time.Time {
	return time.Now().UTC().Truncate(time.Microsecond)
}

// Ptr is a convenience function making it possible to construct initialized
// pointer values in an expression.
func Ptr[T any](val T) *T {
	return &val
}

// Value attempts to dereference a pointer type, returning T's zero value if
// the pointer is nil.
func Value[T any](ptr *T) T {
	var val T
	if ptr != nil {
		val = *ptr
	}
	return val
}

// Must is a convenience function that panics if an error is not nil.
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// Optional is a generic type wrapper that lets us signal value presence or
// absence through an explicit property, including presence of a nil value.
type Optional[T any] struct {
	Val     T
	Present bool
}

// Some returns a populated Optional value.
func Some[T any](value T) Optional[T] {
	return Optional[T]{
		Val:     value,
		Present: true,
	}
}

// None returns an empty Optional value.
func None[T any]() Optional[T] {
	return Optional[T]{
		Val:     *new(T),
		Present: false,
	}
}

// PtrToOptional converts a pointer to an Optional; nil means absent.
func PtrToOptional[T any](val *T) Optional[T] {
	if val == nil {
		return None[T]()
	}
	return Some(*val)
}

// Get returns the value and whether it is present.
func (o Optional[T]) Get() (T, bool) {
	return o.Val, o.Present
}
