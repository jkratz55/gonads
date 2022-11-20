package result

type Result[T any] struct {
	val *T
	err error
}

func (r Result[T]) IsOk() bool {
	return r.err == nil
}

func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// func (r Result[T]) Map(fn func())

func (r Result[T]) Unwrap() T {
	if r.err != nil {
		panic("")
	}
	return *r.val
}

func From[T any](val T, err error) Result[T] {
	return Result[T]{
		val: &val,
		err: err,
	}
}

func Ok[T any](val T) Result[T] {
	return Result[T]{
		val: &val,
		err: nil,
	}
}

func Err[T any](err error) Result[T] {
	return Result[T]{
		val: nil,
		err: err,
	}
}
