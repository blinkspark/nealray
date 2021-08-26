package util

import (
	"errors"
	"os"
)

func PathExists(fpath string) bool {
	_, err := os.Lstat(fpath)
	return !errors.Is(err, os.ErrNotExist)
}
