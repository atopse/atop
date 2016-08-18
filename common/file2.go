package common

import (
	"os"
	"path/filepath"
	"runtime"
)

// SearchFile from the runtime
func SearchFile(name string) (string, error) {
	if name == "" {
		return "", os.ErrNotExist
	}
	f, _ := os.Getwd()
	if f == "" {
		_, file, _, _ := runtime.Caller(1)
		f = filepath.Dir(file)
	}
	dir := filepath.Dir(f)
	var fileName string
	for {
		if dir == filepath.VolumeName(dir) {
			return "", os.ErrNotExist
		}
		fileName = filepath.Join(dir, name)
		if _, err := os.Stat(fileName); err == nil {
			return fileName, nil
		} else if os.IsNotExist(err) == false {
			return "", err
		}
		// Parent dir.
		dir = filepath.Dir(dir)
		if dir == fileName {
			return "", os.ErrNotExist
		}
	}
}
