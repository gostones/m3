package misc

import (
	"fmt"
	"os"
)

// CreateDir returns true if dir does not exist and was created successfully
// or false if it already exists; otherwise error
func CreateDir(dir string) (bool, error) {
	fi, err := os.Stat(dir)

	if os.IsNotExist(err) {
		if errDir := os.MkdirAll(dir, 0755); errDir != nil {
			return false, errDir
		}
		return true, nil
	}

	if fi.Mode().IsRegular() {
		return false, fmt.Errorf("%v exists as file", dir)
	}

	return false, nil
}

// ExistRegularFile check file exists and is regular
func ExistRegularFile(filename string) bool {
	if fi, err := os.Stat(filename); os.IsNotExist(err) || !fi.Mode().IsRegular() {
		return false
	} else {
		return true
	}
}
