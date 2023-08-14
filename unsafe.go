package gobin

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"unsafe"
)

var (
	marshalUnsafeInt    func(v int, bs []byte) (int, error)
	unmarshalUnsafeInt  func(bs []byte) (v int, n int, err error)
	marshalUnsafeUint   func(v uint, bs []byte) (int, error)
	unmarshalUnsafeUint func(bs []byte) (v uint, n int, err error)
)

func init() {
	switch strconv.IntSize {
	case 32:
		marshalUnsafeInt = marshalUnsafeInteger32[int]
		unmarshalUnsafeInt = unmarshalUnsafeInteger32[int]
		marshalUnsafeUint = marshalUnsafeInteger32[uint]
		unmarshalUnsafeUint = unmarshalUnsafeInteger32[uint]
	case 64:
		marshalUnsafeInt = marshalUnsafeInteger64[int]
		unmarshalUnsafeInt = unmarshalUnsafeInteger64[int]
		marshalUnsafeUint = marshalUnsafeInteger64[uint]
		unmarshalUnsafeUint = unmarshalUnsafeInteger64[uint]
	default:
		panic("unsupported int size")
	}
}

var (
	ErrNotEnoughSpace = errors.New("not enough space")
	ErrInvalidBool    = errors.New("invalid bool value")
	ErrNegativeLength = errors.New("negative length")
)

func marshalUnsafeInteger8[T Integer8](t T, bs []byte) (int, error) {
	if len(bs) < 1 {
		return 0, ErrNotEnoughSpace
	}
	*(*T)(unsafe.Pointer(&bs[0])) = t
	return 1, nil
}

func marshalUnsafeInteger16[T Integer16](t T, bs []byte) (int, error) {
	if len(bs) < 2 {
		return 0, ErrNotEnoughSpace
	}
	*(*T)(unsafe.Pointer(&bs[0])) = t
	return 2, nil
}

func marshalUnsafeInteger32[T Integer32](t T, bs []byte) (int, error) {
	if len(bs) < 4 {
		return 0, ErrNotEnoughSpace
	}
	*(*T)(unsafe.Pointer(&bs[0])) = t
	return 4, nil
}

func marshalUnsafeInteger64[T Integer64](t T, bs []byte) (int, error) {
	if len(bs) < 8 {
		return 0, ErrNotEnoughSpace
	}
	*(*T)(unsafe.Pointer(&bs[0])) = t
	return 8, nil
}

func unmarshalUnsafeInteger8[T Integer8](bs []byte) (t T, n int, err error) {
	if len(bs) < 1 {
		return 0, 0, ErrNotEnoughSpace
	}
	return *(*T)(unsafe.Pointer(&bs[0])), 1, nil
}

func unmarshalUnsafeInteger16[T Integer16](bs []byte) (t T, n int, err error) {
	if len(bs) < 2 {
		return 0, 0, ErrNotEnoughSpace
	}
	return *(*T)(unsafe.Pointer(&bs[0])), 2, nil
}

func unmarshalUnsafeInteger32[T Integer32](bs []byte) (t T, n int, err error) {
	if len(bs) < 4 {
		return 0, 0, ErrNotEnoughSpace
	}
	return *(*T)(unsafe.Pointer(&bs[0])), 4, nil
}

func unmarshalUnsafeInteger64[T Integer64](bs []byte) (t T, n int, err error) {
	if len(bs) < 8 {
		return 0, 0, ErrNotEnoughSpace
	}
	return *(*T)(unsafe.Pointer(&bs[0])), 8, nil
}

var _ Marshaler = Unsafe{}

type Unsafe struct{}

func (Unsafe) MarshalBool(v bool, bs []byte) (int, error) {
	if len(bs) < 1 {
		return 0, ErrNotEnoughSpace
	}
	*(*bool)(unsafe.Pointer(&bs[0])) = v
	return 1, nil
}

func (Unsafe) UnmarshalBool(bs []byte) (v bool, n int, err error) {
	if len(bs) < 1 {
		return false, 0, ErrNotEnoughSpace
	}
	if bs[0] > 1 {
		err = ErrInvalidBool
		return
	}
	return *(*bool)(unsafe.Pointer(&bs[0])), 1, nil
}

func (Unsafe) MarshalInt(v int, bs []byte) (int, error) {
	return marshalUnsafeInt(v, bs)
}

func (Unsafe) UnmarshalInt(bs []byte) (v int, n int, err error) {
	return unmarshalUnsafeInt(bs)
}

func (Unsafe) MarshalInt8(v int8, bs []byte) (int, error) {
	return marshalUnsafeInteger8(v, bs)
}

func (Unsafe) UnmarshalInt8(bs []byte) (v int8, n int, err error) {
	return unmarshalUnsafeInteger8[int8](bs)
}

func (Unsafe) MarshalInt16(v int16, bs []byte) (int, error) {
	return marshalUnsafeInteger16(v, bs)
}

func (Unsafe) UnmarshalInt16(bs []byte) (v int16, n int, err error) {
	return unmarshalUnsafeInteger16[int16](bs)
}

func (Unsafe) MarshalInt32(v int32, bs []byte) (int, error) {
	return marshalUnsafeInteger32(v, bs)
}

func (Unsafe) UnmarshalInt32(bs []byte) (v int32, n int, err error) {
	return unmarshalUnsafeInteger32[int32](bs)
}

func (Unsafe) MarshalInt64(v int64, bs []byte) (int, error) {
	return marshalUnsafeInteger64(v, bs)
}

func (Unsafe) UnmarshalInt64(bs []byte) (v int64, n int, err error) {
	return unmarshalUnsafeInteger64[int64](bs)
}

func (Unsafe) MarshalUint(v uint, bs []byte) (int, error) {
	return marshalUnsafeUint(v, bs)
}

func (Unsafe) UnmarshalUint(bs []byte) (v uint, n int, err error) {
	return unmarshalUnsafeUint(bs)
}
func (Unsafe) MarshalUint8(v uint8, bs []byte) (int, error) {
	return marshalUnsafeInteger8(v, bs)
}

func (Unsafe) UnmarshalUint8(bs []byte) (v uint8, n int, err error) {
	return unmarshalUnsafeInteger8[uint8](bs)
}

func (Unsafe) MarshalUint16(v uint16, bs []byte) (int, error) {
	return marshalUnsafeInteger16(v, bs)
}

func (Unsafe) UnmarshalUint16(bs []byte) (v uint16, n int, err error) {
	return unmarshalUnsafeInteger16[uint16](bs)
}

func (Unsafe) MarshalUint32(v uint32, bs []byte) (int, error) {
	return marshalUnsafeInteger32(v, bs)
}

func (Unsafe) UnmarshalUint32(bs []byte) (v uint32, n int, err error) {
	return unmarshalUnsafeInteger32[uint32](bs)
}

func (Unsafe) MarshalUint64(v uint64, bs []byte) (int, error) {
	return marshalUnsafeInteger64(v, bs)
}

func (Unsafe) UnmarshalUint64(bs []byte) (v uint64, n int, err error) {
	return unmarshalUnsafeInteger64[uint64](bs)
}

func (Unsafe) MarshalString(v string, bs []byte) (n int, err error) {
	n, err = marshalUnsafeInt(len(v), bs)
	if len(bs[n:]) < len(v) {
		return 0, ErrNotEnoughSpace
	}
	return n + copy(bs[n:], unsafe.Slice(unsafe.StringData(v), len(v))), nil
}

func (Unsafe) UnmarshalString(bs []byte) (v string, n int, err error) {
	l, n, err := unmarshalUnsafeInt(bs)
	if err != nil {
		return
	}
	if l < 0 {
		err = ErrNegativeLength
		return
	}
	if len(bs[n:]) < int(l) {
		err = ErrNotEnoughSpace
		return
	}
	c := bs[n : n+l]
	return unsafe.String(&c[0], len(c)), n + l, nil
}

func (Unsafe) MarshalByte(v byte, bs []byte) (int, error) {
	return marshalUnsafeInteger8(v, bs)
}

func (Unsafe) UnmarshalByte(bs []byte) (v byte, n int, err error) {
	return unmarshalUnsafeInteger8[byte](bs)
}

func (Unsafe) MarshalBytes(v []byte, bs []byte) (n int, err error) {
	n, err = marshalUnsafeInt(len(v), bs)
	if err != nil {
		return
	}
	if len(bs[n:]) < len(v) {
		return 0, ErrNotEnoughSpace
	}
	return n + copy(bs[n:], v), nil
}

func (Unsafe) UnmarshalBytes(bs []byte) (v []byte, n int, err error) {
	l, n, err := unmarshalUnsafeInt(bs)
	if err != nil {
		return
	}
	if l < 0 {
		err = ErrNegativeLength
		return
	}
	if len(bs[n:]) < int(l) {
		err = ErrNotEnoughSpace
		return
	}
	return bs[n : n+l], n + l, nil
}

func (Unsafe) MarshalFloat32(v float32, bs []byte) (int, error) {
	return marshalUnsafeInteger32(math.Float32bits(v), bs)
}

func (Unsafe) UnmarshalFloat32(bs []byte) (v float32, n int, err error) {
	uv, n, err := unmarshalUnsafeInteger32[uint32](bs)
	if err != nil {
		return
	}
	return math.Float32frombits(uv), n, nil
}

func (Unsafe) MarshalFloat64(v float64, bs []byte) (int, error) {
	return marshalUnsafeInteger64(math.Float64bits(v), bs)
}

func (Unsafe) UnmarshalFloat64(bs []byte) (v float64, n int, err error) {
	uv, n, err := unmarshalUnsafeInteger64[uint64](bs)
	if err != nil {
		return
	}
	return math.Float64frombits(uv), n, nil
}

func (Unsafe) MarshalBinary() ([]byte, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (Unsafe) UnmarshalBinary([]byte) error {
	return fmt.Errorf("unimplemented")
}
