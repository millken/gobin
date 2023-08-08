package main

import (
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/stretchr/testify/require"
)

func TestOption(t *testing.T) {
	require := require.New(t)
	parser, err := participle.Build[Grammar](
		participle.Lexer(def),
		participle.Unquote(),
		participle.Elide("Whitespace"),
	)
	require.NoError(err)
	data, err := parser.ParseString("", `
	package example

	option go_package = "api"
	option go_int = 3

	const int32 a = 1

	enum Type {
		INT = 0
		DOUBLE = 1
	  }
	`)
	_ = data
	require.NoError(err)
	require.Equal("example", data.Package)
	require.Equal(2, len(data.Option))
	require.Equal("go_package", data.Option[0].Name)
	require.Equal("api", *data.Option[0].Value.String)
	require.Equal("go_int", data.Option[1].Name)
	require.Equal(int64(3), *data.Option[1].Value.Int)
	require.Equal("a", data.Consts[0].Name)
	require.Equal("Int32", data.Consts[0].Type.Scalar.String())
	require.Equal(int(1), *data.Consts[0].Value.Int)
}

func TestConst(t *testing.T) {
	require := require.New(t)
	type Const struct {
		Pos   lexer.Position
		Type  *Type    `"const" @@`
		Name  string   `@Ident`
		Value *Literal `"=" @@`
	}
	type Grammar struct {
		Pos     lexer.Position
		Package string   ` "package" @(Ident ( "." Ident )*)`
		Consts  []*Const `@@*`
	}
	parser := participle.MustBuild[Grammar](
		participle.UseLookahead(2),
		participle.Unquote(),
	)
	data, err := parser.ParseString("", `
	package example
	const int32 a = 1
	const float b = 1.1
	const string c = "hello"
	const bool d = true
	const int64 e = -1
	const double f = 1.0
	const int64 g = 1
	`)

	require.NoError(err)
	_ = data
	require.Equal("example", data.Package)
	require.Equal(7, len(data.Consts))
	require.Equal("a", data.Consts[0].Name)
	require.Equal("Int32", data.Consts[0].Type.Scalar.String())
	require.Equal(int64(1), *data.Consts[0].Value.Int)
	require.Equal("b", data.Consts[1].Name)
	require.Equal("Float", data.Consts[1].Type.Scalar.String())
	require.Equal(float64(1.1), *data.Consts[1].Value.Float)
	require.Equal("c", data.Consts[2].Name)
	require.Equal("String", data.Consts[2].Type.Scalar.String())
	require.Equal("hello", *data.Consts[2].Value.String)
	require.Equal("d", data.Consts[3].Name)
	require.Equal("Bool", data.Consts[3].Type.Scalar.String())
	require.Equal("true", *data.Consts[3].Value.Bool)
	require.Equal("e", data.Consts[4].Name)
	require.Equal("Int64", data.Consts[4].Type.Scalar.String())
	require.Equal(int64(-1), *data.Consts[4].Value.Int)
	require.Equal("f", data.Consts[5].Name)
	require.Equal("Double", data.Consts[5].Type.Scalar.String())
	require.Equal(float64(1.0), *data.Consts[5].Value.Float)
}

func TestEnum(t *testing.T) {
	require := require.New(t)
	parser, err := participle.Build[Grammar](
		participle.Lexer(def),
		participle.Unquote(),
		participle.Elide("Whitespace"),
	)
	require.NoError(err)
	data, err := parser.ParseString("", `
	package example
	enum Type {
		INT = 0
		DOUBLE = 1
	}
	enum Type2 {
		INT = 3
		DOUBLE = 4
	}
	`)
	_ = data
	require.NoError(err)
	require.Equal("Type", data.Enum[0].Name)
	require.Equal(2, len(data.Enum[0].Values))
	require.Equal("INT", data.Enum[0].Values[0].Value.Key)
	require.Equal(int(0), data.Enum[0].Values[0].Value.Value)
	require.Equal("DOUBLE", data.Enum[0].Values[1].Value.Key)
	require.Equal(int(1), data.Enum[0].Values[1].Value.Value)
	require.Equal("Type2", data.Enum[1].Name)
	require.Equal(2, len(data.Enum[1].Values))
	require.Equal("INT", data.Enum[1].Values[0].Value.Key)
	require.Equal(int(3), data.Enum[1].Values[0].Value.Value)
	require.Equal("DOUBLE", data.Enum[1].Values[1].Value.Key)
	require.Equal(int(4), data.Enum[1].Values[1].Value.Value)
}

func TestMessage(t *testing.T) {
	require := require.New(t)
	parser, err := participle.Build[Grammar](
		participle.Lexer(def),
		participle.Unquote(),
		participle.Elide("Whitespace"),
	)
	require.NoError(err)
	data, err := parser.ParseString("", `
	package example
	message SearchRequest {
		string query = 1
		int32 page_number = 2
		int32 result_per_page = 3
		map<string, double> scores = 4
	  
		message Foo {}
	  
		enum Bar {
		  FOO = 0
		}
	  }
	  
	  message SearchResponse {
		string results = 1
	  }
	`)
	_ = data
	require.NoError(err)
	require.Equal(2, len(data.Message))
	msg := data.Message[0]
	require.Equal("SearchRequest", msg.Name)
	require.Equal(6, len(msg.Entries))
	require.Equal("query", msg.Entries[0].Field.Name)
	require.Equal("String", msg.Entries[0].Field.Type.Scalar.String())
	require.Equal(int(1), msg.Entries[0].Field.Tag)
	require.Equal("page_number", msg.Entries[1].Field.Name)
	require.Equal("Int32", msg.Entries[1].Field.Type.Scalar.String())
	require.Equal(int(2), msg.Entries[1].Field.Tag)
	require.Equal("result_per_page", msg.Entries[2].Field.Name)
	require.Equal("Int32", msg.Entries[2].Field.Type.Scalar.String())
	require.Equal(int(3), msg.Entries[2].Field.Tag)
	require.Equal("scores", msg.Entries[3].Field.Name)
	require.Equal("String", msg.Entries[3].Field.Type.Map.Key.Scalar.String())
	require.Equal("Double", msg.Entries[3].Field.Type.Map.Value.Scalar.String())
	require.Equal(int(4), msg.Entries[3].Field.Tag)
	require.Equal("Foo", msg.Entries[4].Message.Name)
	require.Equal("Bar", msg.Entries[5].Enum.Name)
	require.Equal("FOO", msg.Entries[5].Enum.Values[0].Value.Key)
	require.Equal(int(0), msg.Entries[5].Enum.Values[0].Value.Value)
	require.Equal("SearchResponse", data.Message[1].Name)
	require.Equal(1, len(data.Message[1].Entries))
	require.Equal("results", data.Message[1].Entries[0].Field.Name)
	require.Equal("String", data.Message[1].Entries[0].Field.Type.Scalar.String())
	require.Equal(int(1), data.Message[1].Entries[0].Field.Tag)
}
