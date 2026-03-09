package errors

import (
	"strings"
	"testing"
)

func TestEip712DomainNotFoundError_Error(t *testing.T) {
	err := &Eip712DomainNotFoundError{Address: "0x1234567890123456789012345678901234567890"}
	got := err.Error()
	if !strings.Contains(got, "EIP-712 domain") || !strings.Contains(got, "eip712Domain") {
		t.Errorf("Error() = %q, want to contain 'EIP-712 domain' and 'eip712Domain'", got)
	}
}
