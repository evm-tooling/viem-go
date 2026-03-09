package errors

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestCallExecutionError_Error(t *testing.T) {
	addr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	tests := []struct {
		name     string
		err      *CallExecutionError
		contains []string
	}{
		{"with message", &CallExecutionError{Message: "revert"}, []string{"call execution", "revert"}},
		{"with cause", &CallExecutionError{Cause: &BaseError{ShortMessage: "rpc error"}}, []string{"rpc error"}},
		{"bare", &CallExecutionError{}, []string{"call execution failed"}},
		{"with to", &CallExecutionError{Message: "x", To: &addr}, []string{"call execution", "x"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func TestCallExecutionError_Unwrap(t *testing.T) {
	cause := &BaseError{ShortMessage: "cause"}
	err := &CallExecutionError{Cause: cause}
	if err.Unwrap() != cause {
		t.Error("Unwrap should return cause")
	}
}

func TestContractFunctionExecutionError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ContractFunctionExecutionError
		contains []string
	}{
		{"with cause", &ContractFunctionExecutionError{Cause: &BaseError{ShortMessage: "revert"}, FunctionName: "foo"}, []string{"revert"}},
		{"without cause", &ContractFunctionExecutionError{FunctionName: "balanceOf"}, []string{"balanceOf", "unknown error"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func TestContractFunctionRevertedError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ContractFunctionRevertedError
		contains []string
	}{
		{"with reason", &ContractFunctionRevertedError{FunctionName: "transfer", Reason: "insufficient balance"}, []string{"reverted", "insufficient balance"}},
		{"with signature", &ContractFunctionRevertedError{FunctionName: "foo", Signature: "0xabcd"}, []string{"reverted", "0xabcd"}},
		{"bare", &ContractFunctionRevertedError{FunctionName: "bar"}, []string{"bar", "reverted"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func TestContractFunctionZeroDataError_Error(t *testing.T) {
	err := &ContractFunctionZeroDataError{FunctionName: "getData"}
	testErrorContains(t, err, []string{"returned no data", "getData", "0x"})
}

func TestCounterfactualDeploymentFailedError_Error(t *testing.T) {
	addr := common.HexToAddress("0xabc")
	tests := []struct {
		name     string
		err      *CounterfactualDeploymentFailedError
		contains []string
	}{
		{"with factory", &CounterfactualDeploymentFailedError{Factory: &addr}, []string{"counterfactual", "factory", "0x"}},
		{"without factory", &CounterfactualDeploymentFailedError{}, []string{"counterfactual", "Please ensure"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func TestRawContractError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *RawContractError
		contains []string
	}{
		{"with message", &RawContractError{Message: "custom revert"}, []string{"custom revert"}},
		{"with data", &RawContractError{Data: []byte{0x12, 0x34}}, []string{"reverted", "0x1234"}},
		{"bare", &RawContractError{}, []string{"contract reverted"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}
