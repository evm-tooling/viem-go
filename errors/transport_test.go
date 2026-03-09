package errors

import (
	"strings"
	"testing"
)

func TestUrlRequiredError_Error(t *testing.T) {
	err := &UrlRequiredError{}
	got := err.Error()
	if !strings.Contains(got, "No URL was provided") {
		t.Errorf("Error() = %q, want to contain 'No URL was provided'", got)
	}
	if !strings.Contains(got, "viem.sh/docs") {
		t.Errorf("Error() = %q, want to contain docs URL", got)
	}
}
