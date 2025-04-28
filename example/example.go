package example

import (
	"fmt"

	"github.com/millken/gobin"
)

const (
	a int32   = 1
	b float32 = 1.1
	c string  = "hello"
	d bool    = true
	e int64   = 1
	f float64 = 1
	g int64   = 1
)

//go:generate go run github.com/millken/gobin/cmd/gobingen
type SearchRequest struct {
	gobin.Unsafe
	query           string
	page_number     int32
	result_per_page int32
}

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *SearchRequest) MarshalBinary() (data []byte, err error) {
	sz := len(o.query) + 16
	data = make([]byte, sz)
	var offset, n int
	if n, err = o.MarshalString(o.query, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalInt32(o.page_number, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalInt32(o.result_per_page, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if offset != sz {
		return nil, fmt.Errorf("%s size / offset different %d : %d", "Marshal", sz, offset)
	}
	return data, nil
}

// Unmarshal decodes data as conform encoding.BinaryUnmarshaler.
func (o *SearchRequest) UnmarshalBinary(data []byte) error {
	var (
		i, n int
		err  error
	)
	if o.query, i, err = o.UnmarshalString(data[n:]); err != nil {
		return err
	}
	n += i
	if o.page_number, i, err = o.UnmarshalInt32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.result_per_page, i, err = o.UnmarshalInt32(data[n:]); err != nil {
		return err
	}
	n += i

	return nil
}

type SearchResponse struct {
	gobin.Unsafe
	results string
}

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *SearchResponse) MarshalBinary() (data []byte, err error) {
	sz := len(o.results) + 8
	data = make([]byte, sz)
	var offset, n int
	if n, err = o.MarshalString(o.results, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if offset != sz {
		return nil, fmt.Errorf("%s size / offset different %d : %d", "Marshal", sz, offset)
	}
	return data, nil
}

// Unmarshal decodes data as conform encoding.BinaryUnmarshaler.
func (o *SearchResponse) UnmarshalBinary(data []byte) error {
	var (
		i, n int
		err  error
	)
	if o.results, i, err = o.UnmarshalString(data[n:]); err != nil {
		return err
	}
	n += i

	return nil
}

type A struct {
	gobin.Unsafe
	Name     string
	BirthDay uint64
	Phone    []byte
	Siblings int32
	Spouse   bool
	Money    float64
}

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *A) MarshalBinary() (data []byte, err error) {
	sz := len(o.Name) + len(o.Phone) + 37
	data = make([]byte, sz)
	var offset, n int
	if n, err = o.MarshalString(o.Name, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalUint64(o.BirthDay, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalBytes(o.Phone, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalInt32(o.Siblings, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalBool(o.Spouse, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalFloat64(o.Money, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if offset != sz {
		return nil, fmt.Errorf("%s size / offset different %d : %d", "Marshal", sz, offset)
	}
	return data, nil
}

// Unmarshal decodes data as conform encoding.BinaryUnmarshaler.
func (o *A) UnmarshalBinary(data []byte) error {
	var (
		i, n int
		err  error
	)
	if o.Name, i, err = o.UnmarshalString(data[n:]); err != nil {
		return err
	}
	n += i
	if o.BirthDay, i, err = o.UnmarshalUint64(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Phone, i, err = o.UnmarshalBytes(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Siblings, i, err = o.UnmarshalInt32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Spouse, i, err = o.UnmarshalBool(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Money, i, err = o.UnmarshalFloat64(data[n:]); err != nil {
		return err
	}
	n += i

	return nil
}
