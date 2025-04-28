package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strconv"
	"strings"
)

func readSource(filename string, src any) ([]byte, error) {
	if src != nil {
		switch s := src.(type) {
		case string:
			return []byte(s), nil
		case []byte:
			return s, nil
		case *bytes.Buffer:
			// is io.Reader, but src is already available in []byte form
			if s != nil {
				return s.Bytes(), nil
			}
		case io.Reader:
			return io.ReadAll(s)
		}
		return nil, errors.New("invalid source")
	}
	return os.ReadFile(filename)
}

type FieldType struct {
	Name     string
	Kind     string
	KeyType  *FieldType    // 用于 map 的 key 类型
	ElemType *FieldType    // 用于 slice、array、map、pointer 和自定义类型的元素类型
	Size     int           // 用于固定大小数组
	Level    int           // loop level
	Fields   []StructField // 用于嵌套结构体
	Expr     ast.Expr
}

type StructField struct {
	Name string
	Type *FieldType
}

type StructInfo struct {
	Name   string
	Fields []StructField
}

func ParseFiles(paths []string) ([]*StructInfo, error) {
	var structInfos []*StructInfo
	for _, p := range paths {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, p, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		structsInfo, err := parseAstStructs(file)
		if err != nil {
			return nil, err
		}

		structInfos = append(structInfos, structsInfo...)
	}

	return structInfos, nil
}

func parseAstStructs(node *ast.File) ([]*StructInfo, error) {
	var structInfos []*StructInfo
	typeSpecs := make(map[string]ast.Expr)

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			typeSpecs[typeSpec.Name.Name] = typeSpec.Type

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			structInfo := &StructInfo{
				Name:   typeSpec.Name.Name,
				Fields: make([]StructField, 0, len(structType.Fields.List)),
			}

			for _, field := range structType.Fields.List {
				for _, name := range field.Names {
					structInfo.Fields = append(structInfo.Fields, StructField{
						Name: name.Name,
						Type: parseFieldType(field.Type, typeSpecs, 0),
					})
				}
			}

			structInfos = append(structInfos, structInfo)
		}
	}

	return structInfos, nil
}

func parseFieldType(expr ast.Expr, typeSpecs map[string]ast.Expr, level int) *FieldType {
	switch t := expr.(type) {
	case *ast.Ident:
		if typeSpec, ok := typeSpecs[t.Name]; ok {
			return parseFieldType(typeSpec, typeSpecs, level+1)
		}
		return &FieldType{Name: t.Name, Kind: "basic", Level: level}
	case *ast.ArrayType:
		if t.Len == nil {
			if ident, ok := t.Elt.(*ast.Ident); ok && ident.Name == "byte" {
				return &FieldType{Name: "[]byte", Kind: "basic", Level: level}
			}
			return &FieldType{
				Name:     fmt.Sprintf("[]%s", parseFieldType(t.Elt, typeSpecs, level).Name),
				Kind:     "slice",
				ElemType: parseFieldType(t.Elt, typeSpecs, level+1),
				Level:    level,
				Expr:     expr,
			}
		} else {
			size, _ := strconv.Atoi(t.Len.(*ast.BasicLit).Value)
			return &FieldType{
				Kind:     "array",
				Size:     size,
				ElemType: parseFieldType(t.Elt, typeSpecs, level+1),
				Level:    level,
				Expr:     expr,
			}
		}
	case *ast.MapType:
		return &FieldType{
			Name:     fmt.Sprintf("map[%s]%s", parseFieldType(t.Key, typeSpecs, level).Name, parseFieldType(t.Value, typeSpecs, level).Name),
			Kind:     "map",
			KeyType:  parseFieldType(t.Key, typeSpecs, level+1),
			ElemType: parseFieldType(t.Value, typeSpecs, level+1),
			Level:    level,
			Expr:     expr,
		}
	case *ast.StarExpr:
		return &FieldType{
			Kind:     "pointer",
			ElemType: parseFieldType(t.X, typeSpecs, level+1),
			Level:    level,
			Expr:     expr,
		}
	case *ast.StructType:
		fields := make([]StructField, 0, len(t.Fields.List))
		for _, field := range t.Fields.List {
			for _, name := range field.Names {
				fields = append(fields, StructField{
					Name: name.Name,
					Type: parseFieldType(field.Type, typeSpecs, level),
				})
			}
		}
		return &FieldType{
			Kind:   "struct",
			Fields: fields,
			Level:  level,
			Expr:   expr,
		}
	default:
		return &FieldType{Name: fmt.Sprintf("%T", t), Kind: "unknown", Level: level}
	}
}

func printStructInfo(info *StructInfo, depth int) {
	indent := strings.Repeat("  ", depth)
	fmt.Printf("%sStruct: %s\n", indent, info.Name)
	for _, field := range info.Fields {
		fmt.Printf("%s  Field: %s\n", indent, field.Name)
		printFieldType(field.Type, depth+2)
	}
}

func printFieldType(ft *FieldType, depth int) {
	indent := strings.Repeat("  ", depth)
	fmt.Printf("%sType: %s\n", indent, ft.Kind)
	if ft.Name != "" {
		fmt.Printf("%sName: %s\n", indent, ft.Name)
	}
	if ft.Size > 0 {
		fmt.Printf("%sSize: %d\n", indent, ft.Size)
	}
	if ft.KeyType != nil {
		fmt.Printf("%sKey Type:\n", indent)
		printFieldType(ft.KeyType, depth+1)
	}
	if ft.ElemType != nil {
		fmt.Printf("%sElement Type:\n", indent)
		printFieldType(ft.ElemType, depth+1)
	}
	if len(ft.Fields) > 0 {
		fmt.Printf("%sFields:\n", indent)
		for _, field := range ft.Fields {
			fmt.Printf("%s  Field: %s\n", indent, field.Name)
			printFieldType(field.Type, depth+2)
		}
	}
}
