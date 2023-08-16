package example

import (
	"fmt"

	"github.com/millken/gobin"
)

// PackageType is the type of package
type PackageType uint16

const (
	// DATA is a data package
	PackageType_DATA   PackageType = 0
	PackageType_CONFIG PackageType = 1
	PackageType_STATE  PackageType = 2
)

func (o *PackageType) Size() int {
	return 2
}

// MarshalTo writes a wire-format message to w.
func (o *PackageType) MarshalTo(w []byte) (int, error) {
	w[0] = byte(*o)
	w[1] = byte(*o >> 8)
	return 2, nil
}

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *PackageType) MarshalBinary() ([]byte, error) {
	data := make([]byte, 2)
	_, err := o.MarshalTo(data)
	return data, err
}

// Unmarshal decodes data as conform encoding.BinaryUnmarshaler.
func (o *PackageType) UnmarshalBinary(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("invalid data size %d", len(data))
	}
	*o = PackageType(uint16(data[0]) | uint16(data[1])<<8)
	return nil
}

type BinPackage struct {
	gobin.Safe
	Type      *PackageType
	Data      []byte
	Timestamp uint32
	Signature []byte
}

// MarshalBinary encodes o as conform encoding.BinaryMarshaler.
func (o *BinPackage) MarshalBinary() (data []byte, err error) {
	sz := o.Type.Size() + len(o.Data) + len(o.Signature) + 20
	data = make([]byte, sz)
	var offset, n int
	if n, err = o.Type.MarshalTo(data); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalBytes(o.Data, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalUint32(o.Timestamp, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if n, err = o.MarshalBytes(o.Signature, data[offset:]); err != nil {
		return nil, err
	}
	offset += n
	if offset != sz {
		return nil, fmt.Errorf("%s size / offset different %d : %d", "Marshal", sz, offset)
	}
	return data, nil
}

// Unmarshal decodes data as conform encoding.BinaryUnmarshaler.
func (o *BinPackage) UnmarshalBinary(data []byte) error {
	var (
		i, n int
		err  error
	)
	if o.Data, i, err = o.UnmarshalBytes(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Timestamp, i, err = o.UnmarshalUint32(data[n:]); err != nil {
		return err
	}
	n += i
	if o.Signature, i, err = o.UnmarshalBytes(data[n:]); err != nil {
		return err
	}
	n += i

	return nil
}
