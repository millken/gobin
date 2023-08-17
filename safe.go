package gobin

import (
	"fmt"
	"math"
	"strconv"
)

var _ Marshaler = Safe{}

var (
	marshalSafeInt    func(v int, bs []byte) (n int, err error)
	unmarshalSafeInt  func(bs []byte) (v int, n int, err error)
	marshalSafeUint   func(v uint, bs []byte) (n int, err error)
	unmarshalSafeUint func(bs []byte) (v uint, n int, err error)
)

func init() {
	switch strconv.IntSize {
	case 32:
		marshalSafeInt = marshalSafeInteger32[int]
		unmarshalSafeInt = unmarshalSafeInteger32[int]
		marshalSafeUint = marshalSafeInteger32[uint]
		unmarshalSafeUint = unmarshalSafeInteger32[uint]
	case 64:
		marshalSafeInt = marshalSafeInteger64[int]
		unmarshalSafeInt = unmarshalSafeInteger64[int]
		marshalSafeUint = marshalSafeInteger64[uint]
		unmarshalSafeUint = unmarshalSafeInteger64[uint]
	default:
		panic("unsupported int size")
	}
}

func marshalSafeInteger8[T Integer8](t T, bs []byte) (n int, err error) {
	if len(bs) < 1 {
		return 0, ErrNotEnoughSpace
	}
	bs[0] = byte(t)
	return 1, nil
}

func marshalSafeInteger16[T Integer16](t T, bs []byte) (n int, err error) {
	if len(bs) < 2 {
		return 0, ErrNotEnoughSpace
	}
	bs[0] = byte(t)
	bs[1] = byte(t >> 8)
	return 2, nil
}

func marshalSafeInteger32[T Integer32](t T, bs []byte) (n int, err error) {
	if len(bs) < 4 {
		return 0, ErrNotEnoughSpace
	}
	bs[0] = byte(t)
	bs[1] = byte(t >> 8)
	bs[2] = byte(t >> 16)
	bs[3] = byte(t >> 24)
	return 4, nil
}

func marshalSafeInteger64[T Integer64](t T, bs []byte) (n int, err error) {
	if len(bs) < 8 {
		return 0, ErrNotEnoughSpace
	}
	bs[0] = byte(t)
	bs[1] = byte(t >> 8)
	bs[2] = byte(t >> 16)
	bs[3] = byte(t >> 24)
	bs[4] = byte(t >> 32)
	bs[5] = byte(t >> 40)
	bs[6] = byte(t >> 48)
	bs[7] = byte(t >> 56)
	return 8, nil
}

func unmarshalSafeInteger8[T Integer8](bs []byte) (t T, n int, err error) {
	if len(bs) < 1 {
		return 0, 0, ErrNotEnoughSpace
	}
	return T(bs[0]), 1, nil
}

func unmarshalSafeInteger16[T Integer16](bs []byte) (t T, n int, err error) {
	if len(bs) < 2 {
		return 0, 0, ErrNotEnoughSpace
	}
	t = T(bs[0])
	t |= T(bs[1]) << 8
	return t, 2, nil
}

func unmarshalSafeInteger32[T Integer32](bs []byte) (t T, n int, err error) {
	if len(bs) < 4 {
		return 0, 0, ErrNotEnoughSpace
	}
	t = T(bs[0])
	t |= T(bs[1]) << 8
	t |= T(bs[2]) << 16
	t |= T(bs[3]) << 24
	return t, 4, nil
}

func unmarshalSafeInteger64[T Integer64](bs []byte) (t T, n int, err error) {
	if len(bs) < 8 {
		return 0, 0, ErrNotEnoughSpace
	}
	t = T(bs[0])
	t |= T(bs[1]) << 8
	t |= T(bs[2]) << 16
	t |= T(bs[3]) << 24
	t |= T(bs[4]) << 32
	t |= T(bs[5]) << 40
	t |= T(bs[6]) << 48
	t |= T(bs[7]) << 56
	return t, 8, nil
}

type Safe struct{}

func (Safe) MarshalBool(v bool, bs []byte) (n int, err error) {
	if len(bs) < 1 {
		return 0, ErrNotEnoughSpace
	}
	if v {
		bs[0] = 1
	} else {
		bs[0] = 0
	}
	return 1, nil
}

func (Safe) UnmarshalBool(bs []byte) (v bool, n int, err error) {
	if len(bs) < 1 {
		return false, 0, ErrNotEnoughSpace
	}
	switch bs[0] {
	case 0:
		v = false
	case 1:
		v = true
	default:
		err = ErrInvalidBool
	}
	return v, 1, err
}

func (Safe) MarshalInt(v int, bs []byte) (n int, err error) {
	return marshalSafeInt(v, bs)
}

func (Safe) UnmarshalInt(bs []byte) (v int, n int, err error) {
	return unmarshalSafeInt(bs)
}

func (Safe) MarshalInt8(v int8, bs []byte) (n int, err error) {
	return marshalSafeInteger8(v, bs)
}

func (Safe) UnmarshalInt8(bs []byte) (v int8, n int, err error) {
	return unmarshalSafeInteger8[int8](bs)
}

func (Safe) MarshalInt16(v int16, bs []byte) (n int, err error) {
	return marshalSafeInteger16(v, bs)
}

func (Safe) UnmarshalInt16(bs []byte) (v int16, n int, err error) {
	return unmarshalUnsafeInteger16[int16](bs)
}

func (Safe) MarshalInt32(v int32, bs []byte) (n int, err error) {
	return marshalSafeInteger32(v, bs)
}

func (Safe) UnmarshalInt32(bs []byte) (v int32, n int, err error) {
	return unmarshalSafeInteger32[int32](bs)
}

func (Safe) MarshalInt64(v int64, bs []byte) (n int, err error) {
	return marshalSafeInteger64(v, bs)
}

func (Safe) UnmarshalInt64(bs []byte) (v int64, n int, err error) {
	return unmarshalSafeInteger64[int64](bs)
}

func (Safe) MarshalUint(v uint, bs []byte) (n int, err error) {
	return marshalSafeUint(v, bs)
}

func (Safe) UnmarshalUint(bs []byte) (v uint, n int, err error) {
	return unmarshalSafeUint(bs)
}
func (Safe) MarshalUint8(v uint8, bs []byte) (n int, err error) {
	return marshalSafeInteger8(v, bs)
}

func (Safe) UnmarshalUint8(bs []byte) (v uint8, n int, err error) {
	return unmarshalSafeInteger8[uint8](bs)
}

func (Safe) MarshalUint16(v uint16, bs []byte) (n int, err error) {
	return marshalSafeInteger16(v, bs)
}

func (Safe) UnmarshalUint16(bs []byte) (v uint16, n int, err error) {
	return unmarshalSafeInteger16[uint16](bs)
}

func (Safe) MarshalUint32(v uint32, bs []byte) (n int, err error) {
	return marshalSafeInteger32(v, bs)
}

func (Safe) UnmarshalUint32(bs []byte) (v uint32, n int, err error) {
	return unmarshalSafeInteger32[uint32](bs)
}

func (Safe) MarshalUint64(v uint64, bs []byte) (n int, err error) {
	return marshalSafeInteger64(v, bs)
}

func (Safe) UnmarshalUint64(bs []byte) (v uint64, n int, err error) {
	return unmarshalSafeInteger64[uint64](bs)
}

func (Safe) MarshalString(v string, bs []byte) (n int, err error) {
	n, err = marshalSafeInt(len(v), bs)
	if err != nil {
		return
	}
	if len(bs[n:]) < len(v) {
		return 0, ErrNotEnoughSpace
	}
	return n + copy(bs[n:], v), nil
}

func (Safe) UnmarshalString(bs []byte) (v string, n int, err error) {
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
	return string(bs[n : n+l]), n + l, nil
}

func (Safe) MarshalByte(v byte, bs []byte) (n int, err error) {
	return marshalSafeInteger8(v, bs)
}

func (Safe) UnmarshalByte(bs []byte) (v byte, n int, err error) {
	return unmarshalSafeInteger8[byte](bs)
}

func (Safe) MarshalBytes(v []byte, bs []byte) (n int, err error) {
	n, err = marshalSafeInt(len(v), bs)
	if err != nil {
		return
	}
	if len(bs[n:]) < len(v) {
		return 0, ErrNotEnoughSpace
	}
	return n + copy(bs[n:], v), nil
}

func (Safe) UnmarshalBytes(bs []byte) (v []byte, n int, err error) {
	l, n, err := unmarshalUnsafeInt(bs)
	if err != nil {
		return
	}
	if l < 0 {
		err = ErrNegativeLength
		return
	} else if l == 0 {
		return nil, n, nil
	}
	if len(bs[n:]) < int(l) {
		err = ErrNotEnoughSpace
		return
	}
	return bs[n : n+l], n + l, nil
}

func (Safe) MarshalFloat32(v float32, bs []byte) (n int, err error) {
	return marshalSafeInteger32(math.Float32bits(v), bs)
}

func (Safe) UnmarshalFloat32(bs []byte) (v float32, n int, err error) {
	uv, n, err := unmarshalUnsafeInteger32[uint32](bs)
	if err != nil {
		return
	}
	return math.Float32frombits(uv), n, nil
}

func (Safe) MarshalFloat64(v float64, bs []byte) (n int, err error) {
	return marshalSafeInteger64(math.Float64bits(v), bs)
}

func (Safe) UnmarshalFloat64(bs []byte) (v float64, n int, err error) {
	uv, n, err := unmarshalUnsafeInteger64[uint64](bs)
	if err != nil {
		return
	}
	return math.Float64frombits(uv), n, nil
}

func (Safe) MarshalBinary() ([]byte, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (Safe) UnmarshalBinary([]byte) error {
	return fmt.Errorf("unimplemented")
}
