package errors

import (
	"strings"
	"testing"
)

func TestSiweInvalidMessageFieldError_Error(t *testing.T) {
	err := &SiweInvalidMessageFieldError{Field: "nonce"}
	got := err.Error()
	if !strings.Contains(got, "Sign-In with Ethereum") || !strings.Contains(got, "nonce") {
		t.Errorf("Error() = %q, want to contain 'Sign-In with Ethereum' and field", got)
	}
}
