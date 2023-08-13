package main

import (
	"bufio"
	"bytes"
	"errors"
	"go/format"
	"io"
	"os"
	"strings"

	"gobin/parser"
)

type Parser struct {
	buf    *bufio.Reader
	out    *bytes.Buffer
	parser *parser.FileTopLevel
}

func NewParser(out *bytes.Buffer, src any) (*Parser, error) {
	var err error
	p := &Parser{out: out}
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
	consts, structs := splitTopLevelDeclarations(parser.TopLevelDeclarations)
	//parse package
	if err := p.parsePackage(parser.Package.Identifier.String); err != nil {
		return err
	}
	//parse option
	//parse const
	if err := p.parseConst(consts); err != nil {
		return err
	}
	//parse struct
	if err := p.parseStruct(structs); err != nil {
		return err
	}
	//format output
	formatSrc, err := format.Source(p.out.Bytes())
	if err != nil {
		return err
	}
	p.out.Reset()
	p.out.Write(formatSrc)
	return nil
}
func (p *Parser) parsePackage(name string) error {
	err := prologTemplate.ExecuteTemplate(p.out, "prolog", name)
	return err
}

func (p *Parser) parseStruct(structs []parser.Struct) error {
	if len(structs) > 0 {
		if err := structTemplate.ExecuteTemplate(p.out, "struct", structs); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) parseOption() error {
	return nil
}

func (p *Parser) parseConst(consts []parser.Const) error {
	if len(consts) > 0 {
		err := constTemplate.ExecuteTemplate(p.out, "const", consts)
		if err != nil {
			return err
		}
	}
	return nil
}

func splitTopLevelDeclarations(topLevelDeclarations []parser.TopLevelDeclaration) ([]parser.Const, []parser.Struct) {
	consts := []parser.Const{}
	structs := []parser.Struct{}
	for _, topLevelDeclaration := range topLevelDeclarations {
		parser.TopLevelDeclarationExhaustiveSwitch(
			topLevelDeclaration,
			func(topLevelDeclaration parser.Const) {
				consts = append(consts, topLevelDeclaration)
			},
			func(topLevelDeclaration parser.Struct) {
				structs = append(structs, topLevelDeclaration)
			},
		)
	}
	return consts, structs
}
