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

	option go_package = "api"
	option go_int = 3

	const int32 a = 1
	`
	p, err := NewParser(out, src)
	require.NoError(err)
	require.NoError(p.Parse())
	fmt.Println(out.String())
}
