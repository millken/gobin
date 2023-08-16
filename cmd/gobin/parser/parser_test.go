package parser_test

import (
	"os"
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
  //comment 1
  option go_marshal = "unsafe"
  /* 
  comment 2
  abc
  */
  option go_int = 3

  //comment 12
  const int32 a = 1
  /*
multi line comment
multi line comment
  */
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
	assert.Equal[parser.Type](t, parser.Int32, *stru.Fields[0].Type.Type)
	assert.Equal[string](t, "b", stru.Fields[1].Name.String)
	assert.Equal[parser.Type](t, parser.Float, *stru.Fields[1].Type.Type)
}

func TestSchema(t *testing.T) {
	data, err := parser.ParseString(`
  package example
  struct hole  {
	// Lat is the latitude of the cup.
	double lat
	// Lon is the longitude of the cup.
	double lon
	// Par is the difficulty index.
	uint8 par
	// Water marks the presence of water.
	bool water 
	// Sand marks the presence of sand.
	bool sand 
}
  // Course is the grounds where the game of golf is played.
  struct course  {
	uint64 ID    
	string name
    hole holes [repeated = true] 
	bytes image  
    string tags  [repeated = true]
}
	`)
	assert.NoError(t, err)
	assert.Equal(t, "example", data.Package.Identifier.String)
	idx := 0
	idx++
}

func TestPebble(t *testing.T) {
	input, err := os.ReadFile("../testdata/pebble.gobin")
	assert.NoError(t, err)
	data, err := parser.ParseBytes(input)
	assert.NoError(t, err)
	assert.Equal(t, "pebble", data.Package.Identifier.String)
	idx := 0
	enum := data.TopLevelDeclarations[idx].(parser.Enum)
	assert.Equal[string](t, "PackageType", enum.Name.String)
	assert.Equal[int](t, 3, len(enum.Values))
	assert.Equal[string](t, "DATA", enum.Values[0].Value)
	assert.Equal[string](t, "CONFIG", enum.Values[1].Value)
	assert.Equal[string](t, "STATE", enum.Values[2].Value)
}

func TestParserGrammar(t *testing.T) {
	expected := `FileTopLevel = Package TopLevelDeclaration* .
Package = <comment>* "package" Name .
Name = <ident> .
TopLevelDeclaration = Option | Struct | Const .
Option = <comment>* "option" Name "=" Literal .
Literal = LiteralFloat | LiteralInt | LiteralString | LiteralBool | LiteralNull .
LiteralFloat = <float> .
LiteralInt = <int> .
LiteralString = <string> .
LiteralBool = "true" | "false" .
LiteralNull = "null" .
Struct = <comment>* "struct" Name "{" StructField* "}" .
StructField = <comment>* StructType Name ("[" StructOption ("," StructOption)* "]")? .
StructType = Type | <ident> .
StructOption = (("(" <ident> ("." <ident>)* ")") | (<ident> ("." <ident>)*)) "=" Literal .
Const = <comment>* "const" Type Name "=" Literal .`
	grammar, err := parser.Grammar()
	//t.Log(grammar)
	assert.NoError(t, err)
	assert.Equal(t, expected, grammar)
}
