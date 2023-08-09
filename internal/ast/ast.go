package ast

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type AST struct {
	Pos lexer.Position
	// Comments []string `@Comment*`
	Package *Package `@@`
	// Entries          Entries        `parser:"@@*"`
	// Schema           bool           `parser:""`
}

type Package struct {
	Pos  lexer.Position
	Name string `"package" @Ident`
}

var (
	basicLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"Comment", `(?:(?://|#)[^\n]*)|/\*.*?\*/`},
		{"String", `"[^"]*"`},
		// {"Number", `[-+]?(\d*\.)?\d+`},
		{`Ident`, `[a-zA-Z][a-zA-Z0-9]*`},
		// {"Punct", `[,.<>(){}=:]`},
		// {"EOL", `[\n\r]+`},
		{"Whitespace", `[ \t]+`},
	})

	basicParser = participle.MustBuild[AST](
		participle.Lexer(basicLexer),
		//participle.CaseInsensitive("Ident"),
		participle.Unquote(),
		participle.Elide("Whitespace"),
		participle.UseLookahead(2),
	)
)
