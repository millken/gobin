package parser

import (
	"regexp"
	"text/scanner"

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

func ParseBytes(input []byte) (*FileTopLevel, error) {
	p, err := parser()
	if err != nil {
		return nil, err
	}
	ast, err := p.ParseBytes("", input)
	if err != nil {
		return nil, err
	}
	return ast, nil
}

func parser() (*participle.Parser[FileTopLevel], error) {
	lex := lexer.NewTextScannerLexer(func(s *scanner.Scanner) {
		s.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanChars | scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanComments
	})
	return participle.Build[FileTopLevel](participle.Lexer(lex),
		topLevelDeclarationUnion,
		literalUnion,
		participle.Map(stripComment, "Comment"), // for multi language
		participle.Unquote())
}

var stripCommentRe = regexp.MustCompile(`^//\s*|^/\*|\*/$`)

func stripComment(token lexer.Token) (lexer.Token, error) {
	token.Value = stripCommentRe.ReplaceAllString(token.Value, "")
	return token, nil
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
	Comments   string `@Comment?`
	Identifier Name   `"package" @@`
}

var topLevelDeclarationUnion = participle.Union[TopLevelDeclaration](Option{}, Struct{}, Const{}, Enum{})

type Option struct {
	Comments string  `@Comment?`
	Name     Name    `"option" @@`
	Value    Literal `"=" @@`
}

func (o Option) sealedTopLevelDeclaration() {}

type Const struct {
	Comments string  `@Comment?`
	Type     *Type   `"const" @@`
	Name     Name    `@@`
	Value    Literal `"=" @@`
}

func (c Const) sealedTopLevelDeclaration() {}

type Enum struct {
	Comments string      `@Comment?`
	Name     Name        `"enum" @@`
	Values   []EnumValue `"{" @@* "}"`
}

func (c Enum) sealedTopLevelDeclaration() {}

type EnumValue struct {
	Comments string `@Comment?`
	Value    string `@Ident`
}

type Struct struct {
	Comments string        `@Comment?`
	Name     Name          `"struct" @@`
	Fields   []StructField `"{" @@* "}"`
}

type StructField struct {
	Comments string          `@Comment?`
	Type     *StructType     `@@`
	Name     Name            `@@`
	Options  []*StructOption `( "[" @@ ( "," @@ )* "]" )?`
}

type StructType struct {
	Type      *Type   `@@`
	Reference *string `| @Ident`
}

type StructOption struct {
	Name  string  `( "(" @Ident @( "." Ident )* ")" | @Ident @( "." @Ident )* )`
	Value Literal `"=" @@`
}

func (s Struct) sealedTopLevelDeclaration() {}

func TopLevelDeclarationExhaustiveSwitch(
	topLevelDeclaration TopLevelDeclaration,
	caseOption func(topLevelDeclaration Option),
	caseConst func(topLevelDeclaration Const),
	caseEnum func(topLevelDeclaration Enum),
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
	enum, ok := topLevelDeclaration.(Enum)
	if ok {
		caseEnum(enum)
		return
	}
	struc, ok := topLevelDeclaration.(Struct)
	if ok {
		caseStruct(struc)
		return
	}
}
