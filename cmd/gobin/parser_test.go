package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestParser(t *testing.T) {
	out := &bytes.Buffer{}
	src := `
	package example

	option go_marshal = "unsafe"
	option go_int = 3

	//comment int32 a = 1
	const int32 a = 1
	/* 
		comment 3
		coment 30 
	*/
	const float b = 1.1
	const string c = "hello"
	const bool d = true
	const int64 e = 1
	const double f = 1.0
	const int64 g = 1
	struct Person {
		string Name
		uint64 BirthDay
		bytes Phone
		int32 Siblings
		bool Spouse

		double Money
	}

	// SearchRequest is a request message for SearchService.Search.
	struct SearchRequest {
		// query is a search query string.
		string query
		Person person
		// page_number is a page number.
		int32 page_number
		// result_per_page is a result per page.
		int32 result_per_page
	  }
	  
	  struct SearchResponse {
		string results
	  }
	`
	p, err := NewParser(out, src)
	assert.NoError(t, err)
	assert.NoError(t, p.Parse())
	fmt.Println(out.String())
}

func TestPebbleTemple(t *testing.T) {
	input, err := os.ReadFile("./testdata/pebble.gobin")
	assert.NoError(t, err)

	out := &bytes.Buffer{}
	var opt []option
	opt = []option{WithFormatted()}
	p, err := NewParser(out, input, opt...)
	assert.NoError(t, err)
	assert.NoError(t, p.Parse())
	fmt.Println(out.String())
}
