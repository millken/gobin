package parser

import "github.com/alecthomas/participle/v2/lexer"

type Node struct {
	Pos    lexer.Position
	EndPos lexer.Position
}

type FileTopLevel struct {
	Package              Package               `@@`
	TopLevelDeclarations []TopLevelDeclaration `@@*`
}
type Name struct {
	Node
	String string `@Ident`
}

type TopLevelDeclaration interface {
	sealedTopLevelDeclaration()
}

type Package struct {
	Identifier Name `"package" @@`
}

type Struct struct {
	Name     Name   `"struct" @@`
	Generics []Name `("<" (@@ ("," @@)*)? ">")?`
}

func (s Struct) sealedTopLevelDeclaration() {}
