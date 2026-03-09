package errors

import (
	"strings"
	"testing"
)

func TestAccountNotFoundError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *AccountNotFoundError
		contains []string
	}{
		{"basic", &AccountNotFoundError{}, []string{"Could not find an Account"}},
		{"with docs", &AccountNotFoundError{DocsPath: "/docs/account"}, []string{"Could not find", "viem.sh/docs"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func TestAccountTypeNotSupportedError_Error(t *testing.T) {
	err := &AccountTypeNotSupportedError{Type: "json-rpc"}
	got := err.Error()
	if !strings.Contains(got, "not supported") || !strings.Contains(got, "json-rpc") {
		t.Errorf("Error() = %q, want to contain 'not supported' and 'json-rpc'", got)
	}
}
