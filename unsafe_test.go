package gobin

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func testA[T any](i T, m MarshallerFn[T], u UnmarshallerFn[T], nn int, r *require.Assertions) {
	bs := make([]byte, 1)
	n := m(i, bs)
	r.Equal(nn, n)
	_, n, err := u(bs)
	r.ErrorIs(err, ErrNotEnoughSpace)
	r.Equal(0, n)
}
func TestUnsafe(t *testing.T) {
	r := require.New(t)
	t.Run("float32", func(t *testing.T) {
		var tests = []struct {
			f float32
			n int
		}{
			{0.0, 4},
			{1.0, 4},
			{math.MaxFloat32, 4},
			{math.SmallestNonzeroFloat32, 4},
		}
		for _, test := range tests {
			bs := make([]byte, 4)
			n := MarshalFloat32(test.f, bs)
			r.Equal(test.n, n)
			f2, n, err := UnmarshalFloat32(bs)
			r.NoError(err)
			r.Equal(test.n, n)
			r.Equal(test.f, f2)
		}
	})
	t.Run("float32 should return ErrNotEnoughSpace if there is no space in bs", func(t *testing.T) {
		var f float32 = -1.0
		bs := make([]byte, 2)
		n := MarshalFloat32(f, bs)
		r.Equal(4, n)
		f2, n, err := UnmarshalFloat32(bs)
		r.ErrorIs(err, ErrNotEnoughSpace)
		r.Equal(0, n)
		r.Equal(float32(0), f2)
	})
	t.Run("float64", func(t *testing.T) {
		var tests = []struct {
			f float64
			n int
		}{
			{0.0, 8},
			{1.0, 8},
			{math.MaxFloat64, 8},
			{math.SmallestNonzeroFloat64, 8},
			{math.Inf(1), 8},
			{math.Inf(-1), 8},
			{math.NaN(), 8},
		}
		for _, test := range tests {
			bs := make([]byte, 8)
			n := MarshalFloat64(test.f, bs)
			r.Equal(test.n, n)
			f2, n, err := UnmarshalFloat64(bs)
			r.NoError(err)
			r.Equal(test.n, n)
			if math.IsNaN(test.f) {
				r.True(math.IsNaN(f2))
			} else {
				r.Equal(test.f, f2)
			}
		}
	})
	t.Run("int", func(t *testing.T) {
		t.Run("int", func(t *testing.T) {
			var i int = 1
			bs := make([]byte, 8)
			n := MarshalInt(i, bs)
			r.Equal(8, n)
			i2, n, err := UnmarshalInt(bs)
			r.NoError(err)
			r.Equal(8, n)
			r.Equal(i, i2)
		})

		t.Run("int should return ErrNotEnoughSpace if there is no space in bs", func(t *testing.T) {
			var i int = 1
			//testA[int8](int8(i), MarshalInt8, UnmarshalInt8, 1, r)
			testA[int16](int16(i), MarshalInt16, UnmarshalInt16, 2, r)
			testA[int32](int32(i), MarshalInt32, UnmarshalInt32, 4, r)
			testA[int64](int64(i), MarshalInt64, UnmarshalInt64, 8, r)
		})
		t.Run("int8", func(t *testing.T) {
			var i int8 = 1
			bs := make([]byte, 1)
			n := MarshalInt8(i, bs)
			r.Equal(1, n)
			i2, n, err := UnmarshalInt8(bs)
			r.NoError(err)
			r.Equal(1, n)
			r.Equal(i, i2)
		})
		t.Run("int16", func(t *testing.T) {
			var i int16 = 1
			bs := make([]byte, 2)
			n := MarshalInt16(i, bs)
			r.Equal(2, n)
			i2, n, err := UnmarshalInt16(bs)
			r.NoError(err)
			r.Equal(2, n)
			r.Equal(i, i2)
		})
		t.Run("int32", func(t *testing.T) {
			var i int32 = 1
			bs := make([]byte, 4)
			n := MarshalInt32(i, bs)
			r.Equal(4, n)
			i2, n, err := UnmarshalInt32(bs)
			r.NoError(err)
			r.Equal(4, n)
			r.Equal(i, i2)
		})
		t.Run("int64", func(t *testing.T) {
			var i int64 = 1
			bs := make([]byte, 8)
			n := MarshalInt64(i, bs)
			r.Equal(8, n)
			i2, n, err := UnmarshalInt64(bs)
			r.NoError(err)
			r.Equal(8, n)
			r.Equal(i, i2)
		})
	})
	t.Run("uint", func(t *testing.T) {
		t.Run("uint", func(t *testing.T) {
			var i uint = 1
			bs := make([]byte, 8)
			n := MarshalUint(i, bs)
			r.Equal(8, n)
			i2, n, err := UnmarshalUint(bs)
			r.NoError(err)
			r.Equal(8, n)
			r.Equal(i, i2)
		})
		t.Run("uint should return ErrNotEnoughSpace if there is no space in bs", func(t *testing.T) {
			var i int = 1
			//testA[int8](int8(i), MarshalInt8, UnmarshalInt8, 1, r)
			testA[uint16](uint16(i), MarshalUint16, UnmarshalUint16, 2, r)
			testA[uint32](uint32(i), MarshalUint32, UnmarshalUint32, 4, r)
			testA[uint64](uint64(i), MarshalUint64, UnmarshalUint64, 8, r)
		})
		t.Run("uint8", func(t *testing.T) {
			var i uint8 = 1
			bs := make([]byte, 1)
			n := MarshalUint8(i, bs)
			r.Equal(1, n)
			i2, n, err := UnmarshalUint8(bs)
			r.NoError(err)
			r.Equal(1, n)
			r.Equal(i, i2)
		})
		t.Run("uint16", func(t *testing.T) {
			var i uint16 = 1
			bs := make([]byte, 2)
			n := MarshalUint16(i, bs)
			r.Equal(2, n)
			i2, n, err := UnmarshalUint16(bs)
			r.NoError(err)
			r.Equal(2, n)
			r.Equal(i, i2)
		})
		t.Run("uint32", func(t *testing.T) {
			var i uint32 = 1
			bs := make([]byte, 4)
			n := MarshalUint32(i, bs)
			r.Equal(4, n)
			i2, n, err := UnmarshalUint32(bs)
			r.NoError(err)
			r.Equal(4, n)
			r.Equal(i, i2)
		})
		t.Run("uint64", func(t *testing.T) {
			var i uint64 = 1
			bs := make([]byte, 8)
			n := MarshalUint64(i, bs)
			r.Equal(8, n)
			i2, n, err := UnmarshalUint64(bs)
			r.NoError(err)
			r.Equal(8, n)
			r.Equal(i, i2)
		})
	})
	t.Run("string", func(t *testing.T) {
		s := "hello world"
		bs := make([]byte, 100)
		n := MarshalString(s, bs)
		r.Equal(19, n)
		s2, n, err := UnmarshalString(bs)
		r.NoError(err)
		r.Equal(19, n)
		r.Equal(s, s2)
	})
	t.Run("bool", func(t *testing.T) {
		b := true
		bs := make([]byte, 1)
		n := MarshalBool(b, bs)
		r.Equal(1, n)
		b2, n, err := UnmarshalBool(bs)
		r.NoError(err)
		r.Equal(1, n)
		r.Equal(b, b2)
	})
	t.Run("byte", func(t *testing.T) {
		b := byte(1)
		bs := make([]byte, 1)
		n := MarshalByte(b, bs)
		r.Equal(1, n)
		b2, n, err := UnmarshalByte(bs)
		r.NoError(err)
		r.Equal(1, n)
		r.Equal(b, b2)
	})

	t.Run("bytes", func(t *testing.T) {
		bs := []byte("hello world")
		bs2 := make([]byte, 100)
		n := MarshalBytes(bs, bs2)
		r.Equal(19, n)
		bs3, n, err := UnmarshalBytes(bs2)
		r.NoError(err)
		r.Equal(19, n)
		r.Equal(bs, bs3)
	})
}
