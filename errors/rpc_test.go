package errors

import (
	"strings"
	"testing"
)

func TestRpcError_Error(t *testing.T) {
	cause := &BaseError{ShortMessage: "rpc failure"}
	err := NewParseRpcError(cause)
	testErrorContains(t, err, []string{"Invalid JSON", "rpc failure"})
}

func TestRpcError_Unwrap(t *testing.T) {
	cause := &BaseError{ShortMessage: "x"}
	err := NewInternalRpcError(cause)
	if err.Unwrap() != cause {
		t.Error("Unwrap should return cause")
	}
}

func TestNewParseRpcError(t *testing.T) {
	err := NewParseRpcError(nil)
	testErrorContains(t, err, []string{"Invalid JSON", "parsing"})
	if err.Code != RpcCodeParse {
		t.Errorf("Code = %d, want %d", err.Code, RpcCodeParse)
	}
}

func TestNewTransactionRejectedRpcError(t *testing.T) {
	err := NewTransactionRejectedRpcError(nil)
	testErrorContains(t, err, []string{"Transaction creation failed"})
}

func TestNewUserRejectedRequestError(t *testing.T) {
	err := NewUserRejectedRequestError(nil)
	testErrorContains(t, err, []string{"User rejected"})
}

func TestNewProviderDisconnectedError(t *testing.T) {
	err := NewProviderDisconnectedError(nil)
	testErrorContains(t, err, []string{"disconnected from all chains"})
}

func TestNewUnknownRpcError(t *testing.T) {
	err := NewUnknownRpcError(nil)
	got := err.Error()
	if !strings.Contains(got, "unknown RPC error") {
		t.Errorf("Error() = %q, want to contain 'unknown RPC error'", got)
	}
}
