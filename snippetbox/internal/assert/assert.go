package assert

import "testing"

// func to test if two items are equivalent
// the types must be comparable
func Equal[T comparable](t *testing.T, actual, expected T) {
	// tells the testing module that this func is a test helper and not a func that should be tested itself
	// this will make sure the calling function is reported instead of this function itself
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}
