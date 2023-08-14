package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

func Grammar() (string, error) {
	p, err := parser()
	if err != nil {
		return "", err
	}
	return p.String(), nil
}

func ParseString(input string) (*FileTopLevel, error) {
	p, err := parser()
	if err != nil {
		return nil, err
	}
	ast, err := p.ParseString("", input)
	if err != nil {
		return nil, err
	}
	return ast, nil
}

func parser() (*participle.Parser[FileTopLevel], error) {
	return participle.Build[FileTopLevel](topLevelDeclarationUnion, literalUnion, participle.Unquote())
}

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

var topLevelDeclarationUnion = participle.Union[TopLevelDeclaration](Option{}, Struct{}, Const{})

type Option struct {
	Name  Name    `"option" @@`
	Value Literal `"=" @@`
}

func (o Option) sealedTopLevelDeclaration() {}

type Const struct {
	Type  *Type   `"const" @@`
	Name  Name    `@@`
	Value Literal `"=" @@`
}

func (c Const) sealedTopLevelDeclaration() {}

type Struct struct {
	Name   Name          `"struct" @@`
	Fields []StructField `"{" @@* "}"`
}

type StructField struct {
	Type *Type `@@`
	Name Name  `@@`
}

func (s Struct) sealedTopLevelDeclaration() {}

func TopLevelDeclarationExhaustiveSwitch(
	topLevelDeclaration TopLevelDeclaration,
	caseOption func(topLevelDeclaration Option),
	caseConst func(topLevelDeclaration Const),
	caseStruct func(topLevelDeclaration Struct),
) {
	opt, ok := topLevelDeclaration.(Option)
	if ok {
		caseOption(opt)
		return
	}
	cons, ok := topLevelDeclaration.(Const)
	if ok {
		caseConst(cons)
		return
	}
	struc, ok := topLevelDeclaration.(Struct)
	if ok {
		caseStruct(struc)
		return
	}
}