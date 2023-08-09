package ast

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackage(t *testing.T) {
	require := require.New(t)
	data, err := basicParser.ParseString("", `
	package example
	`)
	require.NoError(err)
	require.Equal("example", data.Package.Name)
}
