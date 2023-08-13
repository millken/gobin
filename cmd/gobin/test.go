package main

const (
	a  int32   = 1
	b  float32 = 1.1
	c  string  = "hello"
	d  bool    = true
	e  int64   = 1
	af int16   = 3
	f  float64 = 1
	g  int64   = 1
)

var (
	gag []byte = []byte{0x01, 0x02}
)

type A struct {
	a int32
	B struct {
		b int32
	}
}

func fff() {
	var a A
	a.a = 1

}
