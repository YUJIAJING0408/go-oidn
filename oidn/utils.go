package oidn

import "unsafe"

// cString converts a Go string to a null-terminated C string (*byte).
// The returned pointer is valid for the duration of the call it's used in.
func cString(s string) *byte {
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	return &buf[0]
}

// goString converts a C string pointer to a Go string.
func goString(cstr *byte) string {
	if cstr == nil {
		return ""
	}
	ptr := unsafe.Pointer(cstr)
	var length int
	for {
		if *(*byte)(unsafe.Add(ptr, length)) == 0 {
			break
		}
		length++
	}
	return string(unsafe.Slice(cstr, length))
}
