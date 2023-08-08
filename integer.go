package gobin

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
type MarshallerFn[T any] func(t T, bs []byte) (n int)

// UnmarshallerFn is a functional implementation of the Unmarshaller interface.
type UnmarshallerFn[T any] func(bs []byte) (t T, n int, err error)
