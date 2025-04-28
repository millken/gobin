package main

import (
	"fmt"
	"go/format"
	"testing"
)

func TestGen(t *testing.T) {
	g := &Generator{
		GoFile:      "./demo",
		OutName:     "./demo/demo_bin.go",
		IsDir:       true,
		Types:       []string{"B"},
		StructInfos: make([]*StructInfo, 0),
	}
	// err := g.Run()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	err := g.Parse(g.GoFile, g.IsDir)
	if err != nil {
		t.Fatal(err)
	}
	code, err := g.GenerateUnmarshal()
	if err != nil {
		t.Fatal(err)
	}
	code, err = format.Source(code)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s", code)
}
