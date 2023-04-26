package os

import (
	"os"
	"path/filepath"
	"strings"
)

func MustGetAppName() string {
	exec, err := os.Executable()
	if err != nil {
		panic(err)
	}
	_, file := filepath.Split(exec)
	return strings.TrimSuffix(file, filepath.Ext(file))
}
