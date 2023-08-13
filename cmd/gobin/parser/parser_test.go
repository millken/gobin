package parser_test

import (
	"testing"

	"gobin/parser"

	"github.com/alecthomas/assert/v2"
)

type TestCode struct {
	Name    string
	Content string
}

func TestParseString(t *testing.T) {
	testCases := []TestCode{
		{
			Name:    "basic",
			Content: `package example`,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := parser.ParseString(testCase.Content)
			assert.NoError(t, err)
		})
	}
}

func TestPackage(t *testing.T) {
	data, err := parser.ParseString(`
  package example
  option go_marshal = "unsafe"
  option go_int = 3
  const int32 a = 1
  const float b = 1.1
  struct Foo {
	int32 a
	float b
  }
	`)
	assert.NoError(t, err)
	assert.Equal(t, "example", data.Package.Identifier.String)
	idx := 0
	opts := data.TopLevelDeclarations[idx].(parser.Option)
	assert.Equal(t, "go_marshal", opts.Name.String)
	assert.Equal(t, "unsafe", opts.Value.(parser.LiteralString).Value)
	idx++
	opts = data.TopLevelDeclarations[idx].(parser.Option)
	assert.Equal(t, "go_int", opts.Name.String)
	assert.Equal(t, int(3), opts.Value.(parser.LiteralInt).Value)
	idx++
	cons := data.TopLevelDeclarations[idx].(parser.Const)
	assert.Equal(t, "a", cons.Name.String)
	assert.Equal(t, "Int32", cons.Type.String())
	assert.Equal(t, int(1), cons.Value.(parser.LiteralInt).Value)
	idx++
	cons = data.TopLevelDeclarations[idx].(parser.Const)
	assert.Equal[string](t, "b", cons.Name.String)
	assert.Equal[string](t, "Float", cons.Type.String())
	assert.Equal[float64](t, 1.1, cons.Value.(parser.LiteralFloat).Value)
	idx++
	stru := data.TopLevelDeclarations[idx].(parser.Struct)
	assert.Equal[string](t, "Foo", stru.Name.String)
	assert.Equal[int](t, 2, len(stru.Fields))
	assert.Equal[string](t, "a", stru.Fields[0].Name.String)
	assert.Equal[parser.Type](t, parser.Int32, *stru.Fields[0].Type)
	assert.Equal[string](t, "b", stru.Fields[1].Name.String)
	assert.Equal[parser.Type](t, parser.Float, *stru.Fields[1].Type)

}

func TestParserGrammar(t *testing.T) {
	expected := `FileTopLevel = Package TopLevelDeclaration* .
Package = "package" Name .
Name = <ident> .
TopLevelDeclaration = Struct | Const .
Struct = "struct" Name ("<" (Name ("," Name)*)? ">")? .
Const = "const" Type Name "=" Literal .
Literal = LiteralFloat | LiteralInt | LiteralString | LiteralBool | LiteralNull .
LiteralFloat = <float> .
LiteralInt = <int> .
LiteralString = <string> .
LiteralBool = "true" | "false" .
LiteralNull = "null" .`
	grammar, err := parser.Grammar()
	//t.Log(grammar)
	assert.NoError(t, err)
	assert.Equal(t, expected, grammar)
}
