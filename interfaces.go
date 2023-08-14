package gobin

import "encoding"

type Marshaler interface {
	encoding.BinaryMarshaler
	MarshalBool(bool, []byte) (int, error)
	MarshalInt8(int8, []byte) (int, error)
	MarshalInt16(int16, []byte) (int, error)
	MarshalInt32(int32, []byte) (int, error)
	MarshalInt64(int64, []byte) (int, error)
	MarshalUint8(uint8, []byte) (int, error)
	MarshalUint16(uint16, []byte) (int, error)
	MarshalUint32(uint32, []byte) (int, error)
	MarshalUint64(uint64, []byte) (int, error)
	MarshalFloat32(float32, []byte) (int, error)
	MarshalFloat64(float64, []byte) (int, error)
	MarshalString(string, []byte) (int, error)
	MarshalByte(byte, []byte) (int, error)
	MarshalBytes([]byte, []byte) (int, error)
}

type Unmarshaler interface {
	encoding.BinaryUnmarshaler
	UnmarshalBool([]byte) (bool, int, error)
	UnmarshalInt8([]byte) (int8, int, error)
	UnmarshalInt16([]byte) (int16, int, error)
	UnmarshalInt32([]byte) (int32, int, error)
	UnmarshalInt64([]byte) (int64, int, error)
	UnmarshalUint8([]byte) (uint8, int, error)
	UnmarshalUint16([]byte) (uint16, int, error)
	UnmarshalUint32([]byte) (uint32, int, error)
	UnmarshalUint64([]byte) (uint64, int, error)
	UnmarshalFloat32([]byte) (float32, int, error)
	UnmarshalFloat64([]byte) (float64, int, error)
	UnmarshalString([]byte) (string, int, error)
	UnmarshalByte([]byte) (byte, int, error)
	UnmarshalBytes([]byte) ([]byte, int, error)
}

// Integer64 is a constraint that permits any 64-bit integer type.
type Integer64 interface {
	~uint | ~uint64 | ~int | ~int64
}

// Integer64 is a constraint that permits any 32-bit integer type.
type Integer32 interface {
	~uint | ~uint32 | ~int | ~int32
}

// Integer64 is a constraint that permits any 16-bit integer type.
type Integer16 interface {
	~uint16 | ~int16
}

// Integer64 is a constraint that permits any 8-bit integer type.
type Integer8 interface {
	~uint8 | ~int8
}

// Num64 is a constraint that permits any 64-bit integer or float type.
type Num64 interface {
	~float64 | Integer64
}

// Num32 is a constraint that permits any 32-bit integer or float type.
type Num32 interface {
	~float32 | Integer32
}

// MarshallerFn is a functional implementation of the Marshaller interface.
type MarshallerFn[T any] func(t T, bs []byte) (n int, err error)

// UnmarshallerFn is a functional implementation of the Unmarshaller interface.
type UnmarshallerFn[T any] func(bs []byte) (t T, n int, err error)
