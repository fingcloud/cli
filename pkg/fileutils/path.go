package fileutils

import (
	"os"
)

func FileExists(path string) bool {
	_, err := os.Open(path)
	return os.IsNotExist(err)
}
