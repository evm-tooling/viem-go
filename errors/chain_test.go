package errors

import (
	"strings"
	"testing"
)

func TestChainDoesNotSupportContract_Error(t *testing.T) {
	blockNum := uint64(1000)
	blockCreated := uint64(2000)
	tests := []struct {
		name     string
		err      *ChainDoesNotSupportContract
		contains []string
	}{
		{"basic", &ChainDoesNotSupportContract{ChainName: "mainnet", ContractName: "multicall3"}, []string{"mainnet", "multicall3"}},
		{"block not deployed", &ChainDoesNotSupportContract{ChainName: "x", ContractName: "y", BlockNumber: &blockNum, BlockCreated: &blockCreated}, []string{"not deployed until block 2000", "current block 1000"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func TestChainMismatchError_Error(t *testing.T) {
	err := &ChainMismatchError{ChainID: 1, ChainName: "mainnet", CurrentChainID: 137}
	got := err.Error()
	if !strings.Contains(got, "1") || !strings.Contains(got, "137") || !strings.Contains(got, "mainnet") {
		t.Errorf("Error() = %q, want to contain chain IDs and name", got)
	}
}

func TestChainNotFoundError_Error(t *testing.T) {
	err := &ChainNotFoundError{}
	testErrorContains(t, err, []string{"No chain was provided"})
}

func TestClientChainNotConfiguredError_Error(t *testing.T) {
	err := &ClientChainNotConfiguredError{}
	testErrorContains(t, err, []string{"No chain was provided to the Client"})
}

func TestInvalidChainIdError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *InvalidChainIdError
		contains []string
	}{
		{"with id", &InvalidChainIdError{ChainID: int64Ptr(999)}, []string{"999", "invalid"}},
		{"without id", &InvalidChainIdError{}, []string{"invalid"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func int64Ptr(n int64) *int64 {
	return &n
}
