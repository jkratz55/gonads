package option

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkratz55/gonads"
)

type person struct {
	FirstName  string         `json:"firstName,omitempty"`
	MiddleName Option[string] `json:"middleName,omitempty"`
	LastName   string         `json:"lastName,omitempty"`
	Gender     Option[string] `json:"gender,omitempty"`
}

func (p person) String() string {
	return "I am a string"
}

func TestSome_PreventNil(t *testing.T) {
	assert.Panics(t, func() {
		_ = Some[error](nil)
	})
}

func TestSome(t *testing.T) {
	var opt Option[string]
	assert.NotPanics(t, func() {
		opt = Some("Billy Bob")
	})
	assert.True(t, opt.exists)
	assert.Equal(t, "Billy Bob", opt.val)

	var opt2 Option[fmt.Stringer]
	var p fmt.Stringer
	p = person{
		FirstName:  "Billy",
		MiddleName: Some("Jane"),
		LastName:   "Bob",
		Gender:     Some("MALE"),
	}
	assert.NotPanics(t, func() {
		opt2 = Some(p)
	})
	assert.True(t, opt2.exists)
}

func TestNone(t *testing.T) {
	var opt Option[string]
	assert.NotPanics(t, func() {
		opt = None[string]()
	})
	assert.False(t, opt.exists)
}

func TestFromNillable(t *testing.T) {
	var p *person
	opt := FromNillable(p)
	assert.False(t, opt.exists)

	var realPerson person
	opt = FromNillable(&realPerson)
	assert.True(t, opt.exists)
	assert.Equal(t, realPerson, opt.val)
}

func TestPtrFromNillable(t *testing.T) {

	p := new(person)
	p.FirstName = "Billy"
	p.LastName = "Bob"

	opt := PtrFromNillable(p)
	assert.True(t, opt.exists)
	assert.Equal(t, p, opt.val)
}

func TestOption_IsSome(t *testing.T) {
	opt := Some("Billy Bob")
	assert.True(t, opt.exists)
	assert.True(t, opt.IsSome())
}

func TestOption_IsNone(t *testing.T) {
	opt := None[string]()
	assert.False(t, opt.exists)
	assert.True(t, opt.IsNone())
}

func TestOption_IfSome(t *testing.T) {
	called := false
	opt := Some("Billy Bob")
	opt.IfSome(func(val string) {
		assert.Equal(t, "Billy Bob", val)
		called = true
	})
	assert.True(t, called)

	called = false
	opt2 := None[string]()
	opt2.IfSome(func(val string) {
		called = true
	})
	assert.False(t, called)
}

func TestOption_IfNone(t *testing.T) {
	called := false
	opt := None[string]()
	opt.IfNone(func() {
		called = true
	})
	assert.True(t, called)

	called = false
	opt = Some("Billy")
	opt.IfNone(func() {
		called = true
	})
	assert.False(t, called)
}

func TestOption_Filter(t *testing.T) {

	tests := []struct {
		name     string
		opt      Option[string]
		pred     gonads.Predicate[string]
		expected Option[string]
	}{
		{
			name: "None Option",
			opt:  None[string](),
			pred: func(val string) bool {
				return strings.Contains(val, "Billy")
			},
			expected: None[string](),
		},
		{
			name: "Some But Doesn't Meet Predicate",
			opt:  Some("Joe Joe"),
			pred: func(val string) bool {
				return strings.Contains(val, "Billy")
			},
			expected: None[string](),
		},
		{
			name: "Some And Meets Predicate",
			opt:  Some("Billy Bob"),
			pred: func(val string) bool {
				return strings.Contains(val, "Billy")
			},
			expected: Some("Billy Bob"),
		},
	}

	for _, test := range tests {
		actual := test.opt.Filter(test.pred)
		assert.Equal(t, test.expected, actual, fmt.Sprintf("Test %s failed", test.name))
	}
}

func TestOption_Get(t *testing.T) {
	tests := []struct {
		name           string
		opt            Option[string]
		expectedValue  string
		expectedExists bool
	}{
		{
			name:           "None",
			opt:            None[string](),
			expectedValue:  "",
			expectedExists: false,
		},
		{
			name:           "Some",
			opt:            Some("Billy Bob"),
			expectedValue:  "Billy Bob",
			expectedExists: true,
		},
	}

	for _, test := range tests {
		actualVal, actualExists := test.opt.Get()
		assert.Equal(t, test.expectedValue, actualVal)
		assert.Equal(t, test.expectedExists, actualExists)
	}
}

func TestOption_Unwrap(t *testing.T) {
	tests := []struct {
		name        string
		shouldPanic bool
		expected    string
		opt         Option[string]
	}{
		{
			name:        "None - Panic",
			shouldPanic: true,
			expected:    "",
			opt:         None[string](),
		},
		{
			name:        "Some - No Panic",
			shouldPanic: false,
			expected:    "Billy Bob",
			opt:         Some("Billy Bob"),
		},
	}

	for _, test := range tests {
		if test.shouldPanic {
			assert.Panics(t, func() {
				test.opt.Unwrap()
			})
		} else {
			assert.NotPanics(t, func() {
				actual := test.opt.Unwrap()
				assert.Equal(t, test.expected, actual)
			})
		}
	}
}

func TestOption_UnwrapOrDefault(t *testing.T) {
	tests := []struct {
		name       string
		opt        Option[string]
		defaultVal string
		expected   string
	}{
		{
			name:       "None - Default",
			opt:        None[string](),
			defaultVal: "Jilly Bob",
			expected:   "Jilly Bob",
		},
		{
			name:       "Some - Opt Value",
			opt:        Some("Billy Bob"),
			defaultVal: "Jilly Bob",
			expected:   "Billy Bob",
		},
	}

	for _, test := range tests {
		actual := test.opt.UnwrapOrDefault(test.defaultVal)
		assert.Equal(t, test.expected, actual)
	}
}

func TestOption_UnwrapOrElse(t *testing.T) {
	tests := []struct {
		name     string
		opt      Option[string]
		fn       gonads.Supplier[string]
		expected string
	}{
		{
			name: "None",
			opt:  None[string](),
			fn: func() string {
				return "Jill"
			},
			expected: "Jill",
		},
		{
			name: "Some",
			opt:  Some("Billy Bob"),
			fn: func() string {
				return "Jill"
			},
			expected: "Billy Bob",
		},
	}

	for _, test := range tests {
		actual := test.opt.UnwrapOrElse(test.fn)
		assert.Equal(t, test.expected, actual)
	}
}

func TestOption_Expect_Panic(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			msg := r.(string)
			assert.Equal(t, "missing string property", msg)
		}
	}()

	opt := None[string]()
	opt.Expect("missing string property")
}

func TestOption_Expect(t *testing.T) {
	opt := Some("Billy Bob")
	assert.Equal(t, "Billy Bob", opt.Expect("oppps missing value"))
}

func TestOption_MarshalJSON(t *testing.T) {

	p := person{
		FirstName:  "Billy",
		MiddleName: None[string](),
		LastName:   "Bob",
		Gender:     Some("MALE"),
	}
	data, err := json.Marshal(p)
	assert.NoError(t, err)
	assert.Equal(t, []byte("{\"firstName\":\"Billy\",\"middleName\":null,\"lastName\":\"Bob\",\"gender\":\"MALE\"}"), data)
}

func TestOption_UnmarshalJSON(t *testing.T) {

	data := []byte("{\"firstName\":\"Billy\",\"middleName\":null,\"lastName\":\"Bob\",\"gender\":\"MALE\"}")
	var p person

	err := json.Unmarshal(data, &p)
	assert.NoError(t, err)
	assert.Equal(t, person{
		FirstName:  "Billy",
		MiddleName: None[string](),
		LastName:   "Bob",
		Gender:     Some("MALE"),
	}, p)
}
