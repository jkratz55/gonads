package option

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestOption_MarshalJSON(t *testing.T) {

	res, _ := json.Marshal(person{})
	fmt.Println(string(res))

	p := person{
		FirstName:  "Billy",
		MiddleName: None[string](),
		LastName:   "Bob",
		Gender:     Some("MALE"),
	}
	data, err := json.Marshal(p)
	assert.NoError(t, err)

	fmt.Println(string(data))
}

func TestOption_UnmarshalJSON(t *testing.T) {

	data := []byte("{\"firstName\":\"Billy\",\"middleName\":null,\"lastName\":\"Bob\",\"gender\":\"MALE\"}")
	var p person

	err := json.Unmarshal(data, &p)
	assert.NoError(t, err)
}
