package demo

import (
	"fmt"

	"github.com/millken/gobin"
)

type HomeBlock struct {
	ID string `json:"id"`
}

//gobin:binary
type GetTab struct {
	gobin.Unsafe
	Code  string `json:"code"`
	Datas struct {
		HomeBlock HomeBlock `json:"home_block"`
		Block     []struct {
			ID         string `json:"id"`
			TargetType string `json:"target_type"`
			Videos     []struct {
				Vid string `json:"vid"`
			} `json:"videos"`
		} `json:"block"`
	} `json:"datas"`
}

//gobin:binary
type A struct {
	gobin.Unsafe
	B0 string
	B1 int
	B2 struct {
		B20 string
		B21 int
	}
	B3 []int
	B4 struct {
		sid  int `gobin:"id"`
		b2   []byte
		Name string `gobin:"name"`
		Sub3 struct {
			Id int
			B3 []byte
		}
	}
}

func (o *A) SizeBinary() int {
	size := 0

	// B0
	size += 8
	size += len(o.B0)
	// B1
	size += 8

	// B2
	size += 8
	size += len(o.B2.B20)
	size += 8

	// B3
	size += 8
	for _, v := range o.B3 {
		_ = v

		size += 8

	}

	// B4
	size += 8

	size += 9
	size += len(o.B4.b2)
	size += 8
	size += len(o.B4.Name)
	size += 8

	size += 9
	size += len(o.B4.Sub3.B3)

	return size
}

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *A) MarshalBinary() (data []byte, err error) {
	sz := o.SizeBinary()
	data = make([]byte, sz)
	if n, err := o.MarshalTo(data); err != nil {
		return nil, err
	} else if n != sz {
		return nil, fmt.Errorf("%s size / offset different %d : %d", "Marshal", sz, n)
	}
	return data, nil
}

// MarshalTo encodes o as conform encoding.BinaryMarshaler.
func (o *A) MarshalTo(data []byte) (int, error) {
	var (
		offset, n int
		err       error
	)

	// B0
	if n, err = o.MarshalString(o.B0, data[offset:]); err != nil {
		return 0, err
	}
	offset += n
	// B1
	if n, err = o.MarshalInt(o.B1, data[offset:]); err != nil {
		return 0, err
	}
	offset += n
	// B2
	if n, err = o.MarshalString(o.B2.B20, data[offset:]); err != nil {
		return 0, err
	}
	offset += n
	if n, err = o.MarshalInt(o.B2.B21, data[offset:]); err != nil {
		return 0, err
	}
	offset += n

	// B3
	if n, err = o.MarshalInt(len(o.B3), data[offset:]); err != nil { // length
		return 0, err
	}
	offset += n
	for _, v := range o.B3 {
		if n, err = o.MarshalInt(v, data[offset:]); err != nil {
			return 0, err
		}
		offset += n
	}

	// B4
	if n, err = o.MarshalInt(o.B4.sid, data[offset:]); err != nil {
		return 0, err
	}
	offset += n
	if n, err = o.MarshalBytes(o.B4.b2, data[offset:]); err != nil {
		return 0, err
	}
	offset += n
	if n, err = o.MarshalString(o.B4.Name, data[offset:]); err != nil {
		return 0, err
	}
	offset += n
	if n, err = o.MarshalInt(o.B4.Sub3.Id, data[offset:]); err != nil {
		return 0, err
	}
	offset += n
	if n, err = o.MarshalBytes(o.B4.Sub3.B3, data[offset:]); err != nil {
		return 0, err
	}
	offset += n

	return offset, nil
}

// UnmarshalBinary decodes o as conform encoding.BinaryUnmarshaler.
func (o *A) UnmarshalBinary(data []byte) error {
	_, err := o.UnmarshalFrom(data)
	return err
}

// UnmarshalFrom decodes o as conform encoding.BinaryUnmarshaler.
func (o *A) UnmarshalFrom(data []byte) (int, error) {
	var (
		i, n, l int
		err     error
	)

	// B0
	if o.B0, i, err = o.UnmarshalString(data[n:]); err != nil {
		return 0, err
	}
	n += i

	// B1
	if o.B1, i, err = o.UnmarshalInt(data[n:]); err != nil {
		return 0, err
	}
	n += i

	// B2
	if o.B2.B20, i, err = o.UnmarshalString(data[n:]); err != nil {
		return 0, err
	}
	n += i
	if o.B2.B21, i, err = o.UnmarshalInt(data[n:]); err != nil {
		return 0, err
	}
	n += i

	// B3
	if l, i, err = o.UnmarshalInt(data[n:]); err != nil { // length
		return 0, err
	}
	n += i

	o.B3 = make([]int, l)
	for i0 := range o.B3 {
		if o.B3[i0], i, err = o.UnmarshalInt(data[n:]); err != nil {
			return 0, err
		}
		n += i
	}

	// B4
	if o.B4.sid, i, err = o.UnmarshalInt(data[n:]); err != nil {
		return 0, err
	}
	n += i
	if o.B4.b2, i, err = o.UnmarshalBytes(data[n:]); err != nil {
		return 0, err
	}
	n += i
	if o.B4.Name, i, err = o.UnmarshalString(data[n:]); err != nil {
		return 0, err
	}
	n += i
	if o.B4.Sub3.Id, i, err = o.UnmarshalInt(data[n:]); err != nil {
		return 0, err
	}
	n += i
	if o.B4.Sub3.B3, i, err = o.UnmarshalBytes(data[n:]); err != nil {
		return 0, err
	}
	n += i

	_ = l
	return n, nil
}
