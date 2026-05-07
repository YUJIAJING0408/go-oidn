package oidn

import "errors"

var (
	ErrLibraryNotFound = errors.New("oidn: dynamic library not found")
	ErrNotInitialized  = errors.New("oidn: library not initialized, call Init() first")
)

var libraryPath string
var version = "v2.4.1"

func SetLibraryPath(path string) {
	libraryPath = path
}
func SetVersion(v string) {
	version = v
}
