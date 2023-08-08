package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Parser struct {
	buf    *bufio.Reader
	out    *bytes.Buffer
	parser *Grammar
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
	p.parser, err = parser.ParseBytes("", buf.Bytes())
	if err != nil {
		return err
	}
	//parse package
	if err := p.parsePackage(); err != nil {
		return err
	}
	//parse option
	//parse const
	//parse enum
	//parse message
	return nil
}
func (p *Parser) parsePackage() error {
	//_, err := p.out.WriteString("package " + p.parser.Package + "\n")
	err := prologTemplate.ExecuteTemplate(p.out, "prolog", p.parser.Package)
	return err
}

func (p *Parser) parseOption() error {
	return nil
}

func (p *Parser) parseConst() error {
	type constT struct {
		Name  string
		Type  string
		Value string
	}
	var consts []constT
	for _, c := range p.parser.Consts {
		var t string
		var v string
		switch c.Type.Scalar {
		case Int32:
			t = "int32"
			v = fmt.Sprintf("%d", *c.Value.Int)
		case Int64:
			t = "int64"
			v = fmt.Sprintf("%d", *c.Value.Int)
		case Uint32:
			t = "uint32"
			v = fmt.Sprintf("%d", *c.Value.Int)
		case Uint64:
			t = "uint64"
			v = fmt.Sprintf("%d", *c.Value.Int)
		case Float:
			t = "float32"
			v = fmt.Sprintf("%f", *c.Value.Float)
		case Double:
			t = "float64"
			v = fmt.Sprintf("%f", *c.Value.Float)
		case String:
			t = "string"
		case Bytes:
			t = "[]byte"
		case Bool:
			t = "bool"
		default:
			return errors.New("invalid type")
		}
		consts = append(consts, constT{
			Name:  c.Name,
			Type:  t,
			Value: v,
		})
	}
	return nil
}
