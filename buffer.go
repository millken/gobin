package gobin

import (
	"io"
	"sync"
	"unsafe"
)

// buffer.go manages Buffer and its responsibilities.
// Its existence is justified by the performance increase gained
// over more general implementations and the fact that we allow
// these buffers to be pooled, which reduces our allocation
// profile quite significantly.

// Buffer is used to pass on to the encoders Marshal methods.
type Buffer struct {
	Bytes []byte
}

var _ io.Writer = &Buffer{} // commit to compatibility with io.Writer

// Write a chunk of bytes to the buffer
func (b *Buffer) Write(v []byte) (int, error) {
	b.Bytes = append(b.Bytes, v...)
	return len(v), nil
}

// WriteString writes a string to the buffer
func (b *Buffer) WriteString(v string) {
	b.Bytes = append(b.Bytes, v...)
}

// WriteByte writes a single byte into the output buffer
func (b *Buffer) WriteByte(v byte) error {
	b.Bytes = append(b.Bytes, v)
	return nil
}

// Reset allows this to be reused by emptying
func (b *Buffer) Reset() {
	b.Bytes = b.Bytes[:0]
}

func (b *Buffer) String() string {
	return *(*string)(unsafe.Pointer(&b.Bytes))
}

// WriteTo writes the contents of our buffer to an io.Writer
func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.Bytes)
	return int64(n), err
}

var bufpool = sync.Pool{
	New: func() interface{} { return &Buffer{} },
}

// NewBufferFromPool returns a pointer to a zerod Buffer. This may be retrieved from a
// pool. When you're done with it, call 'ReturnToPool'.
func NewBufferFromPool() *Buffer {
	b := bufpool.Get().(*Buffer)
	b.Reset()
	return b
}

// NewBufferFromPoolWithCap returns a pointer to a zero'd Buffer with its underlying
// capacity set. This may be retrieved from a pool. When you're done with it, call 'ReturnToPool'.
func NewBufferFromPoolWithCap(size int) *Buffer {
	b := bufpool.Get().(*Buffer)

	if c := cap(b.Bytes); c < size {
		b.Bytes = make([]byte, 0, size)
	} else if c > 0 {
		b.Reset()
	}

	return b
}

// ReturnToPool puts this instance back in the underlying pool. Reading from or using this instance
// in any way after calling this is invalid.
func (b *Buffer) ReturnToPool() {
	bufpool.Put(b)
}

// const _size = 1024 // by default, create 1 KiB buffers

// // Buffer is a thin wrapper around a byte slice. It's intended to be pooled, so
// // the only way to construct one is via a Pool.
// type Buffer struct {
// 	bs   []byte
// 	pool sync.Pool
// }

// // AppendByte writes a single byte to the Buffer.
// func (b *Buffer) AppendByte(v byte) {
// 	b.bs = append(b.bs, v)
// }

// // AppendString writes a string to the Buffer.
// func (b *Buffer) AppendString(s string) {
// 	b.bs = append(b.bs, s...)
// }

// // AppendInt appends an integer to the underlying buffer (assuming base 10).
// func (b *Buffer) AppendInt(i int64) {
// 	b.bs = strconv.AppendInt(b.bs, i, 10)
// }

// // AppendTime appends the time formatted using the specified layout.
// func (b *Buffer) AppendTime(t time.Time, layout string) {
// 	b.bs = t.AppendFormat(b.bs, layout)
// }

// // AppendUint appends an unsigned integer to the underlying buffer (assuming
// // base 10).
// func (b *Buffer) AppendUint(i uint64) {
// 	b.bs = strconv.AppendUint(b.bs, i, 10)
// }

// // AppendBool appends a bool to the underlying buffer.
// func (b *Buffer) AppendBool(v bool) {
// 	b.bs = strconv.AppendBool(b.bs, v)
// }

// // AppendFloat appends a float to the underlying buffer. It doesn't quote NaN
// // or +/- Inf.
// func (b *Buffer) AppendFloat(f float64, bitSize int) {
// 	b.bs = strconv.AppendFloat(b.bs, f, 'f', -1, bitSize)
// }

// // Len returns the length of the underlying byte slice.
// func (b *Buffer) Len() int {
// 	return len(b.bs)
// }

// // Cap returns the capacity of the underlying byte slice.
// func (b *Buffer) Cap() int {
// 	return cap(b.bs)
// }

// // Bytes returns a mutable reference to the underlying byte slice.
// func (b *Buffer) Bytes() []byte {
// 	return b.bs
// }

// // String returns a string copy of the underlying byte slice.
// func (b *Buffer) String() string {
// 	return unsafe.String(unsafe.SliceData(b.bs), len(b.bs))
// }

// // Reset resets the underlying byte slice. Subsequent writes re-use the slice's
// // backing array.
// func (b *Buffer) Reset() {
// 	b.bs = b.bs[:0]
// }

// // Write implements io.Writer.
// func (b *Buffer) Write(bs []byte) (int, error) {
// 	b.bs = append(b.bs, bs...)
// 	return len(bs), nil
// }

// // WriteByte writes a single byte to the Buffer.
// //
// // Error returned is always nil, function signature is compatible
// // with bytes.Buffer and bufio.Writer
// func (b *Buffer) WriteByte(v byte) error {
// 	b.AppendByte(v)
// 	return nil
// }

// // WriteString writes a string to the Buffer.
// //
// // Error returned is always nil, function signature is compatible
// // with bytes.Buffer and bufio.Writer
// func (b *Buffer) WriteString(s string) (int, error) {
// 	b.AppendString(s)
// 	return len(s), nil
// }

// func (b *Buffer) WriteInterface(value any) {
// 	switch fValue := value.(type) {
// 	case string:
// 		if needsQuote(fValue) {
// 			b.WriteString(strconv.Quote(fValue))
// 		} else {
// 			b.WriteString(fValue)
// 		}
// 	case int:
// 		b.AppendInt(int64(fValue))
// 	case int8:
// 		b.AppendInt(int64(fValue))
// 	case int16:
// 		b.AppendInt(int64(fValue))
// 	case int32:
// 		b.AppendInt(int64(fValue))
// 	case int64:
// 		b.AppendInt(int64(fValue))
// 	case uint:
// 		b.AppendUint(uint64(fValue))
// 	case uint8:
// 		b.AppendUint(uint64(fValue))
// 	case uint16:
// 		b.AppendUint(uint64(fValue))
// 	case uint32:
// 		b.AppendUint(uint64(fValue))
// 	case uint64:
// 		b.AppendUint(uint64(fValue))
// 		b.AppendUint(uint64(fValue))
// 	case float32:
// 		b.AppendFloat(float64(fValue), 64)
// 	case float64:
// 		b.AppendFloat(float64(fValue), 64)
// 	case bool:
// 		b.AppendBool(fValue)
// 	case error:
// 		b.WriteString(fValue.Error())
// 	case []byte:
// 		b.Write(fValue)
// 	case time.Time:
// 		b.AppendTime(fValue, time.RFC3339Nano)
// 	case time.Duration:
// 		b.AppendString(fValue.String())
// 	case json.Number:
// 		b.AppendString(fValue.String())
// 	default:
// 		js, err := json.Marshal(fValue)
// 		if err != nil {
// 			fmt.Fprintln(os.Stderr, err)
// 		} else {
// 			b.Write(js)
// 		}
// 	}
// }

// // TrimNewline trims any final "\n" byte from the end of the buffer.
// func (b *Buffer) TrimNewline() {
// 	if i := len(b.bs) - 1; i >= 0 {
// 		if b.bs[i] == '\n' {
// 			b.bs = b.bs[:i]
// 		}
// 	}
// }

// // WriteNewLine writes a new line to the buffer if it's needed.
// func (b *Buffer) WriteNewLine() {
// 	if length := b.Len(); length > 0 && b.bs[length-1] != '\n' {
// 		b.WriteByte('\n') // nolint:errcheck
// 	}
// }

// func needsQuote(s string) bool {
// 	for i := range s {
// 		c := s[i]
// 		if c < 0x20 || c > 0x7e || c == ' ' || c == '\\' || c == '"' || c == '\n' || c == '\r' {
// 			return true
// 		}
// 	}
// 	return false
// }

// // bufPool is used to reuse issued-request buffers across writes to brokers.
// type bufPool struct{ p *sync.Pool }

// func newBufPool() bufPool {
// 	return bufPool{
// 		p: &sync.Pool{New: func() any { r := make([]byte, 1<<10); return &r }},
// 	}
// }

// func (p bufPool) get() []byte  { return (*p.p.Get().(*[]byte))[:0] }
// func (p bufPool) put(b []byte) { p.p.Put(&b) }
