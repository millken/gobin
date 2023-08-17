package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io"
	"os"
	"strconv"
	"strings"

	"gobin/parser"
)

type option func(*Parser) error

func WithFormatted() option {
	return func(p *Parser) error {
		p.formatted = true
		return nil
	}
}

type Parser struct {
	buf    *bufio.Reader
	out    *bytes.Buffer
	parser *parser.FileTopLevel
	option map[string]parser.Literal

	formatted bool
}

func NewParser(out *bytes.Buffer, src any, opts ...option) (*Parser, error) {
	var err error
	p := &Parser{out: out, option: make(map[string]parser.Literal)}
	for _, opt := range opts {
		if err = opt(p); err != nil {
			return nil, err
		}
	}
	if src != nil {
		switch s := src.(type) {
		case string:
			p.buf = bufio.NewReader(strings.NewReader(s))
		case []byte:
			p.buf = bufio.NewReader(bytes.NewReader(s))
		case *bytes.Buffer:
			p.buf = bufio.NewReader(s)
		case io.Reader:
			p.buf = bufio.NewReader(s)
		default:
			err = errors.New("invalid source")
		}
		return p, err
	} else {
		p.buf = bufio.NewReader(os.Stdin)
	}
	return p, err
}

func (p *Parser) Error() error {
	return nil
}

func (p *Parser) Parse() error {
	buf := &bytes.Buffer{}
	_, err := buf.ReadFrom(p.buf)
	if err != nil {
		return err
	}
	parser, err := parser.ParseString(buf.String())
	if err != nil {
		return err
	}
	options, consts, enums, structs := splitTopLevelDeclarations(parser.TopLevelDeclarations)
	//parse package
	if err := p.parsePackage(parser.Package.Identifier.String); err != nil {
		return errors.New("parsePackage error: " + err.Error())
	}
	//parse option
	if err := p.parseOption(options); err != nil {
		return errors.New("parseOption error: " + err.Error())
	}
	//parse const
	if err := p.parseConst(consts); err != nil {
		return errors.New("parseConst error: " + err.Error())
	}

	// parse enum
	if err := p.parseEnum(enums); err != nil {
		return errors.New("parseEnum error: " + err.Error())
	}
	//parse struct
	if err := p.parseStruct(structs); err != nil {
		return errors.New("parseStruct error: " + err.Error())
	}
	if p.formatted {
		//format output
		formatSrc, err := format.Source(p.out.Bytes())
		if err != nil {
			return errors.New("format error: " + err.Error())
		}
		p.out.Reset()
		p.out.Write(formatSrc)
	} else {
		p.out.Write([]byte("\n"))
	}

	return nil
}
func (p *Parser) parsePackage(name string) error {
	err := prologTemplate.ExecuteTemplate(p.out, "prolog", name)
	return err
}

func (p *Parser) parseStruct(structs []parser.Struct) error {
	if len(structs) > 0 {
		if err := structTemplate.ExecuteTemplate(p.out, "struct", map[string]any{"Structs": structs, "Options": p.option}); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) parseOption(options []parser.Option) error {
	for _, option := range options {
		p.option[option.Name.String] = option.Value
	}
	return nil
}

func (p *Parser) parseConst(consts []parser.Const) error {
	if len(consts) > 0 {
		err := constTemplate.ExecuteTemplate(p.out, "const", map[string]any{"Consts": consts, "Options": p.option})
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseEnum(enums []parser.Enum) error {
	if len(enums) > 0 {
		err := enumTemplate.ExecuteTemplate(p.out, "enum", map[string]any{"Enums": enums, "Options": p.option})
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateLiteral(literal parser.Literal) string {
	result := ""
	parser.LiteralExhaustiveSwitch(
		literal,
		func(literal float64) { result = fmt.Sprintf("%f", literal) },
		func(literal int) { result = fmt.Sprintf("%d", literal) },
		func(literal string) { result = literal },
		func(literal bool) { result = strconv.FormatBool(literal) },
		func() { result = "nil" },
	)
	return result
}

func splitTopLevelDeclarations(topLevelDeclarations []parser.TopLevelDeclaration) ([]parser.Option, []parser.Const, []parser.Enum, []parser.Struct) {
	options := []parser.Option{}
	consts := []parser.Const{}
	structs := []parser.Struct{}
	enums := []parser.Enum{}

	for _, topLevelDeclaration := range topLevelDeclarations {
		parser.TopLevelDeclarationExhaustiveSwitch(
			topLevelDeclaration,
			func(topLevelDeclaration parser.Option) {
				options = append(options, topLevelDeclaration)
			},
			func(topLevelDeclaration parser.Const) {
				topLevelDeclaration.Name.String = UpperFirst(topLevelDeclaration.Name.String)
				consts = append(consts, topLevelDeclaration)
			},
			func(topLevelDeclaration parser.Enum) {
				topLevelDeclaration.Name.String = UpperFirst(topLevelDeclaration.Name.String)
				for i, field := range topLevelDeclaration.Values {
					field.Value = UpperFirst(field.Value)
					topLevelDeclaration.Values[i] = field
				}
				enums = append(enums, topLevelDeclaration)
			},
			func(topLevelDeclaration parser.Struct) {
				topLevelDeclaration.Name.String = UpperFirst(topLevelDeclaration.Name.String)
				for i, field := range topLevelDeclaration.Fields {
					field.Name.String = UpperFirst(field.Name.String)
					topLevelDeclaration.Fields[i] = field
				}
				structs = append(structs, topLevelDeclaration)
			},
		)
	}
	return options, consts, enums, structs
}
