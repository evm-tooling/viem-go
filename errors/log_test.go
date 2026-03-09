package errors

import (
	"strings"
	"testing"
)

func TestFilterTypeNotSupportedError_Error(t *testing.T) {
	err := &FilterTypeNotSupportedError{Type: "foo"}
	got := err.Error()
	if !strings.Contains(got, "not supported") || !strings.Contains(got, "foo") {
		t.Errorf("Error() = %q, want to contain 'not supported' and 'foo'", got)
	}
}
