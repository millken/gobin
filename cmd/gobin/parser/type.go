package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Type int

const (
	None Type = iota
	Double
	Float
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Bool
	String
	Bytes
)

var typeToString = map[Type]string{
	None: "None", Double: "Double", Float: "Float", Int: "Int", Int8: "Int8", Int16: "Int16", Int32: "Int32", Int64: "Int64", Uint: "Uint", Uint8: "Uint8", Uint16: "Uint16", Uint32: "Uint32", Uint64: "Uint64", Bool: "Bool", String: "String", Bytes: "Bytes",
}

func (t Type) String() string   { return typeToString[t] }
func (t Type) GoString() string { return typeToGoType[t] }

var stringToType = map[string]Type{
	"none": None, "double": Double, "float": Float, "int": Int, "int8": Int8, "int16": Int16, "int32": Int32, "int64": Int64, "uint": Uint, "uint8": Uint8, "uint16": Uint16, "uint32": Uint32, "uint64": Uint64, "bool": Bool, "string": String, "bytes": Bytes,
}

var typeToGoType = map[Type]string{
	None: "None", Double: "float64", Float: "float32", Int: "int", Int8: "int8", Int16: "int16", Int32: "int32", Int64: "int64", Uint: "uint", Uint8: "uint8", Uint16: "uint16", Uint32: "uint32", Uint64: "uint64", Bool: "bool", String: "string", Bytes: "[]byte",
}

func (t Type) Size() int {
	switch t {
	case Double, Int64, Uint64:
		return 8
	case Float, Int32, Uint32:
		return 4
	case Int, Int8, Uint, Uint8:
		return 1
	case Int16, Uint16:
		return 2
	}
	return 0
}

func (t *Type) Parse(lex *lexer.PeekingLexer) error {
	token := lex.Peek()
	v, ok := stringToType[token.Value]
	if !ok {
		return participle.NextMatch
	}
	lex.Next()
	*t = v
	return nil
}
