package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	require := require.New(t)
	out := &bytes.Buffer{}
	src := `
	package example

	option go_marshal = "unsafe"
	option go_int = 3

	const int32 a = 1
	const float b = 1.1
	const string c = "hello"
	const bool d = true
	const int64 e = -1
	const double f = 1.0
	const int64 g = 1
	message SearchRequest {
		string query = 1
		int32 page_number = 2
		int32 result_per_page = 3
	  
		message Foo {}
	  
		enum Bar {
		  FOO = 0
		}
	  }
	  
	  message SearchResponse {
		string results = 1
	  }
	`
	p, err := NewParser(out, src)
	require.NoError(err)
	require.NoError(p.Parse())
	fmt.Println(out.String())
}
