package tests

import (
	"testing"
)

func FailOnError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("error creating testfile %s\n", err)
	}
}
