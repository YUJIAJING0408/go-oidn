//go:build linux

package oidn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ebitengine/purego"
)

var defaultLibNames = []string{
	"libOpenImageDenoise.so",
	"libOpenImageDenoise.so.2",
}

func load(libDir string) (uintptr, error) {
	for _, name := range defaultLibNames {
		fullPath := filepath.Join(libDir, name)
		h, err := purego.Dlopen(fullPath, purego.RTLD_LAZY)
		if err == nil {
			return h, nil
		}
	}
	return 0, ErrLibraryNotFound
}

func loadLibrary() (uintptr, error) {
	candidates := []struct {
		source string
		path   string
	}{
		{"SetLibraryPath", libraryPath},
		{"OIDN_LIB_PATH", os.Getenv("OIDN_LIB_PATH")},
	}

	// 动态候选：当前目录、可执行目录、系统目录
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, struct{ source, path string }{"current dir", wd})
	}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, struct{ source, path string }{"executable dir", filepath.Dir(exe)})
	}
	// 常见系统安装路径
	for _, sysDir := range []string{"/usr/lib", "/usr/local/lib"} {
		candidates = append(candidates, struct{ source, path string }{"system dir", sysDir})
	}

	for _, c := range candidates {
		if c.path == "" {
			continue
		}
		h, err := load(c.path)
		if err == nil {
			return h, nil
		}
		// 只对非系统默认路径打印错误，避免干扰
		if c.source != "current dir" && c.source != "executable dir" && c.source != "system dir" {
			fmt.Fprintf(os.Stderr, "[oidn] %s load failed from %s: %v\n", c.source, c.path, err)
		}
	}

	return 0, ErrLibraryNotFound
}
