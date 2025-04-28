package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"text/template"
)

// FieldType, StructField, and StructInfo definitions remain the same as in the previous example

const sizeTmpl = `
// SizeBinary returns the size of the serialized object
func (s *{{.Name}}) SizeBinary() int {
	size := 0
	{{- range .Fields}}
	//{{.Name}} size
	{{template "sizeField" .}}
	{{- end -}}
	return size
}

{{define "sizeField"}}
	{{- if eq .Type.Kind "basic" -}}
	size += {{ BasicSize .Type.Name }}
	{{- if or (eq .Type.Name "string") (eq .Type.Name "[]byte")}}
		size += len(s.{{.Name}})
	{{- end}}
	{{else if eq .Type.Kind "slice" -}}
	size += 8// int for length
	for _, v := range s.{{.Name}} {
		{{- template "sizeValue" .Type.ElemType -}}
	}
	{{else if eq .Type.Kind "array"}}
	for _, v := range s.{{.Name}} {
		{{template "sizeValue" .Type.ElemType}}
	}
	{{else if eq .Type.Kind "map"}}
	size += 8 // int for length
	for k, v := range s.{{.Name}} {
		{{template "sizeValue" .Type.KeyType}}
		// size += {{.Type.KeyType.Name}}Size1
		{{template "sizeValue" .Type.ElemType}}
	}
	{{else if eq .Type.Kind "pointer"}}
	size += 1 // uint8 for nil check
	if s.{{.Name}} != nil {
		{{template "sizeValue" .Type.ElemType}}
	}
	{{else if eq .Type.Kind "struct"}}
	size += s.{{.Name}}.SizeBinary()
	{{end}}
{{end}}

{{define "sizeValue"}}
	{{if eq .Kind "basic"}}
	size += {{ BasicSize .Name }}
	{{else if eq .Kind "slice"}}
	size += 4 // uint32 for length
	for _, elem := range v {
		{{template "sizeValue" .ElemType}}
	}
	{{else if eq .Kind "array"}}
	for _, elem := range v {
		{{template "sizeValue" .ElemType}}
	}
	{{else if eq .Kind "map"}}
	size += 4 // uint32 for length
	for k, elem := range v {
		
		// size += {{.KeyType.Name}}   Size4
		{{template "sizeValue" .ElemType}}
	}
	{{else if eq .Kind "pointer"}}
	size += 1 // uint8 for nil check
	if v != nil {
		{{template "sizeValue" .ElemType}}
	}
	{{else if eq .Kind "struct"}}
	size += v.SizeBinary()
	{{end}}
{{end}}
`

const marshalTmpl = `
// MarshalBinary implements the encoding.BinaryMarshaler interface
func (s *{{.Name}}) MarshalBinary() ([]byte, error) {
	size := s.SizeBinary()
	buf := make([]byte, size)
	_, err := s.MarshalBinaryTo(buf)
	return buf, err
}

// MarshalBinaryTo marshals the object into a pre-allocated byte slice
func (s *{{.Name}}) MarshalBinaryTo(buf []byte) (int, error) {
	var offset int
	{{range .Fields}}
	{{template "marshalField" dict "Field" . "Parent" "s" "Offset" "offset"}}
	{{end}}
	return offset, nil
}

{{define "marshalField"}}
	{{if eq .Field.Type.Kind "basic"}}
	if n, err := Marshal{{.Field.Type.Name}}({{.Parent}}.{{.Field.Name}}, buf[{{.Offset}}:]); err != nil {
		return 0, err
	} else {
		{{.Offset}} += n
	}
	{{else if eq .Field.Type.Kind "slice"}}
	if n, err := binary.PutUvarint(buf[{{.Offset}}:], uint64(len({{.Parent}}.{{.Field.Name}}))); err != nil {
		return 0, err
	} else {
		{{.Offset}} += n
	}
	for _, v := range {{.Parent}}.{{.Field.Name}} {
		{{template "marshalValue" dict "Type" .Field.Type.ElemType "Value" "v" "Offset" .Offset}}
	}
	{{else if eq .Field.Type.Kind "array"}}
	for _, v := range {{.Parent}}.{{.Field.Name}} {
		{{template "marshalValue" dict "Type" .Field.Type.ElemType "Value" "v" "Offset" .Offset}}
	}
	{{else if eq .Field.Type.Kind "map"}}
	if n, err := binary.PutUvarint(buf[{{.Offset}}:], uint64(len({{.Parent}}.{{.Field.Name}}))); err != nil {
		return 0, err
	} else {
		{{.Offset}} += n
	}
	for k, v := range {{.Parent}}.{{.Field.Name}} {
		if n, err := Marshal{{.Field.Type.KeyType.Name}}(k, buf[{{.Offset}}:]); err != nil {
			return 0, err
		} else {
			{{.Offset}} += n
		}
		{{template "marshalValue" dict "Type" .Field.Type.ElemType "Value" "v" "Offset" .Offset}}
	}
	{{else if eq .Field.Type.Kind "pointer"}}
	if {{.Parent}}.{{.Field.Name}} == nil {
		buf[{{.Offset}}] = 0
		{{.Offset}}++
	} else {
		buf[{{.Offset}}] = 1
		{{.Offset}}++
		{{template "marshalValue" dict "Type" .Field.Type.ElemType "Value" (printf "*%s.%s" .Parent .Field.Name) "Offset" .Offset}}
	}
	{{else if eq .Field.Type.Kind "struct"}}
	if n, err := {{.Parent}}.{{.Field.Name}}.MarshalBinaryTo(buf[{{.Offset}}:]); err != nil {
		return 0, err
	} else {
		{{.Offset}} += n
	}
	{{end}}
{{end}}

{{define "marshalValue"}}
	{{if eq .Type.Kind "basic"}}
	if n, err := Marshal{{.Type.Name}}({{.Value}}, buf[{{.Offset}}:]); err != nil {
		return 0, err
	} else {
		{{.Offset}} += n
	}
	{{else if eq .Type.Kind "slice"}}
	if n, err := binary.PutUvarint(buf[{{.Offset}}:], uint64(len({{.Value}}))); err != nil {
		return 0, err
	} else {
		{{.Offset}} += n
	}
	for _, elem := range {{.Value}} {
		{{template "marshalValue" dict "Type" .Type.ElemType "Value" "elem" "Offset" .Offset}}
	}
	{{else if eq .Type.Kind "array"}}
	for _, elem := range {{.Value}} {
		{{template "marshalValue" dict "Type" .Type.ElemType "Value" "elem" "Offset" .Offset}}
	}
	{{else if eq .Type.Kind "map"}}
	if n, err := binary.PutUvarint(buf[{{.Offset}}:], uint64(len({{.Value}}))); err != nil {
		return 0, err
	} else {
		{{.Offset}} += n
	}
	for k, elem := range {{.Value}} {
		if n, err := Marshal{{.Type.KeyType.Name}}(k, buf[{{.Offset}}:]); err != nil {
			return 0, err
		} else {
			{{.Offset}} += n
		}
		{{template "marshalValue" dict "Type" .Type.ElemType "Value" "elem" "Offset" .Offset}}
	}
	{{else if eq .Type.Kind "pointer"}}
	if {{.Value}} == nil {
		buf[{{.Offset}}] = 0
		{{.Offset}}++
	} else {
		buf[{{.Offset}}] = 1
		{{.Offset}}++
		{{template "marshalValue" dict "Type" .Type.ElemType "Value" (printf "*%s" .Value) "Offset" .Offset}}
	}
	{{else if eq .Type.Kind "struct"}}
	if n, err := {{.Value}}.MarshalBinaryTo(buf[{{.Offset}}:]); err != nil {
		return 0, err
	} else {
		{{.Offset}} += n
	}
	{{end}}
{{end}}
`

const unmarshalTmpl = `
// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (s *{{.Name}}) UnmarshalBinary(data []byte) error {
	_, err := s.UnmarshalBinaryFrom(data)
	return err
}

// UnmarshalBinaryFrom unmarshals the object from a byte slice
func (s *{{.Name}}) UnmarshalBinaryFrom(data []byte) (int, error) {
	var offset int
	{{range .Fields}}
	{{template "unmarshalField" dict "Field" . "Parent" "s" "Offset" "offset"}}
	{{end}}
	return offset, nil
}

{{define "unmarshalField"}}
	{{if eq .Field.Type.Kind "basic"}}
	if v, n, err := Unmarshal{{.Field.Type.Name}}(data[{{.Offset}}:]); err != nil {
		return 0, err
	} else {
		{{.Parent}}.{{.Field.Name}} = v
		{{.Offset}} += n
	}
	{{else if eq .Field.Type.Kind "slice"}}
	if length, n := binary.Uvarint(data[{{.Offset}}:]); n <= 0 {
		return 0, fmt.Errorf("error reading length")
	} else {
		{{.Offset}} += n
		{{.Parent}}.{{.Field.Name}} = make([]{{template "typeString" .Field.Type.ElemType}}, length)
		for i := range {{.Parent}}.{{.Field.Name}} {
			{{template "unmarshalValue" dict "Type" .Field.Type.ElemType "Parent" .Parent "Field" .Field "Index" "i" "Offset" .Offset}}
		}
	}
	{{else if eq .Field.Type.Kind "array"}}
	for i := range {{.Parent}}.{{.Field.Name}} {
		{{template "unmarshalValue" dict "Type" .Field.Type.ElemType "Parent" .Parent "Field" .Field "Index" "i" "Offset" .Offset}}
	}
	{{else if eq .Field.Type.Kind "map"}}
	if length, n := binary.Uvarint(data[{{.Offset}}:]); n <= 0 {
		return 0, fmt.Errorf("error reading length")
	} else {
		{{.Offset}} += n
		{{.Parent}}.{{.Field.Name}} = make(map[{{.Field.Type.KeyType.Name}}]{{template "typeString" .Field.Type.ElemType}}, length)
		for i := uint64(0); i < length; i++ {
			var k {{.Field.Type.KeyType.Name}}
			if v, n, err := Unmarshal{{.Field.Type.KeyType.Name}}(data[{{.Offset}}:]); err != nil {
				return 0, err
			} else {
				k = v
				{{.Offset}} += n
			}
			{{template "unmarshalValue" dict "Type" .Field.Type.ElemType "Parent" .Parent "Field" .Field "Index" "k" "Offset" .Offset}}
		}
	}
	{{else if eq .Field.Type.Kind "pointer"}}
	if data[{{.Offset}}] == 1 {
		{{.Offset}}++
		var v {{template "typeString" .Field.Type.ElemType}}
		{{template "unmarshalValue" dict "Type" .Field.Type.ElemType "Parent" "&v" "Field" .Field "Offset" .Offset}}
		{{.Parent}}.{{.Field.Name}} = &v
	} else {
		{{.Offset}}++
		{{.Parent}}.{{.Field.Name}} = nil
	}
	{{else if eq .Field.Type.Kind "struct"}}
	if n, err := {{.Parent}}.{{.Field.Name}}.UnmarshalBinaryFrom(data[{{.Offset}}:]); err != nil {
		return 0, err
	} else {
		{{.Offset}} += n
	}
	{{end}}
{{end}}

{{define "unmarshalValue"}}
	{{if eq .Type.Kind "basic"}}
	if v, n, err := Unmarshal{{.Type.Name}}(data[{{.Offset}}:]); err != nil {
		return 0, err
	} else {
		{{.Parent}}.{{.Field.Name}}{{if .Index}}[{{.Index}}]{{end}} = v
		{{.Offset}} += n
	}
	{{else if eq .Type.Kind "slice"}}
	if length, n := binary.Uvarint(data[{{.Offset}}:]); n <= 0 {
		return 0, fmt.Errorf("error reading length")
	} else {
		{{.Offset}} += n
		{{.Parent}}.{{.Field.Name}}{{if .Index}}[{{.Index}}]{{end}} = make([]{{template "typeString" .Type.ElemType}}, length)
		for i := range {{.Parent}}.{{.Field.Name}}{{if .Index}}[{{.Index}}]{{end}} {
			{{template "unmarshalValue" dict "Type" .Type.ElemType "Parent" .Parent "Field" .Field "Index" (printf "%s][i" .Index) "Offset" .Offset}}
		}
	}
	{{else if eq .Type.Kind "array"}}
	for i := range {{.Parent}}.{{.Field.Name}}{{if .Index}}[{{.Index}}]{{end}} {
		{{template "unmarshalValue" dict "Type" .Type.ElemType "Parent" .Parent "Field" .Field "Index" (printf "%s][i" .Index) "Offset" .Offset}}
	}
	{{else if eq .Type.Kind "map"}}
	if length, n := binary.Uvarint(data[{{.Offset}}:]); n <= 0 {
		return 0, fmt.Errorf("error reading length")
	} else {
		{{.Offset}} += n
		{{.Parent}}.{{.Field.Name}}{{if .Index}}[{{.Index}}]{{end}} = make(map[{{.Type.KeyType.Name}}]{{template "typeString" .Type.ElemType}}, length)
		for i := uint64(0); i < length; i++ {
			var k {{.Type.KeyType.Name}}
			if v, n, err := Unmarshal{{.Type.KeyType.Name}}(data[{{.Offset}}:]); err != nil {
				return 0, err
			} else {
				k = v
				{{.Offset}} += n
			}
			{{template "unmarshalValue" dict "Type" .Type.ElemType "Parent" .Parent "Field" .Field "Index" (printf "%s][k" .Index) "Offset" .Offset}}
		}
	}
	{{else if eq .Type.Kind "pointer"}}
	if data[{{.Offset}}] == 1 {
		{{.Offset}}++
		var v {{template "typeString" .Type.ElemType}}
		{{template "unmarshalValue" dict "Type" .Type.ElemType "Parent" "&v" "Field" .Field "Offset" .Offset}}
		{{.Parent}}.{{.Field.Name}}{{if .Index}}[{{.Index}}]{{end}} = &v
	} else {
		{{.Offset}}++
		{{.Parent}}.{{.Field.Name}}{{if .Index}}[{{.Index}}]{{end}} = nil
	}
	{{else if eq .Type.Kind "struct"}}
	if n, err := {{.Parent}}.{{.Field.Name}}{{if .Index}}[{{.Index}}]{{end}}.UnmarshalBinaryFrom(data[{{.Offset}}:]); err != nil {
		return 0, err
	} else {
		{{.Offset}} += n
	}
	{{end}}
{{end}}

{{define "typeString"}}
	{{if eq .Kind "basic"}}{{.Name}}
	{{else if eq .Kind "slice"}}[]{{template "typeString" .ElemType}}
	{{else if eq .Kind "array"}}[{{.Size}}]{{template "typeString" .ElemType}}
	{{else if eq .Kind "map"}}map[{{.KeyType.Name}}]{{template "typeString" .ElemType}}
	{{else if eq .Kind "pointer"}}*{{template "typeString" .ElemType}}
	{{else if eq .Kind "struct"}}{{.Name}}
	{{end}}
{{end}}
`

func generateMarshalUnmarshal(structInfos []*StructInfo) (string, error) {
	var buf bytes.Buffer

	buf.WriteString("package main\n\n")
	buf.WriteString("import (\n\t\"encoding/binary\"\n\t\"fmt\"\n)\n\n")

	funcMap := template.FuncMap{
		"BasicSize": func(t string) int {
			v := basicTypes.Get(t)
			if v != nil {
				return v.Size
			}
			return 4
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	}

	sizeTmplObj := template.Must(template.New("size").Funcs(funcMap).Parse(sizeTmpl))
	// marshalTmplObj := template.Must(template.New("marshal").Funcs(funcMap).Parse(marshalTmpl))
	// unmarshalTmplObj := template.Must(template.New("unmarshal").Funcs(funcMap).Parse(unmarshalTmpl))

	for _, si := range structInfos {
		if err := sizeTmplObj.Execute(&buf, si); err != nil {
			return "", err
		}
		// if err := marshalTmplObj.Execute(&buf, si); err != nil {
		// 	return "", err
		// }
		// if err := unmarshalTmplObj.Execute(&buf, si); err != nil {
		// 	return "", err
		// }
	}

	// Format the generated code
	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		return "", err
	}

	return string(formattedCode), nil
}

func main() {
	// Example StructInfo (you would typically get this from your AST parser)
	structInfos := []*StructInfo{
		{
			Name: "A",
			Fields: []StructField{
				{Name: "A1", Type: &FieldType{Kind: "basic", Name: "int"}},
				{Name: "A2", Type: &FieldType{Kind: "slice", ElemType: &FieldType{Kind: "basic", Name: "int"}}},
				{Name: "A3", Type: &FieldType{Kind: "slice", ElemType: &FieldType{Kind: "slice", ElemType: &FieldType{Kind: "basic", Name: "int"}}}},
				{Name: "A4", Type: &FieldType{Kind: "slice", ElemType: &FieldType{Kind: "slice", ElemType: &FieldType{Kind: "slice", ElemType: &FieldType{Kind: "basic", Name: "int"}}}}},
				{Name: "A5", Type: &FieldType{Kind: "map", KeyType: &FieldType{Kind: "basic", Name: "string"}, ElemType: &FieldType{Kind: "basic", Name: "int"}}},
				{Name: "A6", Type: &FieldType{Kind: "array", Size: 5, ElemType: &FieldType{Kind: "basic", Name: "int"}}},
				{Name: "A7", Type: &FieldType{Kind: "array", Size: 3, ElemType: &FieldType{Kind: "array", Size: 6, ElemType: &FieldType{Kind: "basic", Name: "int"}}}},
				{Name: "A8", Type: &FieldType{Kind: "basic", Name: "[]byte"}},
				{Name: "A9", Type: &FieldType{Kind: "pointer", ElemType: &FieldType{Kind: "basic", Name: "int"}}},
				{Name: "A10", Type: &FieldType{Kind: "pointer", ElemType: &FieldType{Kind: "pointer", ElemType: &FieldType{Kind: "basic", Name: "string"}}}},
				{Name: "B1", Type: &FieldType{Kind: "basic", Name: "B1"}},
				{Name: "B2", Type: &FieldType{Kind: "pointer", ElemType: &FieldType{Kind: "struct", Name: "B2"}}},
			},
		},
	}

	code, err := generateMarshalUnmarshal(structInfos)
	if err != nil {
		fmt.Printf("Error generating code: %v\n", err)
		return
	}

	fmt.Println(code)

	// Optionally, write to a file
	if err := os.WriteFile("generated_marshal.go", []byte(code), 0644); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
}

type baseType struct {
	Name string
	Size int
	Type string
}

type baseTypes []baseType

var basicTypes = baseTypes{
	{"int64", 8, "Int64"},
	{"uint64", 8, "Uint64"},
	{"float64", 8, "Float64"},
	{"int", 8, "Int"},
	{"uint", 8, "Uint"},
	{"int32", 4, "Int32"},
	{"uint32", 4, "Uint32"},
	{"float32", 4, "Float32"},
	{"bool", 1, "Bool"},
	{"int8", 1, "Int8"},
	{"uint8", 1, "Uint8"},
	{"byte", 1, "Byte"},
	{"int16", 2, "Int16"},
	{"uint16", 2, "Uint16"},
	{"string", 8, "String"},
	{"[]byte", 9, "Bytes"},
}

func (b baseTypes) Get(t string) *baseType {
	for _, v := range b {
		if v.Name == t {
			return &v
		}
	}
	return nil
}
