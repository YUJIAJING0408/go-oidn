package oidn

import "errors"

var (
	ErrLibraryNotFound = errors.New("oidn: dynamic library not found")
	ErrNotInitialized  = errors.New("oidn: library not initialized, call Init() first")
)
