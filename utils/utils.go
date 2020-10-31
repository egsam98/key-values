package utils

import (
	"path/filepath"
	"runtime"
	"strings"
)

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	return strings.TrimSuffix(filepath.Dir(b), "utils")
}
