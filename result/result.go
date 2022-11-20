package result

import (
	"github.com/jkratz55/gonads"
	"github.com/jkratz55/gonads/option"
)

type Result[T any] struct {
	val T
	err error
}

func From[T any](val T, err error) Result[T] {
	return Result[T]{
		val: val,
		err: err,
	}
}

func Ok[T any](val T) Result[T] {
	return Result[T]{
		val: val,
		err: nil,
	}
}

func Error[T any](err error) Result[T] {
	var zeroVal T
	return Result[T]{
		val: zeroVal,
		err: err,
	}
}

func (r Result[T]) IsOk() bool {
	return r.err == nil
}

func (r Result[T]) IsErr() bool {
	return r.err != nil
}

func (r Result[T]) IfOk(fn gonads.Consumer[T]) {
	if r.err == nil {
		fn(r.val)
	}
}

func (r Result[T]) IfError(fn func(err error)) {
	if r.err != nil {
		fn(r.err)
	}
}

func (r Result[T]) Ok() option.Option[T] {
	if r.err != nil {
		return option.Some(r.val)
	}
	return option.None[T]()
}

func (r Result[T]) Error() option.Option[error] {
	if r.err != nil {
		return option.Some(r.err)
	}
	return option.None[error]()
}

func (r Result[T]) Get() (T, error) {
	return r.val, r.err
}

func (r Result[T]) Unwrap() T {
	if r.err != nil {
		panic("cannot unwrap Result when Error")
	}
	return r.val
}

func (r Result[T]) UnwrapOrDefault(defaultVal T) T {
	if r.err != nil {
		return defaultVal
	}
	return r.val
}

func (r Result[T]) UnwrapOrElse(fn gonads.Supplier[T]) T {
	if r.err != nil {
		return fn()
	}
	return r.val
}

func (r Result[T]) Expect(msg string) T {
	if r.err != nil {
		panic(msg)
	}
	return r.val
}

func Map[T, R any](res Result[T], fn func(T) R) Result[R] {
	if res.err != nil {
		return Error[R](res.err)
	}
	return Ok(fn(res.val))
}
