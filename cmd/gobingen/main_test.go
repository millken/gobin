package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	var buf bytes.Buffer
	g := &Generator{
		GoFile: "main.go",
		out:    &buf,
		Types:  []string{"A"},
	}
	g.Parse("example.go", nil)
	t.Log(buf.String())
	out := &bytes.Buffer{}

	for typeName, v := range g.Structs {
		fmt.Fprintf(out, "func (o *%s) Size() int {\n", typeName)
		fmt.Fprintln(out, "sz := 0")
		var n int
		for _, f := range v {
			fType := basicTypes.Get(f.Type)
			fieldName := f.Name
			if f.Parent != "" {
				fieldName = f.Parent + "." + f.Name
			}
			if fType != nil {
				if !fType.IsArray {
					n += fType.Size
					if !fType.IsFixed {
						fmt.Fprintf(out, "sz += len(o.%s)\n", fieldName)
					}
				} else if fType.IsFixed {
					n += 8
					fmt.Fprintf(out, "sz += len(o.%s) * %d\n", fieldName, fType.Size)
				} else { // slice
					n += 8
					fmt.Fprintf(out, `for _, v := range o.%s {
						sz += len(v) + %d
					}
					`, fieldName, fType.Size)
				}
			}
		}
		if n > 0 {
			fmt.Fprintf(out, "sz += %d\n", n)
		}
		fmt.Fprintln(out, "return sz")
		fmt.Fprintln(out, "}")
		fmt.Fprintln(out)
		// MarshalTo
		fmt.Fprintf(out, "func (o *%s) MarshalTo(data []byte) (int, error) {\n", typeName)
		fmt.Fprintln(out, "var offset,n int")
		fmt.Fprintln(out, "var err error")
		for _, f := range v {
			fType := basicTypes.Get(f.Type)
			fieldName := f.Name
			if f.Parent != "" {
				fieldName = f.Parent + "." + f.Name
			}
			if fType != nil {
				if !fType.IsArray {
					fmt.Fprintf(out, `if n, err = o.Marshal%s(o.%s, data[offset:]); err != nil {
						return 0, err
					}
					offset += n
					`, fType.Type, fieldName)
				} else if fType.IsFixed {
					fmt.Fprintf(out, `if n, err = o.MarshalInt(len(o.%s), data[offset:]); err != nil {
						return 0, err
					}
					offset += n
					`, fieldName)
					fmt.Fprintf(out, `for _, v := range o.%s {
						if n, err = o.Marshal%s(v, data[offset:]); err != nil {
							return 0, err
						}
						offset += n
					}
					`, fieldName, fType.Type)
				} else if fType.IsArray && !fType.IsFixed {
					fmt.Fprintf(out, `if n, err = o.MarshalInt(len(o.%s), data[offset:]); err != nil {
					return 0, err
				}
				offset += n
				for _, v := range o.%s {
					if n, err = o.Marshal%s(v, data[offset:]); err != nil {
						return 0, err
					}
					offset += n
				}
				`, fieldName, fieldName, fType.Type)
				} else {
					fmt.Fprintf(out, "n, err = o.%s.MarshalTo(data[offset:])\n", fieldName)
					fmt.Fprintln(out, "if err != nil {")
					fmt.Fprintln(out, "return 0, err")
					fmt.Fprintln(out, "}")
					fmt.Fprintln(out, "offset += n")
				}
			}
		}
		fmt.Fprintln(out, "return offset, nil")
		fmt.Fprintln(out, "}")
		fmt.Fprintln(out)

		// MarshalBinary
		fmt.Fprintf(out, `// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *%s) MarshalBinary() (data []byte, err error) {
	sz := o.Size()
	data = make([]byte, sz)
	n, err := o.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	if n != sz {
		return nil, fmt.Errorf("%%s size / offset different %%d : %%d", "Marshal", sz, n)
	}
	return data, nil
}
	`, typeName)
		fmt.Fprintln(out)
		// UnmarshalTo
		fmt.Fprintf(out, "func (o *%s) UnmarshalTo(data []byte) (int, error) {\n", typeName)
		fmt.Fprintln(out, "var (")
		fmt.Fprintln(out, "i, n, l int")
		fmt.Fprintln(out, "err error")
		fmt.Fprintln(out, ")")
		for _, f := range v {
			fType := basicTypes.Get(f.Type)
			fieldName := f.Name
			if f.Parent != "" {
				fieldName = f.Parent + "." + f.Name
			}
			if fType != nil {
				if !fType.IsArray {
					fmt.Fprintf(out, `if o.%s, i, err = o.Unmarshal%s(data[n:]); err != nil {
				return 0, err
			}
			n += i
			`, fieldName, fType.Type)
				} else if fType.IsFixed { // fixed array
					fmt.Fprintf(out, `if l, i, err = o.UnmarshalInt(data[n:]); err != nil {
				return 0, err
			}
			n += i
			o.%s = make(%s, l)
			`, fieldName, fType.Name)
					fmt.Fprintf(out, `for j := range o.%s {
				if o.%s[j], i, err = o.Unmarshal%s(data[n:]); err != nil {
					return 0, err
				}
				n += i
			}
			`, fieldName, fieldName, fType.Type)

				} else {
					fmt.Fprintf(out, `if l, i, err = o.UnmarshalInt(data[n:]); err != nil {
				return 0, err
			}
			n += i
			o.%s = make(%s, l)
			for j := range o.%s {
				if o.%s[j], i, err = o.Unmarshal%s(data[n:]); err != nil {
					return 0, err
				}
				n += i
			}
			`, fieldName, fType.Name, fieldName, fieldName, fType.Type)

				}
			}
		}
		fmt.Fprintln(out, "_ = l")
		fmt.Fprintln(out, "return n, nil")
		fmt.Fprintln(out, "}")

		// UnmarshalBinary
		fmt.Fprintf(out, `// Unmarshal decodes data as conform encoding.BinaryUnmarshaler.
func (o *%s) UnmarshalBinary(data []byte) error {
	_, err := o.UnmarshalTo(data)
	return err
}
`, typeName)
	}

	fmt.Println(out.String())
}
