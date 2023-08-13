package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type AST struct {
	Pos      lexer.Position
	Comments []string `@Comment*`
	Package  string   `  "package" @Ident`
	Consts   []*Const `@@*`
	// Schema           bool           `parser:""`
}

type Const struct {
	Pos   lexer.Position
	Type  *Type   `"const" @@`
	Name  string  `@Ident`
	Value Literal `"=" @@`
}

type LiteralExpression struct {
	Pos     lexer.Position
	Literal Literal `@@`
}

func (l LiteralExpression) sealedExpression() {}

var (
	basicLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"Comment", `(?:(?://|#)[^\n]*)|/\*.*?\*/`},
		{"String", `"[^"]*"`},
		// {"Number", `[-+]?(\d*\.)?\d+`},
		{`Ident`, `[a-zA-Z][a-zA-Z0-9]*`},
		{"Punct", `[,.<>(){}=:]`},
		// {"EOL", `[\n\r]+`},
		{"Whitespace", `[ \n\t]+`},
	})

	basicParser = participle.MustBuild[AST](
		participle.Lexer(basicLexer),
		//participle.CaseInsensitive("Ident"),
		participle.Unquote(),
		participle.Elide("Whitespace"), // ignore whitespace tokens
		participle.UseLookahead(2),
	)
)

func ParseString(input string) (*AST, error) {
	ast, err := basicParser.ParseString("", input)
	if err != nil {
		return nil, err
	}
	return ast, nil
}
