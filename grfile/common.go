package grfile

import (
	"os"
	"path/filepath"
)

func IsDirExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}

		return false
	}
	return true
}
func MakeFilePathDirIfNotExist(path string) (err error) {
	dir := filepath.Dir(path)
	if !IsDirExist(dir) {
		err := os.MkdirAll(dir, 0775)

		if err != nil {

			return err

		}
	}
	return

}
