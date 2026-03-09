package errors

import (
	"math/big"
	"strings"
	"testing"
)

func TestExecutionRevertedError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ExecutionRevertedError
		contains []string
	}{
		{"with message", &ExecutionRevertedError{Message: "insufficient funds"}, []string{"Execution reverted", "insufficient funds"}},
		{"with execution reverted prefix", &ExecutionRevertedError{Message: "execution reverted: custom"}, []string{"Execution reverted", "custom"}},
		{"bare", &ExecutionRevertedError{}, []string{"Execution reverted", "unknown reason"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func TestFeeCapTooHighError_Error(t *testing.T) {
	gwei := big.NewInt(1e18)
	tests := []struct {
		name     string
		err      *FeeCapTooHighError
		contains []string
	}{
		{"with maxFeePerGas", &FeeCapTooHighError{MaxFeePerGas: gwei}, []string{"fee cap", "2^256"}},
		{"without", &FeeCapTooHighError{}, []string{"fee cap", "2^256"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func TestFeeCapTooLowError_Error(t *testing.T) {
	err := &FeeCapTooLowError{}
	testErrorContains(t, err, []string{"cannot be lower", "block base fee"})
}

func TestNonceTooHighError_Error(t *testing.T) {
	n := uint64(5)
	err := &NonceTooHighError{Nonce: &n}
	testErrorContains(t, err, []string{"higher than the next", "5"})
}

func TestNonceTooLowError_Error(t *testing.T) {
	err := &NonceTooLowError{}
	testErrorContains(t, err, []string{"lower than the current", "getTransactionCount"})
}

func TestNonceMaxValueError_Error(t *testing.T) {
	err := &NonceMaxValueError{}
	testErrorContains(t, err, []string{"exceeds the maximum"})
}

func TestInsufficientFundsError_Error(t *testing.T) {
	err := &InsufficientFundsError{}
	got := err.Error()
	if !strings.Contains(got, "exceeds the balance") {
		t.Errorf("Error() = %q, want to contain 'exceeds the balance'", got)
	}
	if !strings.Contains(got, "gas * gas fee") {
		t.Errorf("Error() = %q, want to contain 'gas * gas fee'", got)
	}
}

func TestIntrinsicGasTooHighError_Error(t *testing.T) {
	err := &IntrinsicGasTooHighError{}
	testErrorContains(t, err, []string{"exceeds the limit"})
}

func TestIntrinsicGasTooLowError_Error(t *testing.T) {
	err := &IntrinsicGasTooLowError{}
	testErrorContains(t, err, []string{"too low"})
}

func TestTransactionTypeNotSupportedError_Error(t *testing.T) {
	err := &TransactionTypeNotSupportedError{}
	testErrorContains(t, err, []string{"not supported"})
}

func TestTipAboveFeeCapError_Error(t *testing.T) {
	err := &TipAboveFeeCapError{}
	testErrorContains(t, err, []string{"tip", "maxPriorityFeePerGas", "maxFeePerGas"})
}

func TestUnknownNodeError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *UnknownNodeError
		contains []string
	}{
		{"with cause", &UnknownNodeError{Cause: &BaseError{ShortMessage: "rpc failed"}}, []string{"An error occurred", "rpc failed"}},
		{"without cause", &UnknownNodeError{}, []string{"An error occurred", "unknown"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}
