//go:build windows

package oidn

import (
	"fmt"
	"os"
	"path"
	"syscall"
)

var defaultLibNames = "OpenImageDenoise.dll"

func load(libDir string) (uintptr, error) {
	h, err := syscall.LoadLibrary(path.Join(libDir, version, defaultLibNames))
	if err == nil {
		return uintptr(h), nil
	}
	return uintptr(0), err
}

func loadLibrary() (res uintptr, err error) {
	if libraryPath != "" {
		if res, err = load(libraryPath); err == nil {
			return
		}
		// Failed, Not Panic But Try Others
		fmt.Fprintf(os.Stderr, "Error loading library: %s so %s\n", "set library path maybe wrong", err)
	}
	if env := os.Getenv("OIDN_LIB_PATH"); env != "" {
		if res, err = load(env); err == nil {
			return
		}
		// Failed, Not Panic But Try Others
		fmt.Fprintf(os.Stderr, "Error loading library: %s so %s\n", "OIDN_LIB_PATH maybe wrong", err)
	}
	if res, err = load("lib/windows"); err == nil {
		return
	}
	return 0, ErrLibraryNotFound
}
