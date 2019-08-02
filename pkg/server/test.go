package server

import "testing"

func assertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Errorf("Expected %v, but got: %v", expected, actual)
	}
}
