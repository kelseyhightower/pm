package cmd

import (
	"fmt"
	"runtime"
)

func FormatPackageName(name, tag, os, arch string) string {
	if os == "" {
		os = runtime.GOOS
	}
	if arch == "" {
		arch = runtime.GOARCH
	}
	return fmt.Sprintf("%s-%s-%s-%s.tar.gz", name, tag, os, arch)
}
