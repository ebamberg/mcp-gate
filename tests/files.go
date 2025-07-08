package tests

import (
	"errors"
	"os"
	"testing"
)

func AssertFileExist(t *testing.T, filename string) {
	_, err := os.Stat(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			t.Errorf("Failed: file doesn't exists: %s\n", err)
			return
		} else {
			t.Errorf("Failed: error testing file exists: %s\n", err)
		}
	}
}
