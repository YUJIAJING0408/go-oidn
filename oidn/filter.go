package oidn

import (
	"fmt"
	"unsafe"

	"github.com/YUJIAJING0408/go-oidn/internal"
)

// ProgressCallback is the function signature for filter progress callbacks.
type ProgressCallback func(userPtr unsafe.Pointer, n float64) uintptr

// Filter represents an OIDN filter with automatic reference counting.
type Filter struct {
	h uintptr
}

// NewFilter creates a new filter of the specified type (e.g. "RT").
func (d *Device) NewFilter(filterType string) (*Filter, error) {
	h := internal.F.NewFilter(d.h, cString(filterType))
	if h == 0 {
		return nil, fmt.Errorf("oidn: failed to create filter of type %q", filterType)
	}
	return &Filter{h: h}, nil
}

// Release releases the filter.
func (f *Filter) Release() {
	if f.h != 0 {
		internal.F.ReleaseFilter(f.h)
		f.h = 0
	}
}

// SetImage sets an image parameter of the filter.
func (f *Filter) SetImage(name string, buffer *Buffer, format Format, width, height int) {
	internal.F.SetFilterImage(f.h, cString(name), buffer.Handle(), int32(format),
		uintptr(width), uintptr(height), 0, 0, 0)
}

// SetBool sets a boolean parameter of the filter.
func (f *Filter) SetBool(name string, value bool) {
	internal.F.SetFilterBool(f.h, cString(name), value)
}

// SetInt sets an integer parameter of the filter.
func (f *Filter) SetInt(name string, value int) {
	internal.F.SetFilterInt(f.h, cString(name), int32(value))
}

// SetFloat sets a float parameter of the filter.
func (f *Filter) SetFloat(name string, value float32) {
	internal.F.SetFilterFloat(f.h, cString(name), value)
}

// Commit commits all previous changes to the filter.
// Must be called before first executing the filter.
func (f *Filter) Commit() {
	internal.F.CommitFilter(f.h)
}

// Execute runs the filter synchronously.
func (f *Filter) Execute() {
	internal.F.ExecuteFilter(f.h)
}

// Handle returns the underlying C filter handle.
func (f *Filter) Handle() uintptr {
	return f.h
}
