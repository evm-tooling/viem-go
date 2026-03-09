package errors

import (
	"strings"
	"testing"
)

func TestInvalidAddressError_Error(t *testing.T) {
	err := &InvalidAddressError{Address: "0xinvalid"}
	got := err.Error()
	if !strings.Contains(got, "0xinvalid") || !strings.Contains(got, "invalid") {
		t.Errorf("Error() = %q, want to contain address and 'invalid'", got)
	}
	if !strings.Contains(got, "40 hex characters") {
		t.Errorf("Error() = %q, want to contain hex/checksum hint", got)
	}
}
