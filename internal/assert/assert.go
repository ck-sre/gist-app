package assert

import "testing"

func Same[T comparable](t *testing.T, got, exp T) {
	t.Helper()
	if exp != got {
		t.Errorf("Expected %v, got %v", exp, got)
	}
}
