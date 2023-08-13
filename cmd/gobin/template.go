package main

import (
	"fmt"
	"gobin/parser"
	"strconv"
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
				if sz := f.Type.Size(); sz > 0 {
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
				if v, ok := typeToString[*f.Type]; ok {
					ret += fmt.Sprintf(`i += gobin.Marshal%s(o.%s, data[i:])
`, v, f.Name.String)
				} else {
					panic("unknown type")
				}
			}
			return ret
		},
		"StructFieldUnmarshal": func(fields []parser.StructField) string {
			var ret string
			for _, f := range fields {
				if v, ok := typeToString[*f.Type]; ok {
					ret += fmt.Sprintf(`if o.%s, i, err = gobin.Unmarshal%s(data[n:]); err != nil {
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
	}}

	prologTemplate = template.Must(template.New("prolog").Parse(`
package {{ . }}

import (
	"github.com/millken/gobin"
)
`))
	constTemplate = template.Must(template.New("const").Parse(`
const (
	{{- range $c := . }}
	{{ $c.Name.String }} {{ $c.Type.GoString }} = {{ $c.Value.GoString }}
	{{- end }}
)
`))
	structTemplate *template.Template

	initTemplate = template.Must(template.New("init").Parse(`
func init() {
	{{- range $i, $sd := . }}
	parse.AssertUpToDate(&{{ $sd.TableVar }}.s, new({{ $sd.Type }}))
	{{- end }}
}
`))
)

func init() {
	structTemplateTmp := template.New("struct")
	for _, f := range funcMap {
		structTemplateTmp.Funcs(f)
	}
	structTemplate = template.Must(structTemplateTmp.Parse(`
{{- range .}}
type {{.Name.String}} struct {
{{- range .Fields}}
	{{.Name.String}} {{.Type.GoString}}
{{- end}}
}
func (o *{{.Name.String}}) MarshalBinary() (data []byte, err error) {
	sz := {{.Fields | StructFieldLength}}
	data = make([]byte, sz)
	var i int
	{{.Fields | StructFieldMarshal}}
	return data, nil
}

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
