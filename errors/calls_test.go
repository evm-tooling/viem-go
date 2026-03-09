package errors

import (
	"strings"
	"testing"
)

func TestBundleFailedError_Error(t *testing.T) {
	err := &BundleFailedError{StatusCode: 500}
	got := err.Error()
	if !strings.Contains(got, "bundle failed") || !strings.Contains(got, "500") {
		t.Errorf("Error() = %q, want to contain 'bundle failed' and status", got)
	}
}
