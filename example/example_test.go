package example

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func Equal[T comparable](t *testing.T, expected, actual T) {
	t.Helper()

	if expected != actual {
		t.Errorf("want: %v; got: %v", expected, actual)
	}
}

// ObjectsAreEqual determines if two objects are considered equal.
//
// This function does no assertion of any kind.
func ObjectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}
func NoError(t *testing.T, err error, msgAndArgs ...interface{}) bool {
	t.Helper()
	if err != nil {
		t.Errorf(fmt.Sprintf("Received unexpected error:\n%+v", err), msgAndArgs...)
		return false
	}
	return true
}
func TestSearchRequest(t *testing.T) {
	a := &SearchRequest{
		query:           "hello",
		page_number:     1,
		result_per_page: 10,
	}
	data, err := a.MarshalBinary()
	NoError(t, err)
	t.Log(data)
	b := &SearchRequest{}
	err = b.UnmarshalBinary(data)
	NoError(t, err)
	t.Log(b)
	Equal(t, a, b)
}

func TestSearchResponse(t *testing.T) {
	a := &SearchResponse{
		results: "hello",
	}
	data, err := a.MarshalBinary()
	NoError(t, err)
	t.Log(data)
	b := &SearchResponse{}
	err = b.UnmarshalBinary(data)
	NoError(t, err)
	t.Log(b)
	Equal(t, a, b)
}

func TestA(t *testing.T) {
	a := &A{
		Name:     "hello",
		BirthDay: 33,
		Phone:    []byte("123"),
		Siblings: 44,
		Spouse:   true,
		Money:    1444.12324,
		//Children: []string{"a", "b"},
	}
	data, err := a.MarshalBinary()
	t.Logf("%x", data)
	assert.NoError(t, err)
	b := &A{}
	err = b.UnmarshalBinary(data)
	assert.NoError(t, err)
	assert.Equal(t, a, b)
}

/*
BenchmarkSearchRequest-8   	 7282776	       141.0 ns/op	      48 B/op	       2 allocs/op 2023/8/13
*/
func BenchmarkSearchRequest(b *testing.B) {
	a := &SearchRequest{
		query:           "hello",
		page_number:     1,
		result_per_page: 10,
	}
	for i := 0; i < b.N; i++ {
		data, err := a.MarshalBinary()
		if err != nil {
			b.Fatal(err)
		}
		c := &SearchRequest{}
		err = c.UnmarshalBinary(data)
		if err != nil {
			b.Fatal(err)
		}
		assert.Equal(b, a, c)
	}
}
