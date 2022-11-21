package result

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkratz55/gonads"
	"github.com/jkratz55/gonads/option"
)

func TestFrom(t *testing.T) {

	testErr := errors.New("test error")

	tests := []struct {
		name     string
		val      string
		err      error
		expected Result[string]
	}{
		{
			name: "Ok",
			val:  "Billy Bob",
			err:  nil,
			expected: Result[string]{
				val: "Billy Bob",
				err: nil,
			},
		},
		{
			name: "Err",
			val:  "",
			err:  testErr,
			expected: Result[string]{
				val: "",
				err: testErr,
			},
		},
	}

	for _, test := range tests {
		actual := From(test.val, test.err)
		assert.Equal(t, test.expected, actual)
	}
}

func TestOk(t *testing.T) {
	res := Ok("Billy Bob")
	assert.NoError(t, res.err)
	assert.Equal(t, "Billy Bob", res.val)
}

func TestError(t *testing.T) {
	testErr := errors.New("test error")
	res := Error[string](testErr)
	assert.Error(t, testErr, res.err)
	assert.Equal(t, "", res.val)
}

func TestResult_IsOk(t *testing.T) {
	testErr := errors.New("test error")
	tests := []struct {
		name     string
		res      Result[string]
		expected bool
	}{
		{
			name:     "Error",
			res:      Error[string](testErr),
			expected: false,
		},
		{
			name:     "Ok",
			res:      Ok("Billy Bob"),
			expected: true,
		},
	}

	for _, test := range tests {
		actual := test.res.IsOk()
		assert.Equal(t, test.expected, actual)
	}
}

func TestResult_IsErr(t *testing.T) {
	testErr := errors.New("test error")
	tests := []struct {
		name     string
		res      Result[string]
		expected bool
	}{
		{
			name:     "Error",
			res:      Error[string](testErr),
			expected: true,
		},
		{
			name:     "Ok",
			res:      Ok("Billy Bob"),
			expected: false,
		},
	}

	for _, test := range tests {
		actual := test.res.IsErr()
		assert.Equal(t, test.expected, actual)
	}
}

func TestResult_IfOk(t *testing.T) {
	called := false
	testErr := errors.New("test error")
	tests := []struct {
		name       string
		res        Result[string]
		fn         gonads.Consumer[string]
		shouldCall bool
	}{
		{
			name: "Ok",
			res:  Ok("Billy Bob"),
			fn: func(val string) {
				called = true
			},
			shouldCall: true,
		},
		{
			name: "Error",
			res:  Error[string](testErr),
			fn: func(val string) {
				called = true
			},
			shouldCall: false,
		},
	}

	for _, test := range tests {
		called = false
		test.res.IfOk(test.fn)
		assert.Equal(t, test.shouldCall, called)
	}
}

func TestResult_IfError(t *testing.T) {
	called := false
	testErr := errors.New("test error")
	tests := []struct {
		name       string
		res        Result[string]
		fn         func(err error)
		shouldCall bool
	}{
		{
			name: "Ok",
			res:  Ok("Billy Bob"),
			fn: func(err error) {
				called = true
			},
			shouldCall: false,
		},
		{
			name: "Error",
			res:  Error[string](testErr),
			fn: func(err error) {
				called = true
			},
			shouldCall: true,
		},
	}

	for _, test := range tests {
		called = false
		test.res.IfError(test.fn)
		assert.Equal(t, test.shouldCall, called)
	}
}

func TestResult_Ok(t *testing.T) {
	testErr := errors.New("test error")
	tests := []struct {
		name     string
		res      Result[string]
		expected option.Option[string]
	}{
		{
			name:     "Ok",
			res:      Ok("Billy Bob"),
			expected: option.Some("Billy Bob"),
		},
		{
			name:     "Error",
			res:      Error[string](testErr),
			expected: option.None[string](),
		},
	}

	for _, test := range tests {
		actual := test.res.Ok()
		assert.Equal(t, test.expected, actual)
	}
}

func TestResult_Error(t *testing.T) {
	testErr := errors.New("test error")
	tests := []struct {
		name     string
		res      Result[string]
		expected option.Option[error]
	}{
		{
			name:     "Ok",
			res:      Ok("Billy Bob"),
			expected: option.None[error](),
		},
		{
			name:     "Error",
			res:      Error[string](testErr),
			expected: option.Some(testErr),
		},
	}

	for _, test := range tests {
		actual := test.res.Error()
		assert.Equal(t, test.expected, actual)
	}
}

func TestResult_Get(t *testing.T) {
	res := Result[string]{
		val: "Billy Bob",
		err: nil,
	}
	assert.NoError(t, res.err)
	assert.Equal(t, "Billy Bob", res.val)
}

func TestResult_Unwrap(t *testing.T) {
	testErr := errors.New("test error")
	ok := Ok("Billy Bob")
	assert.NotPanics(t, func() {
		val := ok.Unwrap()
		assert.Equal(t, "Billy Bob", val)
	})

	err := Error[string](testErr)
	assert.Panics(t, func() {
		err.Unwrap()
	})
}

func TestResult_UnwrapOrDefault(t *testing.T) {
	tests := []struct {
		name       string
		res        Result[string]
		defaultVal string
		expected   string
	}{
		{
			name:       "Ok",
			res:        Ok("Billy Bob"),
			defaultVal: "Silly Jilly",
			expected:   "Billy Bob",
		},
		{
			name:       "Error",
			res:        Error[string](errors.New("errrrrrrrrr")),
			defaultVal: "Silly Jilly",
			expected:   "Silly Jilly",
		},
	}

	for _, test := range tests {
		actual := test.res.UnwrapOrDefault(test.defaultVal)
		assert.Equal(t, test.expected, actual)
	}
}

func TestResult_UnwrapOrElse(t *testing.T) {
	tests := []struct {
		name     string
		res      Result[string]
		fn       gonads.Supplier[string]
		expected string
	}{
		{
			name: "Ok",
			res:  Ok("Billy Bob"),
			fn: func() string {
				return "Silly Jilly"
			},
			expected: "Billy Bob",
		},
		{
			name: "Error",
			res:  Error[string](errors.New("Errrrrrrrrr")),
			fn: func() string {
				return "Silly Jilly"
			},
			expected: "Silly Jilly",
		},
	}

	for _, test := range tests {
		actual := test.res.UnwrapOrElse(test.fn)
		assert.Equal(t, test.expected, actual)
	}
}

func TestResult_Expect(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic!")
		} else {
			msg := r.(string)
			assert.Equal(t, "critical operation failed", msg)
		}
	}()

	res := Error[string](errors.New("errrrrrrr"))
	res.Expect("critical operation failed")
}

func TestMap(t *testing.T) {
	ok := Ok(10)
	res := Map(ok, func(val int) int {
		return val * 2
	})
	assert.Equal(t, 20, res.val)

	err := Error[int](errors.New("err"))
	res = Map(err, func(val int) int {
		return val * 2
	})
	assert.Error(t, res.err)
}
