package oidn

import (
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/YUJIAJING0408/go-oidn/internal"
)

// Buffer represents an OIDN buffer with automatic reference counting.
type Buffer struct {
	h uintptr
}

// NewBuffer creates a new buffer of the given size in bytes (host storage by default).
func (d *Device) NewBuffer(byteSize int) (*Buffer, error) {
	h := internal.F.NewBuffer(d.h, uintptr(byteSize))
	if h == 0 {
		return nil, fmt.Errorf("oidn: failed to create buffer")
	}
	return &Buffer{h: h}, nil
}

// Release releases the buffer.
func (b *Buffer) Release() {
	if b.h != 0 {
		internal.F.ReleaseBuffer(b.h)
		b.h = 0
	}
}

// Write copies data from the host slice into the buffer.
func (b *Buffer) Write(data []float32) {
	size := uintptr(len(data)) * uintptr(binary.Size(float32(0)))
	internal.F.WriteBuffer(b.h, 0, size, unsafe.Pointer(&data[0]))
}

// Read reads data from the buffer into a newly allocated float32 slice.
func (b *Buffer) Read(count int) []float32 {
	data := make([]float32, count)
	size := uintptr(count) * uintptr(binary.Size(float32(0)))
	internal.F.ReadBuffer(b.h, 0, size, unsafe.Pointer(&data[0]))
	return data
}

// Handle returns the underlying C buffer handle.
func (b *Buffer) Handle() uintptr {
	return b.h
}
