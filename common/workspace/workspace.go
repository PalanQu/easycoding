package workspace

import (
	"path/filepath"
	"runtime"
)

func GetWorkspace() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../..")
}
