package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
)

var literalUnion = participle.Union[Literal](LiteralFloat{}, LiteralInt{}, LiteralString{}, LiteralBool{}, LiteralNull{})

type Literal interface {
	sealedLiteral()
	GoString() string
}

type LiteralFloat struct {
	Value float64 `@Float`
}

func (literal LiteralFloat) sealedLiteral() {}
func (literal LiteralFloat) GoString() string {
	return fmt.Sprintf("%v", literal.Value)
}

type LiteralInt struct {
	Value int `@Int`
}

func (literal LiteralInt) sealedLiteral() {}
func (literal LiteralInt) GoString() string {
	return fmt.Sprintf("%d", literal.Value)
}

type LiteralString struct {
	Value string `@String`
}

func (literal LiteralString) sealedLiteral() {}
func (literal LiteralString) GoString() string {
	return fmt.Sprintf("\"%s\"", literal.Value)
}

type LiteralBool struct {
	Value bool `@"true"`
	False bool `| @"false"`
}

func (literal LiteralBool) sealedLiteral() {}
func (literal LiteralBool) GoString() string {
	return fmt.Sprintf("%v", literal.Value)
}

type LiteralNull struct {
	Value bool `@"null"`
}

func (literal LiteralNull) sealedLiteral() {}
func (literal LiteralNull) GoString() string {
	return fmt.Sprintf("%v", literal.Value)
}

func LiteralExhaustiveSwitch(
	literal Literal,
	caseFloat func(literal float64),
	caseInt func(literal int),
	caseString func(literal string),
	caseBool func(literal bool),
	caseNull func(),
) {
	litFloat, ok := literal.(LiteralFloat)
	if ok {
		caseFloat(litFloat.Value)
		return
	}
	litInt, ok := literal.(LiteralInt)
	if ok {
		caseInt(litInt.Value)
		return
	}
	litString, ok := literal.(LiteralString)
	if ok {
		caseString(litString.Value)
		return
	}
	litBool, ok := literal.(LiteralBool)
	if ok {
		caseBool(litBool.Value)
		return
	}
	_, ok = literal.(LiteralNull)
	if ok {
		caseNull()
		return
	}
}
