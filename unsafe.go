package gobin

import (
	"errors"
	"math"
	"unsafe"
)

var (
	marshalInt    func(v int, bs []byte) (n int)
	unmarshalInt  func(bs []byte) (v int, n int, err error)
	marshalUint   func(v uint, bs []byte) (n int)
	unmarshalUint func(bs []byte) (v uint, n int, err error)
)

func init() {
	switch unsafe.Sizeof(int(0)) {
	case 4:
		marshalInt = marshalInteger32[int]
		unmarshalInt = unmarshalInteger32[int]
		marshalUint = marshalInteger32[uint]
		unmarshalUint = unmarshalInteger32[uint]
	case 8:
		marshalInt = marshalInteger64[int]
		unmarshalInt = unmarshalInteger64[int]
		marshalUint = marshalInteger64[uint]
		unmarshalUint = unmarshalInteger64[uint]
	default:
		panic("unsupported int size")
	}
}

var (
	ErrNotEnoughSpace = errors.New("not enough space")
	ErrInvalidBool    = errors.New("invalid bool value")
	ErrNegativeLength = errors.New("negative length")
)

func marshalInteger8[T Integer8](t T, bs []byte) (n int) {
	*(*T)(unsafe.Pointer(&bs[0])) = t
	return 1
}

func marshalInteger16[T Integer16](t T, bs []byte) (n int) {
	*(*T)(unsafe.Pointer(&bs[0])) = t
	return 2
}

func marshalInteger32[T Integer32](t T, bs []byte) (n int) {
	*(*T)(unsafe.Pointer(&bs[0])) = t
	return 4
}

func marshalInteger64[T Integer64](t T, bs []byte) (n int) {
	*(*T)(unsafe.Pointer(&bs[0])) = t
	return 8
}

func unmarshalInteger8[T Integer8](bs []byte) (t T, n int, err error) {
	if len(bs) < 1 {
		return 0, 0, ErrNotEnoughSpace
	}
	return *(*T)(unsafe.Pointer(&bs[0])), 1, nil
}

func unmarshalInteger16[T Integer16](bs []byte) (t T, n int, err error) {
	if len(bs) < 2 {
		return 0, 0, ErrNotEnoughSpace
	}
	return *(*T)(unsafe.Pointer(&bs[0])), 2, nil
}

func unmarshalInteger32[T Integer32](bs []byte) (t T, n int, err error) {
	if len(bs) < 4 {
		return 0, 0, ErrNotEnoughSpace
	}
	return *(*T)(unsafe.Pointer(&bs[0])), 4, nil
}

func unmarshalInteger64[T Integer64](bs []byte) (t T, n int, err error) {
	if len(bs) < 8 {
		return 0, 0, ErrNotEnoughSpace
	}
	return *(*T)(unsafe.Pointer(&bs[0])), 8, nil
}

func MarshalBool(v bool, bs []byte) (n int) {
	*(*bool)(unsafe.Pointer(&bs[0])) = v
	return 1
}

func UnmarshalBool(bs []byte) (v bool, n int, err error) {
	if len(bs) < 1 {
		return false, 0, ErrNotEnoughSpace
	}
	if bs[0] > 1 {
		err = ErrInvalidBool
		return
	}
	return *(*bool)(unsafe.Pointer(&bs[0])), 1, nil
}

// Returns the number of used bytes. It will panic if receives too small bs.
func MarshalByte(v byte, bs []byte) (n int) {
	return marshalInteger8(v, bs)
}

func UnmarshalByte(bs []byte) (v byte, n int, err error) {
	return unmarshalInteger8[byte](bs)
}

func MarshalString(v string, bs []byte) (n int) {
	n = marshalInt(len(v), bs)
	if len(bs[n:]) < len(v) {
		panic(ErrNotEnoughSpace)
	}
	return n + copy(bs[n:], unsafe.Slice(unsafe.StringData(v), len(v)))
}

func UnmarshalString(bs []byte) (v string, n int, err error) {
	l, n, err := unmarshalInt(bs)
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

func MarshalBytes(v []byte, bs []byte) (n int) {
	n = marshalInt(len(v), bs)
	if len(bs[n:]) < len(v) {
		panic(ErrNotEnoughSpace)
	}
	return n + copy(bs[n:], v)
}

func UnmarshalBytes(bs []byte) (v []byte, n int, err error) {
	l, n, err := unmarshalInt(bs)
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

func MarshalFloat32(v float32, bs []byte) (n int) {
	return marshalInteger32(math.Float32bits(v), bs)
}

func UnmarshalFloat32(bs []byte) (v float32, n int, err error) {
	uv, n, err := unmarshalInteger32[uint32](bs)
	if err != nil {
		return
	}
	return math.Float32frombits(uv), n, nil
}

func MarshalFloat64(v float64, bs []byte) (n int) {
	return marshalInteger64(math.Float64bits(v), bs)
}

func UnmarshalFloat64(bs []byte) (v float64, n int, err error) {
	uv, n, err := unmarshalInteger64[uint64](bs)
	if err != nil {
		return
	}
	return math.Float64frombits(uv), n, nil
}

func MarshalUint(v uint, bs []byte) (n int) {
	return marshalUint(v, bs)
}

func UnmarshalUint(bs []byte) (v uint, n int, err error) {
	return unmarshalUint(bs)
}
func MarshalUint8(v uint8, bs []byte) (n int) {
	return marshalInteger8(v, bs)
}

func UnmarshalUint8(bs []byte) (v uint8, n int, err error) {
	return unmarshalInteger8[uint8](bs)
}

func MarshalUint16(v uint16, bs []byte) (n int) {
	return marshalInteger16(v, bs)
}

func UnmarshalUint16(bs []byte) (v uint16, n int, err error) {
	return unmarshalInteger16[uint16](bs)
}

func MarshalUint32(v uint32, bs []byte) (n int) {
	return marshalInteger32(v, bs)
}

func UnmarshalUint32(bs []byte) (v uint32, n int, err error) {
	return unmarshalInteger32[uint32](bs)
}

func MarshalUint64(v uint64, bs []byte) (n int) {
	return marshalInteger64(v, bs)
}

func UnmarshalUint64(bs []byte) (v uint64, n int, err error) {
	return unmarshalInteger64[uint64](bs)
}

func MarshalInt(v int, bs []byte) (n int) {
	return marshalInt(v, bs)
}

func UnmarshalInt(bs []byte) (v int, n int, err error) {
	return unmarshalInt(bs)
}

func MarshalInt8(v int8, bs []byte) (n int) {
	return marshalInteger8(v, bs)
}

func UnmarshalInt8(bs []byte) (v int8, n int, err error) {
	return unmarshalInteger8[int8](bs)
}

func MarshalInt16(v int16, bs []byte) (n int) {
	return marshalInteger16(v, bs)
}

func UnmarshalInt16(bs []byte) (v int16, n int, err error) {
	return unmarshalInteger16[int16](bs)
}

func MarshalInt32(v int32, bs []byte) (n int) {
	return marshalInteger32(v, bs)
}

func UnmarshalInt32(bs []byte) (v int32, n int, err error) {
	return unmarshalInteger32[int32](bs)
}

func MarshalInt64(v int64, bs []byte) (n int) {
	return marshalInteger64(v, bs)
}

func UnmarshalInt64(bs []byte) (v int64, n int, err error) {
	return unmarshalInteger64[int64](bs)
}
