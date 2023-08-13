package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestParser(t *testing.T) {
	out := &bytes.Buffer{}
	src := `
	package example

	option go_marshal = "unsafe"
	option go_int = 3

	const int32 a = 1
	const float b = 1.1
	const string c = "hello"
	const bool d = true
	const int64 e = 1
	const double f = 1.0
	const int64 g = 1
	struct SearchRequest {
		string query
		int32 page_number
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
