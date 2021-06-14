package util

import "os"

func PathValid(fpath string) bool {
	_, err := os.Stat(fpath)
	return err == nil
}
