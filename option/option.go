package option

import (
	"bytes"
	"encoding/json"
	"reflect"

	"github.com/jkratz55/gonads"
)

var jsonNull = []byte("null")

// Option is a data type that represents a container that may or may not contain
// a value.
//
// An Option can be thought of in two states:
//
//	Some - Contains a value
//	None - Does not contain a value
//
// The zero-value isn't usable and an Option needs to be instantiated using one
// of the four factory methods.
//
//		Some - For creating an Option that contains a value
//		None - For creating an Option that is absent of a value
//	 FromNillable - Shorthand that will create Some or None depending on the value
//		PtrFromNillable - When an Option is needed where the value is a pointer
//
// Option supports JSON marshalling and unmarshalling out of the box. However, do
// to the way it is implemented `omitempty` will have no effect and won't prevent
// value from being encoded.
type Option[T any] struct {
	val    T
	exists bool
}

// Some creates an Option instance from a valid value.
func Some[T any](val T) Option[T] {
	// This is a protective guard to prevent misuse of the API.
	// If someone wanted to be a wise guy they could do something like the following:
	// 	opt := Some[error](nil)
	//	assert.True(t, opt.exists)
	//	assert.True(t, opt.IsSome())
	//	opt.IfSome(func(val error) {
	//		fmt.Println(val)
	//	})
	// The code would compile and run, but may have very strange results at runtime
	// and completely defeats the safety this API is trying to offer. One might argue
	// that the case above is someone abusing the API ... but probably better to stop
	// the misuse, at least for now until this code is determined as a performance
	// issue.
	if reflect.TypeOf(val) == nil {
		panic("cannot provide a nil value for Some")
	}
	return Option[T]{
		val:    val,
		exists: true,
	}
}

// None creates an Option instance that contains no value.
func None[T any]() Option[T] {
	return Option[T]{
		exists: false,
	}
}

// FromNillable creates an Option instance from a pointer to a value. If the ptr is
// non-nil will return an Option of Some, otherwise None.
func FromNillable[T any](val *T) Option[T] {
	if val == nil {
		return None[T]()
	}
	return Some(*val)
}

// PtrFromNillable creates an Option instance from a pointer to a value. If the ptr is
// non-nil will return an Option of Some, otherwise None.
//
// PtrFromNillable is similar to FromNillable but the Option type value is a ptr.
func PtrFromNillable[T any](val *T) Option[*T] {
	if val == nil {
		return None[*T]()
	}
	return Some[*T](val)
}

// IsSome returns a boolean indicating if the Option is Some.
func (o Option[T]) IsSome() bool {
	return o.exists
}

// IsNone returns a boolean indicating if the Option is None.
func (o Option[T]) IsNone() bool {
	return !o.exists
}

// IfSome invokes a Consumer func passing the value of the container to the
// Consumer if the Option is Some (contains a value).
func (o Option[T]) IfSome(fn gonads.Consumer[T]) {
	if o.exists {
		fn(o.val)
	}
}

// IfNone invokes the provided closure if the Option container does not contain
// a value (None).
func (o Option[T]) IfNone(fn func()) {
	if !o.exists {
		fn()
	}
}

// Filter returns None Option if the Option is already None. If the Option is
// Some (contains a value) the predicate is invoked. If the predicate returns
// true, returns an Option with the value. Otherwise, returns a None option.
func (o Option[T]) Filter(fn gonads.Predicate[T]) Option[T] {
	if !o.exists {
		return None[T]()
	}
	if fn(o.val) {
		return Some[T](o.val)
	}
	return None[T]()
}

// Get returns the value of the Option container along with a boolean indicating
// if the value is present.
//
// The Get method is similar to other methods and may be considered redundant, but
// it provides a more idiomatic Go approach that consumers may prefer.
func (o Option[T]) Get() (T, bool) {
	return o.val, o.exists
}

// Unwrap returns the value contained with Option, or panics if it doesn't exist.
//
// Because this function may panic, its use is generally discouraged. Instead, it
// is recommended to use UnwrapOrDefault or IfSome. Unwrap can only be safely called
// if IsSome is called and returns true beforehand.
func (o Option[T]) Unwrap() T {
	if !o.exists {
		panic("cannot unwrap none/nil value")
	}
	return o.val
}

// UnwrapOrDefault returns the value contained within Option, or if its None returns
// the default value provided.
func (o Option[T]) UnwrapOrDefault(defaultVal T) T {
	if !o.exists {
		return defaultVal
	}
	return o.val
}

// UnwrapOrElse returns the value contained within the Option or if its None
// executes the provided closure.
func (o Option[T]) UnwrapOrElse(fn gonads.Supplier[T]) T {
	if !o.exists {
		return fn()
	}
	return o.val
}

// Expect unwraps and returns the value contained within Option. It is similar
// in function to Unwrap. It will panic if the Option is None with the message
// provided.
//
// Expect can be useful for use cases where you want to panic because a required
// property, config, value, etc was not present.
func (o Option[T]) Expect(msg string) T {
	if !o.exists {
		panic(msg)
	}
	return o.val
}

// MarshalJSON marshals the Option type to JSON representation.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if !o.exists {
		return json.Marshal(nil)
	}

	b, err := json.Marshal(o.val)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// UnmarshalJSON unmarshalls JSON representation of Option to the Option type.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if len(data) <= 0 || bytes.Equal(data, jsonNull) {
		*o = None[T]()
		return nil
	}

	var v T
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*o = Some(v)
	return nil
}

// Map converts an Option[T] -> Option[R] by invoking the mapper function. If
// the given option is None, then None is returned.
func Map[T, R any](opt Option[T], fn gonads.Function[T, R]) Option[R] {
	if !opt.exists {
		return None[R]()
	}
	return Some(fn(opt.val))
}

// MapOr converts an Option[T] -> Option[R] by invoking the mapper function. If
// the given option is None, the provided fallback is returned.
func MapOr[T, R any](opt Option[T], fallback R, fn gonads.Function[T, R]) R {
	if !opt.exists {
		return fallback
	}
	return fn(opt.val)
}

// FlatMap converts an Option[T] -> Option[R] by invoking the mapper function. FlatMap
// differs from Map in the mapper function returns an Option[R] instead of a value. If
// the given Option is None, then None is returned.
func FlatMap[T, R any](opt Option[T], fn func(T) Option[R]) Option[R] {
	if !opt.exists {
		return None[R]()
	}
	return fn(opt.val)
}

// FlatMapOr converts an Option[T] -> Option[R] by invoking the mapper function. FlatMap
// differs from Map in the mapper function returns an Option[R] instead of a value. If
// the given Option is None, the default fallback value is returned.
func FlatMapOr[T, R any](opt Option[T], fallback R, fn func(T) Option[R]) Option[R] {
	if !opt.exists {
		return Some(fallback)
	}
	return fn(opt.val)
}
