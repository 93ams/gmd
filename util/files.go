package util

import (
	"errors"
	"os"
	"path/filepath"
)

func DesureFile(path string) error {
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		err = nil
	}
	return err
}
func EnsureFile(path string) error {
	dir, path := filepath.Split(path)
	if err := EnsureDir(dir); err != nil {
		return err
	}
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		file, err = os.Create(path)
	}
	if err != nil {
		return err
	}
	return file.Close()
}

func EnsureDir(path string) error { return os.MkdirAll(path, 0655) }
