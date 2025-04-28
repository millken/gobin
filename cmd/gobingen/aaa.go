package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func t11() {
	// 解析源文件
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "a.go", nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	// 遍历 AST，找到结构体定义
	for _, f := range node.Decls {
		if genDecl, ok := f.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						structName := typeSpec.Name.Name
						generateMarshalMethods(structName, structType)
					}
				}
			}
		}
	}
}

func generateMarshalMethods(structName string, structType *ast.StructType) {
	var marshalBuf, unmarshalBuf bytes.Buffer

	// 生成 MarshalBinary 方法
	marshalBuf.WriteString(fmt.Sprintf("func (o *%s) MarshalBinary() ([]byte, error) {\n", structName))
	marshalBuf.WriteString("var buf bytes.Buffer\n")
	for _, field := range structType.Fields.List {
		fieldName := field.Names[0].Name
		fieldType := field.Type
		generateMarshalField(&marshalBuf, fieldName, fieldType, 1)
	}
	marshalBuf.WriteString("return buf.Bytes(), nil\n")
	marshalBuf.WriteString("}\n\n")

	// 生成 UnmarshalBinary 方法
	unmarshalBuf.WriteString(fmt.Sprintf("func (o *%s) UnmarshalBinary(data []byte) error {\n", structName))
	unmarshalBuf.WriteString("buf := bytes.NewReader(data)\n")
	for _, field := range structType.Fields.List {
		fieldName := field.Names[0].Name
		fieldType := field.Type
		generateUnmarshalField(&unmarshalBuf, fieldName, fieldType, 1)
	}
	unmarshalBuf.WriteString("return nil\n")
	unmarshalBuf.WriteString("}\n")

	// 将生成的代码写入文件
	file, err := os.Create("generated.go")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	file.WriteString("package main\n\n")
	file.WriteString("import (\n")
	file.WriteString("\"bytes\"\n")
	file.WriteString("\"encoding/binary\"\n")
	file.WriteString(")\n\n")
	file.WriteString(marshalBuf.String())
	file.WriteString(unmarshalBuf.String())
}

func generateMarshalField(buf *bytes.Buffer, fieldName string, fieldType ast.Expr, level int) {
	switch t := fieldType.(type) {
	case *ast.Ident:
		if t.Name == "int" {
			buf.WriteString(fmt.Sprintf("binary.Write(&buf, binary.LittleEndian, o.%s)\n", fieldName))
		} else {
			fmt.Printf("Unsupported field type: %s\n", t.Name)
		}
	case *ast.ArrayType:
		buf.WriteString(fmt.Sprintf("binary.Write(&buf, binary.LittleEndian, int32(len(o.%s)))\n", fieldName))
		buf.WriteString(fmt.Sprintf("for _, v := range o.%s {\n", fieldName))
		generateMarshalField(buf, "v", t.Elt, level+1)
		buf.WriteString("}\n")
	case *ast.SliceExpr:
		buf.WriteString(fmt.Sprintf("binary.Write(&buf, binary.LittleEndian, int32(len(o.%s)))\n", fieldName))
		buf.WriteString(fmt.Sprintf("for _, v := range o.%s {\n", fieldName))
		generateMarshalField(buf, "v", t.X, level+1)
		buf.WriteString("}\n")
	case *ast.MapType:
		buf.WriteString(fmt.Sprintf("binary.Write(&buf, binary.LittleEndian, int32(len(o.%s)))\n", fieldName))
		buf.WriteString(fmt.Sprintf("for k, v := range o.%s {\n", fieldName))
		buf.WriteString("binary.Write(&buf, binary.LittleEndian, int32(len(k)))\n")
		buf.WriteString("buf.WriteString(k)\n")
		generateMarshalField(buf, "v", t.Value, level+1)
		buf.WriteString("}\n")
	default:
		fmt.Printf("Unsupported field type: %T\n", t)
	}
}

func generateUnmarshalField(buf *bytes.Buffer, fieldName string, fieldType ast.Expr, level int) {
	switch t := fieldType.(type) {
	case *ast.Ident:
		if t.Name == "int" {
			buf.WriteString(fmt.Sprintf("binary.Read(buf, binary.LittleEndian, &o.%s)\n", fieldName))
		} else {
			fmt.Printf("Unsupported field type: %s\n", t.Name)
		}
	case *ast.ArrayType:
		buf.WriteString(fmt.Sprintf("var len%d int32\n", level))
		buf.WriteString(fmt.Sprintf("binary.Read(buf, binary.LittleEndian, &len%d)\n", level))
		buf.WriteString(fmt.Sprintf("o.%s = make([]%s, len%d)\n", fieldName, t.Elt, level))
		buf.WriteString(fmt.Sprintf("for i := 0; i < int(len%d); i++ {\n", level))
		generateUnmarshalField(buf, fmt.Sprintf("o.%s[i]", fieldName), t.Elt, level+1)
		buf.WriteString("}\n")
	case *ast.SliceExpr:
		buf.WriteString(fmt.Sprintf("var len%d int32\n", level))
		buf.WriteString(fmt.Sprintf("binary.Read(buf, binary.LittleEndian, &len%d)\n", level))
		buf.WriteString(fmt.Sprintf("o.%s = make([]%s, len%d)\n", fieldName, t.X, level))
		buf.WriteString(fmt.Sprintf("for i := 0; i < int(len%d); i++ {\n", level))
		generateUnmarshalField(buf, fmt.Sprintf("o.%s[i]", fieldName), t.X, level+1)
		buf.WriteString("}\n")
	case *ast.MapType:
		buf.WriteString(fmt.Sprintf("var len%d int32\n", level))
		buf.WriteString(fmt.Sprintf("binary.Read(buf, binary.LittleEndian, &len%d)\n", level))
		buf.WriteString(fmt.Sprintf("o.%s = make(map[string]%s, len%d)\n", fieldName, t.Value, level))
		buf.WriteString(fmt.Sprintf("for i := 0; i < int(len%d); i++ {\n", level))
		buf.WriteString("var keyLen int32\n")
		buf.WriteString("binary.Read(buf, binary.LittleEndian, &keyLen)\n")
		buf.WriteString("key := make([]byte, keyLen)\n")
		buf.WriteString("buf.Read(key)\n")
		buf.WriteString("var value int\n")
		generateUnmarshalField(buf, "value", t.Value, level+1)
		buf.WriteString(fmt.Sprintf("o.%s[string(key)] = value\n", fieldName))
		buf.WriteString("}\n")
	default:
		fmt.Printf("Unsupported field type: %T\n", t)
	}
}
