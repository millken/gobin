package main

import (
	"fmt"
	"gobin/parser"
	"strconv"
	"strings"
	"text/template"
)

//nolint:gochecknoglobals
var (
	typeToString = map[parser.Type]string{
		parser.String: "String",
		parser.Int:    "Int",
		parser.Int8:   "Int8",
		parser.Int16:  "Int16",
		parser.Int32:  "Int32",
		parser.Int64:  "Int64",
		parser.Uint:   "Uint",
		parser.Uint8:  "Uint8",
		parser.Uint16: "Uint16",
		parser.Uint32: "Uint32",
		parser.Uint64: "Uint64",
		parser.Float:  "Float32",
		parser.Double: "Float64",
		parser.Bool:   "Bool",
		parser.Bytes:  "Bytes",
	}
	funcMap = []template.FuncMap{map[string]interface{}{
		"StructFieldLength": func(fields []parser.StructField) string {
			var n int
			var ret string
			for _, f := range fields {
				if f.Type.Type == nil {
					// reference to another struct
					ret += fmt.Sprintf(`o.%s.Size() + `, f.Name.String)
					continue
				}
				if sz := f.Type.Type.Size(); sz > 0 {
					n += sz
				} else {
					ret += fmt.Sprintf(`len(o.%s) + `, f.Name.String)
					n += strconv.IntSize / 8
				}
			}
			return ret + fmt.Sprintf("%d", n)
		},
		"StructFieldMarshal": func(fields []parser.StructField) string {
			var ret string
			for _, f := range fields {
				if f.Type.Type == nil {
					ret += fmt.Sprintf(`if n, err = o.%s.MarshalTo(data); err != nil {
						return nil, err
					}
					offset += n
					`, f.Name.String)
					continue
				}
				if v, ok := typeToString[*f.Type.Type]; ok {
					ret += fmt.Sprintf(`if n, err = o.Marshal%s(o.%s, data[offset:]); err != nil {
						return nil, err
					}
					offset += n
`, v, f.Name.String)
				} else {
					panic("unknown type")
				}
			}
			ret += `	if offset != sz {
				return nil, fmt.Errorf("%s size / offset different %d : %d", "Marshal", sz, offset)
			}`
			return ret
		},
		"StructFieldUnmarshal": func(fields []parser.StructField) string {
			var ret string
			for _, f := range fields {
				if f.Type.Type == nil {
					// reference to another struct
					continue
				}
				if v, ok := typeToString[*f.Type.Type]; ok {
					ret += fmt.Sprintf(`if o.%s, i, err = o.Unmarshal%s(data[n:]); err != nil {
						return err
					}
					n += i
					`, f.Name.String, v)
				} else {
					panic("unknown type")
				}
			}
			return ret
		},
		"FormatComment": func(comment string) string {
			comments := ""
			for _, c := range strings.Split(comment, "\n") {
				c = strings.TrimSpace(c)
				if c == "" {
					continue
				}
				comments += fmt.Sprintf("// %s\n", c)
			}
			return strings.TrimRight(comments, "\n")
		},
		"GetString": func(s *string) string {
			if s == nil {
				return ""
			}
			return *s
		},
	}}

	prologTemplate = template.Must(template.New("prolog").Parse(`
package {{ . }}

import (
	"fmt"
	"github.com/millken/gobin"
)
`))

	constTemplate, enumTemplate, structTemplate *template.Template

	initTemplate = template.Must(template.New("init").Parse(`
func init() {
	{{- range $i, $sd := . }}
	parse.AssertUpToDate(&{{ $sd.TableVar }}.s, new({{ $sd.Type }}))
	{{- end }}
}
`))
)

func init() {
	constTemplateTmp := template.New("const")
	structTemplateTmp := template.New("struct")
	enumTemplateTmp := template.New("enum")
	for _, f := range funcMap {
		constTemplateTmp.Funcs(f)
		structTemplateTmp.Funcs(f)
		enumTemplateTmp.Funcs(f)
	}
	constTemplate = template.Must(constTemplateTmp.Parse(`
	const (
		{{- range $c := .Consts }}
		{{- if $c.Comments }}
		{{ $c.Comments | FormatComment }}
		{{- end }}
		{{ $c.Name.String }} {{ $c.Type.GoString }} = {{ $c.Value.GoString }}
		{{- end }}
	)
	`))

	enumTemplate = template.Must(enumTemplateTmp.Parse(`
	{{range $parent := .Enums}}
	{{- if $parent.Comments }}
	{{ $parent.Comments | FormatComment }}
	{{- end }}
	type {{$parent.Name.String}} uint16
	const (
		{{- range $i,$v := $parent.Values}}
		{{- if $v.Comments }}
		{{ $v.Comments | FormatComment }}
		{{- end}}
		{{ $parent.Name.String }}_{{$v.Value}} {{$parent.Name.String}} = {{$i}}
		{{- end }}
	)

	func (o *{{$parent.Name.String}}) Size() int {
		return 2
	}

	// MarshalTo writes a wire-format message to w.
	func (o *{{$parent.Name.String}}) MarshalTo(w []byte) (int, error) {
		w[0] = byte(*o)
		w[1] = byte(*o >> 8)
		return 2, nil
	}

	// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
	func (o *{{$parent.Name.String}}) MarshalBinary() ([]byte,error) {
		data := make([]byte, 2)
		_, err := o.MarshalTo(data)
		return data, err
	}

	// Unmarshal decodes data as conform encoding.BinaryUnmarshaler.
func (o *{{$parent.Name.String}}) UnmarshalBinary(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("invalid data size %d", len(data))
	}
	*o = {{$parent.Name.String}}(uint16(data[0]) | uint16(data[1])<<8)
	return nil
}
	{{- end }}
	`))

	structTemplate = template.Must(structTemplateTmp.Parse(`
{{- range .Structs}}
{{- if .Comments }}
{{ .Comments | FormatComment }}
{{- end }}
type {{.Name.String}} struct {
	gobin.{{with $.Options.go_marshal}}{{if eq .Value "unsafe"}}Unsafe{{else}}Safe{{end}}{{else}}Safe{{end}}
{{- range .Fields}}
{{- if .Comments }}
{{ .Comments | FormatComment }}
{{- end }}
{{- if .Type.Type | eq nil }}
	{{.Name.String}} {{with .Type.Reference}} *{{GetString .}} {{end}}
{{- else}}
	{{.Name.String}} {{.Type.Type.GoString}}
{{- end}}
{{- end}}
}

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *{{.Name.String}}) MarshalBinary() (data []byte, err error) {
	sz := {{.Fields | StructFieldLength}}
	data = make([]byte, sz)
	var offset, n int
	{{.Fields | StructFieldMarshal}}
	return data, nil
}

// Unmarshal decodes data as conform encoding.BinaryUnmarshaler.
func (o *{{.Name.String}}) UnmarshalBinary(data []byte) error {
	var (
		i, n int
		err  error
	)
	{{.Fields | StructFieldUnmarshal}}
	return nil
}
{{- end}}
`))
}
