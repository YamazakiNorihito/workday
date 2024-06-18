package helper

import "testing"

func MustSucceed(t *testing.T, action func() error) {
	t.Helper()
	if err := action(); err != nil {
		panic("Action must succeed but failed with error: " + err.Error())
	}
}
