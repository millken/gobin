package gobin

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSafe(t *testing.T) {
	r := require.New(t)
	v := Safe{}
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
			n, err := v.MarshalFloat32(test.f, bs)
			r.NoError(err)
			r.Equal(test.n, n)
			f2, n, err := v.UnmarshalFloat32(bs)
			r.NoError(err)
			r.Equal(test.n, n)
			r.Equal(test.f, f2)
		}
	})
	t.Run("float32 should return ErrNotEnoughSpace if there is no space in bs", func(t *testing.T) {
		var f float32 = -1.0
		bs := make([]byte, 2)
		n, err := v.MarshalFloat32(f, bs)
		r.ErrorIs(err, ErrNotEnoughSpace)
		r.Equal(0, n)
		f2, n, err := v.UnmarshalFloat32(bs)
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
			n, err := v.MarshalFloat64(test.f, bs)
			r.NoError(err)
			r.Equal(test.n, n)
			f2, n, err := v.UnmarshalFloat64(bs)
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
			n, err := v.MarshalInt(i, bs)
			r.NoError(err)
			r.Equal(8, n)
			i2, n, err := v.UnmarshalInt(bs)
			r.NoError(err)
			r.Equal(8, n)
			r.Equal(i, i2)
		})

		t.Run("int should return ErrNotEnoughSpace if there is no space in bs", func(t *testing.T) {
			var i int = 1
			//testA[int8](int8(i), MarshalInt8, UnmarshalInt8, 1, r)
			testA[int16](int16(i), v.MarshalInt16, v.UnmarshalInt16, 2, r)
			testA[int32](int32(i), v.MarshalInt32, v.UnmarshalInt32, 4, r)
			testA[int64](int64(i), v.MarshalInt64, v.UnmarshalInt64, 8, r)
		})
		t.Run("int8", func(t *testing.T) {
			var i int8 = 1
			bs := make([]byte, 1)
			n, err := v.MarshalInt8(i, bs)
			r.NoError(err)
			r.Equal(1, n)
			i2, n, err := v.UnmarshalInt8(bs)
			r.NoError(err)
			r.Equal(1, n)
			r.Equal(i, i2)
		})
		t.Run("int16", func(t *testing.T) {
			var i int16 = 1
			bs := make([]byte, 2)
			n, err := v.MarshalInt16(i, bs)
			r.NoError(err)
			r.Equal(2, n)
			i2, n, err := v.UnmarshalInt16(bs)
			r.NoError(err)
			r.Equal(2, n)
			r.Equal(i, i2)
		})
		t.Run("int32", func(t *testing.T) {
			var i int32 = 1
			bs := make([]byte, 4)
			n, err := v.MarshalInt32(i, bs)
			r.NoError(err)
			r.Equal(4, n)
			i2, n, err := v.UnmarshalInt32(bs)
			r.NoError(err)
			r.Equal(4, n)
			r.Equal(i, i2)
		})
		t.Run("int64", func(t *testing.T) {
			var i int64 = 1
			bs := make([]byte, 8)
			n, err := v.MarshalInt64(i, bs)
			r.NoError(err)
			r.Equal(8, n)
			i2, n, err := v.UnmarshalInt64(bs)
			r.NoError(err)
			r.Equal(8, n)
			r.Equal(i, i2)
		})
	})
	t.Run("uint", func(t *testing.T) {
		t.Run("uint", func(t *testing.T) {
			var i uint = 1
			bs := make([]byte, 8)
			n, err := v.MarshalUint(i, bs)
			r.NoError(err)
			r.Equal(8, n)
			i2, n, err := v.UnmarshalUint(bs)
			r.NoError(err)
			r.Equal(8, n)
			r.Equal(i, i2)
		})
		t.Run("uint should return ErrNotEnoughSpace if there is no space in bs", func(t *testing.T) {
			var i int = 1
			//testA[int8](int8(i), MarshalInt8, UnmarshalInt8, 1, r)
			testA[uint16](uint16(i), v.MarshalUint16, v.UnmarshalUint16, 2, r)
			testA[uint32](uint32(i), v.MarshalUint32, v.UnmarshalUint32, 4, r)
			testA[uint64](uint64(i), v.MarshalUint64, v.UnmarshalUint64, 8, r)
		})
		t.Run("uint8", func(t *testing.T) {
			var i uint8 = 1
			bs := make([]byte, 1)
			n, err := v.MarshalUint8(i, bs)
			r.NoError(err)
			r.Equal(1, n)
			i2, n, err := v.UnmarshalUint8(bs)
			r.NoError(err)
			r.Equal(1, n)
			r.Equal(i, i2)
		})
		t.Run("uint16", func(t *testing.T) {
			var i uint16 = 1
			bs := make([]byte, 2)
			n, err := v.MarshalUint16(i, bs)
			r.NoError(err)
			r.Equal(2, n)
			i2, n, err := v.UnmarshalUint16(bs)
			r.NoError(err)
			r.Equal(2, n)
			r.Equal(i, i2)
		})
		t.Run("uint32", func(t *testing.T) {
			var i uint32 = 1
			bs := make([]byte, 4)
			n, err := v.MarshalUint32(i, bs)
			r.NoError(err)
			r.Equal(4, n)
			i2, n, err := v.UnmarshalUint32(bs)
			r.NoError(err)
			r.Equal(4, n)
			r.Equal(i, i2)
		})
		t.Run("uint64", func(t *testing.T) {
			var i uint64 = 1
			bs := make([]byte, 8)
			n, err := v.MarshalUint64(i, bs)
			r.NoError(err)
			r.Equal(8, n)
			i2, n, err := v.UnmarshalUint64(bs)
			r.NoError(err)
			r.Equal(8, n)
			r.Equal(i, i2)
		})
	})
	t.Run("string", func(t *testing.T) {
		s := "hello world"
		bs := make([]byte, 100)
		n, err := v.MarshalString(s, bs)
		r.NoError(err)
		r.Equal(19, n)
		s2, n, err := v.UnmarshalString(bs)
		r.NoError(err)
		r.Equal(19, n)
		r.Equal(s, s2)
	})
	t.Run("bool", func(t *testing.T) {
		b := true
		bs := make([]byte, 1)
		n, err := v.MarshalBool(b, bs)
		r.NoError(err)
		r.Equal(1, n)
		b2, n, err := v.UnmarshalBool(bs)
		r.NoError(err)
		r.Equal(1, n)
		r.Equal(b, b2)
	})
	t.Run("byte", func(t *testing.T) {
		b := byte(1)
		bs := make([]byte, 1)
		n, err := v.MarshalByte(b, bs)
		r.NoError(err)
		r.Equal(1, n)
		b2, n, err := v.UnmarshalByte(bs)
		r.NoError(err)
		r.Equal(1, n)
		r.Equal(b, b2)
	})

	t.Run("bytes", func(t *testing.T) {
		bs := []byte("hello world")
		bs2 := make([]byte, 100)
		n, err := v.MarshalBytes(bs, bs2)
		r.NoError(err)
		r.Equal(19, n)
		bs3, n, err := v.UnmarshalBytes(bs2)
		r.NoError(err)
		r.Equal(19, n)
		r.Equal(bs, bs3)
	})
}
