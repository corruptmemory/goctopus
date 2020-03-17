package monitor

import (
	"fmt"
	"os"
	"path/filepath"
)

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
func EnsureDir(path string, clean bool) error {
	if s, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0700); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		if !s.IsDir() {
			return fmt.Errorf("file exists at path: '%s'", path)
		} else {
			if clean {
				return RemoveContents(path)
			}
		}
	}
	return nil
}
