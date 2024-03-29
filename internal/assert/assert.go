package assert

import (
	"strings"
	"testing"
)

func Same[T comparable](t *testing.T, got, exp T) {
	t.Helper()
	if exp != got {
		t.Errorf("Expected %v, got %v", exp, got)
	}
}

func StringHas(t *testing.T, real, expected string) {
	t.Helper()
	if !strings.Contains(real, expected) {
		t.Errorf("Got: %q, Expected %q", real, expected)
	}
}

func NilError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Errorf("got: %v; expected nil", got)
	}
}
