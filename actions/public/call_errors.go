package public

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/utils/errors"
)

// CallExecutionError is returned when a call execution fails.
// Re-exported from utils/errors for backwards compatibility.
type CallExecutionError = errors.CallExecutionError

// CounterfactualDeploymentFailedError is returned when a deployless call via
// factory fails to deploy the contract.
type CounterfactualDeploymentFailedError struct {
	Factory *common.Address
}

func (e *CounterfactualDeploymentFailedError) Error() string {
	if e.Factory != nil {
		return fmt.Sprintf("counterfactual deployment failed: factory=%s", e.Factory.Hex())
	}
	return "counterfactual deployment failed"
}

// RawContractError represents a raw contract revert error.
type RawContractError struct {
	Data []byte
}

func (e *RawContractError) Error() string {
	if len(e.Data) > 0 {
		return fmt.Sprintf("contract reverted with data: 0x%x", e.Data)
	}
	return "contract reverted"
}

// InvalidCallParamsError is returned when call parameters are invalid.
type InvalidCallParamsError struct {
	Message string
}

func (e *InvalidCallParamsError) Error() string {
	return fmt.Sprintf("invalid call parameters: %s", e.Message)
}

// ChainNotConfiguredError is returned when a chain is not configured on the client.
type ChainNotConfiguredError struct{}

func (e *ChainNotConfiguredError) Error() string {
	return "chain not configured on client"
}

// ChainDoesNotSupportContractError is returned when a chain doesn't support
// a required contract (e.g., multicall3).
type ChainDoesNotSupportContractError struct {
	ChainID      int64
	ContractName string
	BlockNumber  *uint64
}

func (e *ChainDoesNotSupportContractError) Error() string {
	if e.BlockNumber != nil {
		return fmt.Sprintf("chain %d does not support %s at block %d",
			e.ChainID, e.ContractName, *e.BlockNumber)
	}
	return fmt.Sprintf("chain %d does not support %s", e.ChainID, e.ContractName)
}
