package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackage(t *testing.T) {
	require := require.New(t)
	data, err := basicParser.ParseString("", `
	// This is a comment
// This is another comment
  package example
	
	`)
	require.NoError(err)
	require.Equal("// This is a comment// This is another comment", strings.Join(data.Comments, ""))
	require.Equal("example", data.Package)
}
