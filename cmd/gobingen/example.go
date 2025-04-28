package main

import (
	"github.com/millken/gobin"
)

type A struct {
	A1 int
	A2 []int
	A3 [][]int
	A4 [][][]int
}

type Bstruct struct {
	ID int `gobin:"id"`
}

//go:generate gobingen
type AAA struct {
	gobin.Safe
	A1 string
	A2 int
	A3 []float32
	A4 float64
	// ccc  string
	// B0   interface{}
	// B2   any
	// A1 []string
	// A2 []int
	// a3   map[string]string
	// a4   *string
	// a5   Bstruct
	// a6   *Bstruct
	_   int
	Sub struct {
		sid  int `gobin:"id"`
		b2   []byte
		Name string `gobin:"name"`
		Sub3 struct {
			Id int
			B3 []byte
		}
	}
}
