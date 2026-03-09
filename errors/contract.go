package errors

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// CallExecutionError is returned when a call execution fails.
type CallExecutionError struct {
	Cause   error
	Message string
	To      *common.Address
	Data    []byte
}

func (e *CallExecutionError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("call execution failed: %s", e.Message)
	}
	if e.Cause != nil {
		return fmt.Sprintf("call execution failed: %v", e.Cause)
	}
	return "call execution failed"
}

func (e *CallExecutionError) Unwrap() error {
	return e.Cause
}

// ContractFunctionExecutionError is returned when a contract function execution fails.
type ContractFunctionExecutionError struct {
	Cause           error
	FunctionName    string
	ContractAddress *common.Address
	Sender          *common.Address
	MetaMessages    []string
}

func (e *ContractFunctionExecutionError) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return fmt.Sprintf("An unknown error occurred while executing the contract function \"%s\".", e.FunctionName)
}

func (e *ContractFunctionExecutionError) Unwrap() error {
	return e.Cause
}

// ContractFunctionRevertedError is returned when a contract function reverts.
// ErrorName is the decoded custom error name (e.g. "ERC20InsufficientBalance", "Error", "Panic").
type ContractFunctionRevertedError struct {
	FunctionName string
	ErrorName    string // Decoded error name for custom error extraction (parity with viem)
	Reason       string
	Signature    string
	Data         interface{}
	Raw          []byte
}

func (e *ContractFunctionRevertedError) Error() string {
	if e.Reason != "" && e.Reason != "execution reverted" {
		return fmt.Sprintf("The contract function \"%s\" reverted with the following reason:\n%s", e.FunctionName, e.Reason)
	}
	if e.Signature != "" {
		return fmt.Sprintf("The contract function \"%s\" reverted with the following signature:\n%s", e.FunctionName, e.Signature)
	}
	return fmt.Sprintf("The contract function \"%s\" reverted.", e.FunctionName)
}

// ContractFunctionZeroDataError is returned when a contract function returns no data.
type ContractFunctionZeroDataError struct {
	FunctionName string
}

func (e *ContractFunctionZeroDataError) Error() string {
	msg := fmt.Sprintf("The contract function \"%s\" returned no data (\"0x\").", e.FunctionName)
	msg += "\nThis could be due to any of the following:"
	msg += fmt.Sprintf("\n  - The contract does not have the function \"%s\",", e.FunctionName)
	msg += "\n  - The parameters passed to the contract function may be invalid, or"
	msg += "\n  - The address is not a contract."
	return msg
}

// CounterfactualDeploymentFailedError is returned when counterfactual deployment fails.
type CounterfactualDeploymentFailedError struct {
	Factory *common.Address
}

func (e *CounterfactualDeploymentFailedError) Error() string {
	msg := "Deployment for counterfactual contract call failed"
	if e.Factory != nil {
		msg += fmt.Sprintf(" for factory \"%s\".", e.Factory.Hex())
	} else {
		msg += "."
	}
	msg += "\nPlease ensure:"
	msg += "\n- The `factory` is a valid contract deployment factory (ie. Create2 Factory, ERC-4337 Factory, etc)."
	msg += "\n- The `factoryData` is a valid encoded function call for contract deployment function on the factory."
	return msg
}

// RawContractError represents a raw contract revert error.
type RawContractError struct {
	Data    []byte
	Message string
}

func (e *RawContractError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if len(e.Data) > 0 {
		return fmt.Sprintf("contract reverted with data: 0x%x", e.Data)
	}
	return "contract reverted"
}
