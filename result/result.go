package result

import (
	"github.com/jkratz55/gonads"
	"github.com/jkratz55/gonads/option"
)

// Result is a type representing the result of an operation that can fail.
//
// A Result can be thought of in two states:
//
//	Ok - Operation succeeded and there is a result
//	Error - Operation failed
//
// The zero value isn't usable and Result needs to be instantiated using one of
// the factory methods: From, Ok, or Error.
type Result[T any] struct {
	val T
	err error
}

// From creates a Result from a value and an error value.
func From[T any](val T, err error) Result[T] {
	return Result[T]{
		val: val,
		err: err,
	}
}

// Ok creates a Result representing success
func Ok[T any](val T) Result[T] {
	return Result[T]{
		val: val,
		err: nil,
	}
}

// Error creates a Result representing a failure
func Error[T any](err error) Result[T] {
	var zeroVal T
	return Result[T]{
		val: zeroVal,
		err: err,
	}
}

// IsOk returns a boolean indicating if the result is success or not
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr returns a boolean indicating if the result failed or not
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// IfOk invokes the Consumer function if the Result was successful.
func (r Result[T]) IfOk(fn gonads.Consumer[T]) {
	if r.err == nil {
		fn(r.val)
	}
}

// IfError invokes the provided closure if the Result was a failure.
func (r Result[T]) IfError(fn func(err error)) {
	if r.err != nil {
		fn(r.err)
	}
}

// Ok converts the value of the Result into an Option. If the Result was a failure
// returns None. Otherwise, returns Some(T)
func (r Result[T]) Ok() option.Option[T] {
	if r.err == nil {
		return option.Some(r.val)
	}
	return option.None[T]()
}

// Error converts the error value of the Result into an Option. If the Result was
// successful returns None, otherwise returns Some(error).
func (r Result[T]) Error() option.Option[error] {
	if r.err != nil {
		return option.Some(r.err)
	}
	return option.None[error]()
}

// Get unwraps the Result in a more idiomatic Go way returning the resulting value
// and error.
func (r Result[T]) Get() (T, error) {
	return r.val, r.err
}

// Unwrap returns the resulting value of Result or panics if there was an error.
//
// Since this function may panic its use is generally discouraged. Instead, it is
// recommended to use UnwrapOrDefault, Ok, IfOk, or Get.
func (r Result[T]) Unwrap() T {
	if r.err != nil {
		panic("cannot unwrap Result when Error")
	}
	return r.val
}

// UnwrapOrDefault returns the resulting value of Result or returns the provided
// default value if the Result is an Error.
func (r Result[T]) UnwrapOrDefault(defaultVal T) T {
	if r.err != nil {
		return defaultVal
	}
	return r.val
}

// UnwrapOrElse returns the resulting value of Result or returns the value resulting
// from invoking the provided closure.
func (r Result[T]) UnwrapOrElse(fn gonads.Supplier[T]) T {
	if r.err != nil {
		return fn()
	}
	return r.val
}

// Expect unwraps the value of Result or panics if the Result contains an error.
//
// Expect can be useful for use cases where you want to panic because a required
// operation did not succeed.
func (r Result[T]) Expect(msg string) T {
	if r.err != nil {
		panic(msg)
	}
	return r.val
}

// Map maps a Result[T] -> Result[R] using the provided mapper function. If the Result
// contained an error, an Error is returned with the error value untouched.
func Map[T, R any](res Result[T], fn func(T) R) Result[R] {
	if res.err != nil {
		return Error[R](res.err)
	}
	return Ok(fn(res.val))
}
