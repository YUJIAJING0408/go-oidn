//go:build windows

package oidn

import "syscall"

func loadLibrary() (uintptr, error) {
	names := []string{
		"lib/windows/OpenImageDenoise.dll",
	}
	for _, name := range names {
		h, err := syscall.LoadLibrary(name)
		if err == nil {
			return uintptr(h), nil
		}
	}
	return 0, ErrLibraryNotFound
}
