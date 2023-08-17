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
	IntSize      = strconv.IntSize / 8
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
		"StructFieldGetOption": func(name string, fields []*parser.StructOption) *parser.Literal {
			return getOption(name, fields)
		},
		"StructFieldIsRepeat": func(opts []*parser.StructOption) bool {
			if opts == nil {
				return false
			}
			opt := getOption("repeated", opts)
			return isBool(opt)
		},
		"StructOptionIsBool": func(opt *parser.Literal) bool {
			return isBool(opt)
		},
		"StructFieldLength": func(fields []parser.StructField) string {
			var n int
			var ret string
			for _, f := range fields {
				opt := getOption("repeated", f.Options)
				repeated := isBool(opt)
				if f.Type.Type == nil {
					if repeated {
						//array
						ret += fmt.Sprintf(`
						for _, v := range o.%s {
							sz += v.Size()
						}`, f.Name.String)

						//n is length of array
						//length + array
						n += IntSize
					} else {
						// reference to another struct
						ret += fmt.Sprintf(`
						sz += o.%s.Size()`, f.Name.String)
					}
					continue
				}
				if sz := f.Type.Type.Size(); sz > 0 {
					if repeated {
						//array
						ret += fmt.Sprintf(`
						sz += len(o.%s) * %d`, f.Name.String, sz)
						//n is length of array
						//length + array
						n += IntSize
					} else {
						n += sz
					}

				} else {
					if repeated {
						n += IntSize
						ret += fmt.Sprintf(`
						for _, v := range o.%s {
							sz = sz + len(v) + %d
						}
						`, f.Name.String, IntSize)
					} else {
						ret += fmt.Sprintf(`
					sz += len(o.%s)
					`, f.Name.String)
						n += IntSize
					}
				}
			}
			ret += fmt.Sprintf(`sz += %d`, n)
			return ret
		},
		"StructFieldMarshal": func(fields []parser.StructField) string {
			var ret string
			for _, f := range fields {
				opt := getOption("repeated", f.Options)
				repeated := isBool(opt)
				if f.Type.Type == nil {
					if repeated {
						ret += fmt.Sprintf(`if n, err = o.MarshalInt(len(o.%s), data[offset:]); err != nil {
						return 0, err
						}
						offset += n
						for _, v := range o.%s {
						if n, err = v.MarshalTo(data[offset:]); err != nil {
							return 0, err
						}
						offset += n
					}
					`, f.Name.String, f.Name.String)
					} else {
						ret += fmt.Sprintf(`if n, err = o.%s.MarshalTo(data); err != nil {
						return 0, err
					}
					offset += n
					`, f.Name.String)
					}
					continue
				}
				if v, ok := typeToString[*f.Type.Type]; ok {
					if repeated {
						ret += fmt.Sprintf(`if n, err = o.MarshalInt(len(o.%s), data[offset:]); err != nil {
						return 0, err
						}
						offset += n
						for _, v := range o.%s {
						if n, err = o.Marshal%s(v, data[offset:]); err != nil {
							return 0, err
						}
						offset += n
					}
						`, f.Name.String, f.Name.String, v)
					} else {
						ret += fmt.Sprintf(`if n, err = o.Marshal%s(o.%s, data[offset:]); err != nil {
						return 0, err
					}
					offset += n
		`, v, f.Name.String)
					}
				} else {
					panic("unknown type")
				}
			}
			return ret
		},
		"StructFieldUnmarshal": func(fields []parser.StructField) string {
			var ret string
			for _, f := range fields {
				opt := getOption("repeated", f.Options)
				repeated := isBool(opt)
				if f.Type.Type == nil {
					if repeated {
						ret += fmt.Sprintf(`if l, i, err = o.UnmarshalInt(data[n:]); err != nil {
					return  0, err
				}
				n += i
				if l > 0 {
				o.%s = make([]*%s, l)
				for j := range o.%s {
					o.%s[j] = new(%s)
					if i, err = o.%s[j].UnmarshalTo(data[n:]); err != nil {
						return 0, err
					}
					n += i
				}
			}
				`, f.Name.String, UpperFirst(*f.Type.Reference), f.Name.String, f.Name.String, UpperFirst(*f.Type.Reference), f.Name.String)
					} else {
						ret += fmt.Sprintf(`if i, err = o.%s.UnmarshalTo(data[n:]); err != nil {
					return 0, err
				}
				n += i
				`, f.Name.String)
					}
					continue
				}
				if v, ok := typeToString[*f.Type.Type]; ok {
					if repeated {
						ret += fmt.Sprintf(`if l, i, err = o.UnmarshalInt(data[n:]); err != nil {
					return  0, err
				}
				n += i
				if l > 0 {
				o.%s = make([]%s, l)
				for j := range o.%s {
					if v, m, err := o.Unmarshal%s(data[n:]); err != nil {
						return  0, err
					}else{
						i = m
						o.%s[j] = v
					}
					n += i
				}
			}
				`, f.Name.String, f.Type.Type.GoString(), f.Name.String, v, f.Name.String)
					} else {
						ret += fmt.Sprintf(`if o.%s, i, err = o.Unmarshal%s(data[n:]); err != nil {
					return 0, err
				}
				n += i
				`, f.Name.String, v)
					}
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
			return UpperFirst(*s)
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
)

func getOption(name string, fields []*parser.StructOption) *parser.Literal {
	if fields == nil {
		return nil
	}
	for _, f := range fields {
		if f.Name == name {
			return &f.Value
		}
	}
	return nil
}

func UpperFirst(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func isBool(opt *parser.Literal) bool {
	if opt == nil {
		return false
	}
	value := GenerateLiteral(*opt)
	return value == "true"
}

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

	// UnmarshalTo reads a wire-format message from data.
func (o *{{$parent.Name.String}}) UnmarshalTo(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, fmt.Errorf("invalid data size %d", len(data))
	}
	*o = {{$parent.Name.String}}(uint16(data[0]) | uint16(data[1])<<8)
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
	_, err := o.UnmarshalTo(data)
	return err
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
	{{.Name.String}}{{with .Options}}{{if . | StructFieldIsRepeat}}[]{{end}}{{end}}*{{GetString .Type.Reference}} 
{{- else}}
	{{.Name.String}} {{with .Options}}{{if . | StructFieldIsRepeat}}[]{{end}}{{end}}{{.Type.Type.GoString}}
{{- end}}
{{- end}}
}

func (o *{{.Name.String}}) Size() int {
	var sz int
	{{.Fields | StructFieldLength}}
	return sz
}

func (o *{{.Name.String}}) MarshalTo(data []byte) (int, error) {
	var (
		offset, n int
		err error
	)
	{{.Fields | StructFieldMarshal}}
	return offset, nil
}

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *{{.Name.String}}) MarshalBinary() (data []byte, err error) {
	sz := o.Size()
	data = make([]byte, sz)
	n, err := o.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	if n != sz {
		return nil, fmt.Errorf("%s size / offset different %d : %d", "Marshal", sz, n)
	}
	return data, nil
}

func (o *{{.Name.String}}) UnmarshalTo(data []byte) (int, error) {
	var (
		i, n, l int
		err  error
	)
	{{.Fields | StructFieldUnmarshal}}
	_ = l
	return n, nil
}

// Unmarshal decodes data as conform encoding.BinaryUnmarshaler.
func (o *{{.Name.String}}) UnmarshalBinary(data []byte) error {
	_, err := o.UnmarshalTo(data)
	return err

}
{{- end}}
`))
}
