package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"os"
	"reflect"
	"sort"
	"strings"
	"unicode"
)

type StructField struct {
	Parent string
	Name   string
	Type   string
}

type Generator struct {
	GoFile           string
	IsDir            bool
	PkgPath, PkgName string
	Types            []string
	Structs          map[string][]StructField //保存需要生成代码的结构体

	IncludePrivate           bool
	NoStdMarshalers          bool
	SnakeCase                bool
	LowerCamelCase           bool
	OmitEmpty                bool
	DisallowUnknownFields    bool
	SkipMemberNameUnescaping bool

	out           *bytes.Buffer
	size          *bytes.Buffer
	marshaler     *bytes.Buffer
	unmarshaler   *bytes.Buffer
	OutName       string
	GenBuildFlags string

	StubsOnly   bool
	LeaveTemps  bool
	NoFormat    bool
	SimpleBytes bool
}

func (g *Generator) Parse(fname string, src any) error {
	fmt.Fprintf(g.out, "// Code generated by \"gobingen %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
	fmt.Fprintln(g.out, "package ", g.PkgName)

	if len(g.Types) > 0 {
		fmt.Fprintln(g.out)
		fmt.Fprintln(g.out, "import (")
		fmt.Fprintln(g.out, `  "fmt"`)
		fmt.Fprintln(g.out, `  "github.com/millken/gobin"`)
		fmt.Fprintln(g.out, ")")
	}
	g.size = &bytes.Buffer{}
	g.marshaler = &bytes.Buffer{}
	g.unmarshaler = &bytes.Buffer{}
	g.Structs = make(map[string][]StructField)
	fset := token.NewFileSet()
	p, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	// 遍历文件中的所有声明
	for _, decl := range p.Decls {
		// 检查声明是否为 GenDecl（通用声明）
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			// 遍历 GenDecl 中的所有 Specs
			for _, spec := range genDecl.Specs {
				// 检查 Spec 是否为 TypeSpec（类型声明）
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					// 检查 TypeSpec 是否为 StructType（结构体类型）
					if _, ok := typeSpec.Type.(*ast.StructType); ok {
						// 打印结构体名称
						g.parseStruct(typeSpec)
					}
				}
			}
		}
	}
	return nil
}
func (g *Generator) parseStruct(t *ast.TypeSpec) {
	for _, name := range g.Types {
		if name == t.Name.Name {
			g.Structs[t.Name.Name] = make([]StructField, 0)
			slog.Info("parse struct", "type", t.Name.Name)
			g.processFields(t.Name.Name, "", t.Type.(*ast.StructType).Fields, 1)
			if len(g.Structs[t.Name.Name]) == 0 {
				delete(g.Structs, t.Name.Name)
			} else {

			}
		}
	}
}

func (g *Generator) processFields(structName string, parentName string, fields *ast.FieldList, indentLevel int) {
	indent := ""
	for i := 0; i < indentLevel; i++ {
		indent += "  "
	}
	for _, field := range fields.List {
		// 获取字段名称和类型
		fieldName := ""
		if len(field.Names) > 0 {
			fieldName = field.Names[0].Name

		}
		//检查是否为匿名字段
		if fieldName == "" {
			//跳过gobin
			if sel, ok := field.Type.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "gobin" {
					if sel.Sel.Name == "Safe" || sel.Sel.Name == "Unsafe" {
						continue
					}
				}
			}
			if ident, ok := field.Type.(*ast.Ident); ok {
				fieldName = ident.Name
			}
		}
		//是否为私有字段
		if !g.IncludePrivate && len(fieldName) != 0 && !unicode.IsUpper(rune(fieldName[0])) {
			continue // skip private fields
		}
		//检查字段类型是否为结构体
		//检查字段类型是否为数组
		//检查字段类型是否为切片
		//检查字段类型是否为指针
		//检查字段类型是否为基本类型
		fieldType := fmt.Sprintf("%v", field.Type)

		switch t := field.Type.(type) {
		case *ast.StructType:
			g.processFields(structName, strings.TrimLeft(parentName+"."+fieldName, "."), t.Fields, indentLevel+1)
			return
		case *ast.ArrayType: //固定长度数组
			fieldType = fmt.Sprintf("[]%v", t.Elt)
			g.processArrayField(structName, strings.TrimLeft(parentName+"."+fieldName, "."), t, fieldName, indentLevel)
		case *ast.SliceExpr: //动态长度切片
			fieldType = fmt.Sprintf("[]%v", t.X)
		case *ast.StarExpr: //Pointer
			fieldType = fmt.Sprintf("*%v", t.X)
		case *ast.Ident:
			if t.Obj == nil {
				fmt.Printf("%sField: %s is a basic type\n", indent, fieldName)
			} else {
				fmt.Printf("%sField: %s is a named type\n", indent, fieldName)
			}
		case *ast.MapType:
			keyType := fmt.Sprintf("%v", t.Key)
			valueType := fmt.Sprintf("%v", t.Value)
			fieldType = fmt.Sprintf("map[%v]%v", keyType, valueType)
			fmt.Printf("%sField: %s is a map type\n", indent, fieldName)
		case *ast.InterfaceType, *ast.ChanType, *ast.FuncType:
			continue //skip
		default:
			fmt.Printf("%sField: %s is of unknown type\n", indent, fieldName)

		}
		// fmt.Printf("%sField: %s %s\n", indent, fieldName, fieldType)
		g.processField(structName, strings.TrimLeft(parentName, "."), field, fieldType, fieldName, indentLevel)
	}
}

func (g *Generator) processArrayField(structName string, parentName string, t *ast.ArrayType, n string, indentLevel int) {
	if t.Len == nil {
		return
	}
	if structType, ok := t.Elt.(*ast.StructType); ok {
		g.processFields(structName, parentName+"."+n, structType.Fields, indentLevel+1)
		return
	}
	g.Structs[structName] = append(g.Structs[structName], StructField{Parent: parentName, Name: n, Type: fmt.Sprintf("[]%v", t.Elt)})
}

func (g *Generator) processField(structName string, parentName string, f *ast.Field, t, n string, indentLevel int) {
	if f.Tag != nil {
		tag := reflect.StructTag(strings.Trim(f.Tag.Value, "`"))
		if tag.Get("gobin") == "-" { // skip field
			return
		}
	}
	if n == "_" {
		return // skip anonymous fields
	}
	if structType, ok := f.Type.(*ast.StructType); ok {
		g.processFields(structName, parentName+"."+n, structType.Fields, indentLevel+1)
		return
	}
	g.Structs[structName] = append(g.Structs[structName], StructField{Parent: parentName, Name: n, Type: t})
	// fmt.Fprintf(g.out, "%s %s %v %s \n", structName, parentName, n, t)

}

func (g *Generator) Run() error {
	f, err := os.Create(g.OutName)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "// Code generated by \"gobingen %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
	fmt.Fprintln(f, "package ", g.PkgName)

	if len(g.Types) > 0 {
		fmt.Fprintln(f)
		fmt.Fprintln(f, "import (")
		fmt.Fprintln(f, `  "fmt"`)
		fmt.Fprintln(f, `  "github.com/millken/gobin"`)
		fmt.Fprintln(f, ")")
	}

	sort.Strings(g.Types)
	for _, t := range g.Types {
		fmt.Fprintln(f)
		if !g.NoStdMarshalers {
			fmt.Fprintln(f, "func (", t, ") MarshalJSON() ([]byte, error) { return nil, nil }")
			fmt.Fprintln(f, "func (*", t, ") UnmarshalJSON([]byte) error { return nil }")
		}

		fmt.Fprintln(f, "func (", t, ") MarshalEasyJSON(w *jwriter.Writer) {}")
		fmt.Fprintln(f, "func (*", t, ") UnmarshalEasyJSON(l *jlexer.Lexer) {}")
		fmt.Fprintln(f)
		fmt.Fprintln(f, "type EasyJSON_exporter_"+t+" *"+t)
	}
	return nil
}

type baseType struct {
	Name    string
	Size    int
	Type    string
	IsArray bool
	IsFixed bool
}

type baseTypes []baseType

var basicTypes = baseTypes{
	{"int64", 8, "Int64", false, true},
	{"uint64", 8, "Uint64", false, true},
	{"float64", 8, "Float64", false, true},
	{"int", 8, "Int", false, true},
	{"uint", 8, "Uint", false, true},
	{"int32", 4, "Int32", false, true},
	{"uint32", 4, "Uint32", false, true},
	{"float32", 4, "Float32", false, true},
	{"bool", 1, "Bool", false, true},
	{"int8", 1, "Int8", false, true},
	{"uint8", 1, "Uint8", false, true},
	{"byte", 1, "Byte", false, true},
	{"int16", 2, "Int16", false, true},
	{"uint16", 2, "Uint16", false, true},
	{"string", 8, "String", false, false},
	{"[]byte", 9, "Bytes", false, false},

	{"[]string", 8, "String", true, false},
	{"[]int64", 8, "Int64", true, true},
	{"[]uint64", 8, "Uint64", true, true},
	{"[]float64", 8, "Float64", true, true},

	{"[]int", 8, "Int", true, true},
	{"[]uint", 8, "Uint", true, true},

	{"[]rune", 4, "Int32", true, true},
	{"[]int32", 4, "Int32", true, true},
	{"[]uint32", 4, "Uint32", true, true},
	{"[]float32", 4, "Float32", true, true},
	{"[]bool", 1, "Bool", true, true},
	{"[]int8", 1, "Int8", true, true},
	{"[]uint8", 1, "Uint8", true, true},

	{"[]int16", 2, "Int16", true, true},
	{"[]uint16", 2, "Uint16", true, true},
}

func (b baseTypes) Get(t string) *baseType {
	for _, v := range b {
		if v.Name == t {
			return &v
		}
	}
	return nil
}

func UpperFirst(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}
