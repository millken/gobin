package main

import (
	"bytes"
	"fmt"
	"gobin/parser"
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestParser(t *testing.T) {
	out := &bytes.Buffer{}
	src := `
	package example

	option go_marshal = "unsafe"
	option go_int = 3

	//comment int32 a = 1
	const int32 a = 1
	/* 
		comment 3
		coment 30 
	*/
	const float b = 1.1
	const string c = "hello"
	const bool d = true
	const int64 e = 1
	const double f = 1.0
	const int64 g = 1
	struct Person {
		string Name
		uint64 BirthDay
		bytes Phone
		int32 Siblings
		bool Spouse

		double Money
	}

	// SearchRequest is a request message for SearchService.Search.
	struct SearchRequest {
		// query is a search query string.
		string query
		Person person
		// page_number is a page number.
		int32 page_number
		// result_per_page is a result per page.
		int32 result_per_page
	  }
	  
	  struct SearchResponse {
		string results
	  }
	`
	p, err := NewParser(out, src)
	assert.NoError(t, err)
	assert.NoError(t, p.Parse())
	fmt.Println(out.String())
}

func TestStructField(t *testing.T) {
	var ref string
	ref = "ref"
	var t1 parser.Type
	t1 = parser.Int32

	fields := []parser.StructField{
		{
			Type: &parser.StructType{
				Type: &t1,
			},
			Name: parser.Name{String: "noneopt"},
		},
		{
			Type: &parser.StructType{
				Type: &t1,
			},
			Name: parser.Name{String: "arrint32"},
			Options: []*parser.StructOption{
				{
					Name:  "repeated",
					Value: parser.LiteralBool{Value: true},
				},
			},
		},
		{
			Type: &parser.StructType{
				Type:      nil,
				Reference: &ref,
			},
			Name: parser.Name{String: "fname"},
			Options: []*parser.StructOption{
				{
					Name:  "repeated",
					Value: parser.LiteralBool{Value: true},
				},
			},
		},
	}
	// var n int
	// var ret string
	// for _, f := range fields {
	// 	opt := getOption("repeated", f.Options)
	// 	repeated := isBool(opt)
	// 	if f.Type.Type == nil {
	// 		if repeated {
	// 			//array
	// 			ret += fmt.Sprintf(`len(o.%s) * o.%s.Size() + `, f.Name.String, f.Name.String)
	// 			//n is length of array
	// 			//length + array
	// 			n += strconv.IntSize / 8
	// 		} else {
	// 			// reference to another struct
	// 			ret += fmt.Sprintf(`o.%s.Size() + `, f.Name.String)
	// 		}
	// 		continue
	// 	}
	// 	if sz := f.Type.Type.Size(); sz > 0 {
	// 		if repeated {
	// 			//array
	// 			ret += fmt.Sprintf(`len(o.%s) * %d + `, f.Name.String, sz)
	// 			//n is length of array
	// 			//length + array
	// 			n += strconv.IntSize / 8
	// 		} else {
	// 			n += sz
	// 		}

	// 	} else {
	// 		ret += fmt.Sprintf(`len(o.%s) + `, f.Name.String)
	// 		n += strconv.IntSize / 8
	// 	}
	// }
	// t.Log(ret + fmt.Sprintf("%d", n))

	// 	var ret string
	// 	for _, f := range fields {
	// 		opt := getOption("repeated", f.Options)
	// 		repeated := isBool(opt)
	// 		if f.Type.Type == nil {
	// 			if repeated {
	// 				ret += fmt.Sprintf(`if n, err = o.MarshalInt(len(o.%s), data[offset:]); err != nil {
	// 				return nil, err
	// 				}
	// 				offset += n
	// 				for _, v := range o.%s {
	// 				if n, err = v.MarshalTo(data[offset:]); err != nil {
	// 					return nil, err
	// 				}
	// 				offset += n
	// 			}
	// 			`, f.Name.String, f.Name.String)
	// 			} else {
	// 				ret += fmt.Sprintf(`if n, err = o.%s.MarshalTo(data); err != nil {
	// 				return nil, err
	// 			}
	// 			offset += n
	// 			`, f.Name.String)
	// 			}
	// 			continue
	// 		}
	// 		if v, ok := typeToString[*f.Type.Type]; ok {
	// 			if repeated {
	// 				ret += fmt.Sprintf(`if n, err = o.MarshalInt(len(o.%s), data[offset:]); err != nil {
	// 				return nil, err
	// 				}
	// 				offset += n
	// 				for _, v := range o.%s {
	// 				if n, err = o.Marshal%s(v, data[offset:]); err != nil {
	// 					return nil, err
	// 				}
	// 				offset += n
	// 			}
	// 				`, f.Name.String, f.Name.String, v)
	// 			} else {
	// 				ret += fmt.Sprintf(`if n, err = o.Marshal%s(o.%s, data[offset:]); err != nil {
	// 				return nil, err
	// 			}
	// 			offset += n
	// `, v, f.Name.String)
	// 			}
	// 		} else {
	// 			panic("unknown type")
	// 		}
	// 	}
	// 	ret += `	if offset != sz {
	// 		return nil, fmt.Errorf("%s size / offset different %d : %d", "Marshal", sz, offset)
	// 	}`

	var ret string
	for _, f := range fields {
		opt := getOption("repeated", f.Options)
		repeated := isBool(opt)
		if f.Type.Type == nil {
			if repeated {
				ret += fmt.Sprintf(`if l, i, err = o.UnmarshalInt(data[n:]); err != nil {
			return nil, err
		}
		n += i
		o.%s = make([]*%s, l)
		for j := range o.%s {
			if i, err = o.%s[j].UnmarshalTo(data[n:]); err != nil {
				return nil, err
			}
			n += i
		}
		`, f.Name.String, *f.Type.Reference, f.Name.String, f.Name.String)
			} else {
				ret += fmt.Sprintf(`if i, err = o.%s.UnmarshalTo(data[n:]); err != nil {
			return err
		}
		n += i
		`, f.Name.String)
			}
			continue
		}
		if v, ok := typeToString[*f.Type.Type]; ok {
			if repeated {
				ret += fmt.Sprintf(`if l, i, err = o.UnmarshalInt(data[n:]); err != nil {
			return  err
		}
		n += i
		o.%s = make([]%s, l)
		for j := range o.%s {
			if v, i, err := o.Unmarshal%s(data[n:]); err != nil {
				return  err
			}else{
				o.%s[j] = v
			}
			n += i
		}
		`, f.Name.String, f.Type.Type.GoString(), f.Name.String, v, f.Name.String)
			} else {
				ret += fmt.Sprintf(`if o.%s, i, err = o.Unmarshal%s(data[n:]); err != nil {
			return err
		}
		n += i
		`, f.Name.String, v)
			}
		} else {
			panic("unknown type")
		}
	}
	t.Log(ret)
}

func TestPebbleTemplate(t *testing.T) {
	input, err := os.ReadFile("./testdata/pebble.gobin")
	assert.NoError(t, err)

	out := &bytes.Buffer{}
	var opt []option
	opt = []option{WithFormatted()}
	p, err := NewParser(out, input, opt...)
	assert.NoError(t, err)
	assert.NoError(t, p.Parse())
	fmt.Println(out.String())
}

func TestSimpleTemplate(t *testing.T) {
	input, err := os.ReadFile("./testdata/simple.gobin")
	assert.NoError(t, err)

	out := &bytes.Buffer{}
	var opt []option
	//opt = []option{WithFormatted()}
	p, err := NewParser(out, input, opt...)
	assert.NoError(t, err)
	assert.NoError(t, p.Parse())
	fmt.Println(out.String())
}
