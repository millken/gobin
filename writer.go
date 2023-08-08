package gobin

import (
	"encoding/binary"
	"errors"
)

var (
	// ErrOutOfMemory happens when alloc is required while using NewFixedWriter.
	ErrOutOfMemory = errors.New("out-of-memory, FixedWriter can't reallocate")
)

// Writer holds the encoded, the finished encode can be retrieved by Writer.Bytes()
type Writer struct {
	Memory  []byte
	isFixed bool
	endian  binary.ByteOrder
}

// NewWriter creates a Writer with the given initial capacity.
func NewWriter(capacity int) *Writer {
	return CreateWriter(capacity, false, binary.LittleEndian)
}

// CreateWriter creates a Writer with the given initial capacity, fixed size and endian.
func CreateWriter(capacity int, isFixed bool, endian binary.ByteOrder) *Writer {
	return &Writer{Memory: make([]byte, 0, capacity), isFixed: isFixed, endian: endian}
}

// Alloc allocates n bytes inside.
// It returns the offset and may return error if it's not possible to allocate.
func (w *Writer) Alloc(n uint) (uint, error) {
	ptr := uint(len(w.Memory))
	total := ptr + n
	if total > uint(cap(w.Memory)) {
		if w.isFixed {
			return 0, ErrOutOfMemory
		}
		w.Memory = append(w.Memory, make([]byte, total-uint(len(w.Memory)))...)
	} else {
		w.Memory = w.Memory[:total]
	}
	return ptr, nil
}

// WriteAt copies the given data into the Writer memory.
func (w *Writer) WriteAt(offset uint, data []byte) {
	copy(w.Memory[offset:], data)
}

// Reset will reset the memory length, but keeps the memory capacity.
func (w *Writer) Reset() {
	if len(w.Memory) > 0 {
		w.Memory = w.Memory[:0]
	}
}

func (w *Writer) WriteUint8(offset uint, i uint8) {
	w.Memory[offset] = i
}

func (w *Writer) WriteUint16(offset uint, i uint16) {
	w.endian.PutUint16(w.Memory[offset:], i)
}

func (w *Writer) WriteUint32(offset uint, i uint32) {
	w.endian.PutUint32(w.Memory[offset:], i)
}

func (w *Writer) WriteUint64(offset uint, i uint64) {
	w.endian.PutUint64(w.Memory[offset:], i)
}

// Bytes return the Karmem encoded bytes.
// It doesn't copy the content, and can't be re-used after Reset.
func (w *Writer) Bytes() []byte {
	return w.Memory
}
