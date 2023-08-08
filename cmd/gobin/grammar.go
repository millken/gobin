package main

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	def = lexer.MustSimple([]lexer.SimpleRule{
		{"Int", `\d+`},
		{"Float", `[-+]?\d*\.\d+([eE][-+]?\d+)?`},
		{"Ident", `\w+`},
		{"String", `"[^"]*"`},
		{"Whitespace", `\s+`},
		{"Punct", `[,.<>(){}=:]`},
		{"Comment", `//.*`},
	})
	parser = participle.MustBuild[Grammar](
		//participle.Lexer(def),
		participle.Unquote(),
		participle.UseLookahead(2),
		//participle.Elide("Whitespace"),
	)
)

type Grammar struct {
	Pos     lexer.Position
	Package string     ` "package" @(Ident ( "." Ident )*)`
	Option  []*Option  `@@*`
	Consts  []*Const   `@@*`
	Enum    []*Enum    `@@*`
	Message []*Message `@@*`
	//Entries []*Entry `| ( @@ ";"* )*`
}

type Enum struct {
	Pos lexer.Position

	Name   string       `"enum" @Ident`
	Values []*EnumEntry `"{" ( @@ ( ";" )* )* "}"`
}

type EnumEntry struct {
	Pos lexer.Position

	Value  *EnumValue `  @@`
	Option *Option    `| "option" @@`
}

type EnumValue struct {
	Pos lexer.Position

	Key   string `@Ident`
	Value int    `"=" @( [ "-" ] Int )`

	Options []*Option `( "[" @@ ( "," @@ )* "]" )?`
}
type Scalar int

const (
	None Scalar = iota
	Double
	Float
	Int32
	Int64
	Uint32
	Uint64
	Sint32
	Sint64
	Fixed32
	Fixed64
	SFixed32
	SFixed64
	Bool
	String
	Bytes
)

var scalarToString = map[Scalar]string{
	None: "None", Double: "Double", Float: "Float", Int32: "Int32", Int64: "Int64", Uint32: "Uint32",
	Uint64: "Uint64", Sint32: "Sint32", Sint64: "Sint64", Fixed32: "Fixed32", Fixed64: "Fixed64",
	SFixed32: "SFixed32", SFixed64: "SFixed64", Bool: "Bool", String: "String", Bytes: "Bytes",
}

func (s Scalar) String() string { return scalarToString[s] }

var stringToScalar = map[string]Scalar{
	"double": Double, "float": Float, "int32": Int32, "int64": Int64, "uint32": Uint32, "uint64": Uint64,
	"sint32": Sint32, "sint64": Sint64, "fixed32": Fixed32, "fixed64": Fixed64, "sfixed32": SFixed32,
	"sfixed64": SFixed64, "bool": Bool, "string": String, "bytes": Bytes,
}

func (s *Scalar) Parse(lex *lexer.PeekingLexer) error {
	token := lex.Peek()
	v, ok := stringToScalar[token.Value]
	if !ok {
		return participle.NextMatch
	}
	lex.Next()
	*s = v
	return nil
}

type Type struct {
	Pos lexer.Position

	Scalar    Scalar   `  @@`
	Map       *MapType `| @@`
	Reference string   `| @(Ident ( "." Ident )*)`
}

type MapType struct {
	Pos lexer.Position

	Key   *Type `"map" "<" @@`
	Value *Type `"," @@ ">"`
}
type Const struct {
	Pos   lexer.Position
	Type  *Type    `"const" @@`
	Name  string   `@Ident`
	Value *Literal `"=" @@`
}

type Message struct {
	Pos lexer.Position

	Name    string          `"message" @Ident`
	Entries []*MessageEntry `"{" @@* "}"`
}

type MessageEntry struct {
	Pos lexer.Position

	Enum    *Enum    `( @@`
	Option  *Option  ` | "option" @@`
	Message *Message ` | @@`
	Field   *Field   ` | @@ ) ";"*`
}

type Field struct {
	Pos lexer.Position

	Optional bool `(   @"optional"`
	Required bool `  | @"required"`
	Repeated bool `  | @"repeated" )?`

	Type *Type  `@@`
	Name string `@Ident`
	Tag  int    `"=" @Int`

	Options []*Option `( "[" @@ ( "," @@ )* "]" )?`
}

// Literal is a "union" type, where only one matching value will be present.
type Literal struct {
	Pos       lexer.Position
	Str       *string    `  @String`
	Int       *int       `| @Int`
	Float     *float64   `| @(Float|Int)`
	Bool      *string    `| @( "true" | "false" )`
	Reference *string    `| @Ident ( @"." @Ident )*`
	Minus     *Literal   `| "-" @@`
	List      []*Literal `| "[" ( @@ ","? )* "]"`
	Map       []*MapItem `| "{" ( @@ ","? )* "}"`
}
type MapItem struct {
	Pos   lexer.Position
	Key   *Literal `@@ ":"`
	Value *Literal `@@`
}

func (m *MapItem) GoString() string {
	return fmt.Sprintf("%v: %v", m.Key, m.Value)
}

type Option struct {
	Pos   lexer.Position
	Name  string `"option" @Ident`
	Value *Value `"=" @@`
}

type Value struct {
	Pos lexer.Position

	String    *string `  @String`
	Int       *int64  `| @Int`
	Bool      *bool   `| (@"true" | "false")`
	Reference *string `| @Ident @( "." Ident )*`
	Map       *Map    `| @@`
	Array     *Array  `| @@`
}

type Array struct {
	Pos lexer.Position

	Elements []*Value `"[" ( @@ ( ","? @@ )* )? "]"`
}

type Map struct {
	Pos lexer.Position

	Entries []*MapEntry `"{" ( @@ ( ( "," )? @@ )* )? "}"`
}

type MapEntry struct {
	Pos lexer.Position

	Key   *Value `@@`
	Value *Value `":"? @@`
}
