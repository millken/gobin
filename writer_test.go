package gobin

import (
	"strconv"
	"sync"
	"testing"
)

func TestXxx(t *testing.T) {
	const minPoolExpOf2 = 8
	getHumanReadableSize := func(size int) string {
		if size < 1024 {
			return strconv.Itoa(size) + "B"
		}
		size /= 1024
		if size < 1024 {
			return strconv.Itoa(size) + "KB"
		}
		size /= 1024
		if size < 1024 {
			return strconv.Itoa(size) + "MB"
		}
		size /= 1024
		return strconv.Itoa(size) + "GB"
	}

	var pools [18]*sync.Pool
	for i := range pools {
		bufLen := 1 << (minPoolExpOf2 + i)
		t.Logf("bufSize: %s", getHumanReadableSize(bufLen))
		pools[i] = &sync.Pool{
			New: func() any {
				buf := make([]byte, bufLen)
				return &buf
			},
		}
	}
}
