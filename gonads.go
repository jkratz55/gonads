package gonads

// Predicate represents a predicate (boolean-valued function) of one argument.
type Predicate[T any] func(T) bool

// Consumer represents an operation that accepts a single input argument and
// returns no result. Unlike most other functional interfaces, Consumer is
// expected to operate via side effects.
type Consumer[T any] func(val T)

// Supplier represents a supplier of results.
type Supplier[T any] func() T

// Function represents a function that accepts one argument and produces a result.
type Function[T, R any] func(val T) R
