package public

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/types"
)

// SimulateContractParameters contains the parameters for the SimulateContract action.
// This mirrors viem's SimulateContractParameters type.
type SimulateContractParameters struct {
	// Account is the account attached to the call (msg.sender).
	Account *common.Address

	// Address is the contract address.
	Address common.Address

	// ABI is the contract ABI.
	ABI *abi.ABI

	// FunctionName is the name of the function to call.
	FunctionName string

	// Args are the function arguments.
	Args []any

	// Value is the amount of wei to send with the call.
	Value *big.Int

	// Gas is the gas limit for the call.
	Gas *uint64

	// GasPrice is the gas price for the call (legacy).
	GasPrice *big.Int

	// MaxFeePerGas is the max fee per gas (EIP-1559).
	MaxFeePerGas *big.Int

	// MaxPriorityFeePerGas is the max priority fee per gas (EIP-1559).
	MaxPriorityFeePerGas *big.Int

	// Nonce is the nonce for the call.
	Nonce *uint64

	// BlockNumber is the block number to simulate at.
	// Mutually exclusive with BlockTag.
	BlockNumber *uint64

	// BlockTag is the block tag to simulate at (e.g., "latest", "pending").
	// Mutually exclusive with BlockNumber.
	BlockTag BlockTag

	// StateOverride contains state overrides for the call.
	StateOverride types.StateOverride

	// BlockOverrides contains block-level overrides.
	BlockOverrides *types.BlockOverrides

	// AccessList is the EIP-2930 access list.
	AccessList types.AccessList

	// DataSuffix is optional data to append to the calldata.
	// Useful for adding a "domain" tag.
	DataSuffix []byte
}

// SimulateContractRequest contains the request data for a simulated contract call.
type SimulateContractRequest struct {
	// ABI is the minimized ABI containing only the called function.
	ABI *abi.ABI

	// Address is the contract address.
	Address common.Address

	// FunctionName is the function name.
	FunctionName string

	// Args are the function arguments.
	Args []any

	// Account is the account used for the call.
	Account *common.Address

	// DataSuffix is the optional data suffix.
	DataSuffix []byte

	// Value is the amount of wei to send.
	Value *big.Int

	// Gas is the gas limit.
	Gas *uint64

	// GasPrice is the gas price.
	GasPrice *big.Int

	// MaxFeePerGas is the max fee per gas.
	MaxFeePerGas *big.Int

	// MaxPriorityFeePerGas is the max priority fee per gas.
	MaxPriorityFeePerGas *big.Int

	// Nonce is the nonce.
	Nonce *uint64

	// AccessList is the access list.
	AccessList types.AccessList

	// StateOverride contains state overrides.
	StateOverride types.StateOverride

	// BlockOverrides contains block overrides.
	BlockOverrides *types.BlockOverrides

	// BlockNumber is the block number.
	BlockNumber *uint64

	// BlockTag is the block tag.
	BlockTag BlockTag
}

// SimulateContractReturnType is the return type for the SimulateContract action.
type SimulateContractReturnType struct {
	// Result is the decoded return value from the contract call.
	Result any

	// Request contains the request data that can be used for writeContract.
	Request SimulateContractRequest
}

// SimulateContract simulates/validates a contract interaction.
//
// This is useful for retrieving return data and revert reasons of contract write functions.
// This function does not require gas to execute and does not change the state of the blockchain.
// It is almost identical to ReadContract, but also supports contract write functions.
//
// This is equivalent to viem's `simulateContract` action.
// Internally uses eth_call with ABI-encoded data.
//
// Example:
//
//	erc20ABI, _ := abi.ParseABI(`[{"inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"type":"bool"}],"type":"function"}]`)
//
//	result, err := public.SimulateContract(ctx, client, public.SimulateContractParameters{
//	    Account:      &senderAddr,
//	    Address:      tokenAddr,
//	    ABI:          erc20ABI,
//	    FunctionName: "transfer",
//	    Args:         []any{recipientAddr, big.NewInt(1000000)},
//	})
//
//	if result.Result.(bool) {
//	    // Transfer would succeed, can proceed with actual transaction
//	    fmt.Println("Simulation successful!")
//	}
func SimulateContract(ctx context.Context, client Client, params SimulateContractParameters) (*SimulateContractReturnType, error) {
	if params.ABI == nil {
		return nil, fmt.Errorf("ABI is required")
	}
	if params.FunctionName == "" {
		return nil, fmt.Errorf("functionName is required")
	}

	// Encode function data
	calldata, err := params.ABI.EncodeFunctionData(params.FunctionName, params.Args...)
	if err != nil {
		return nil, &SimulateContractError{
			Cause:        err,
			Address:      params.Address,
			FunctionName: params.FunctionName,
			Args:         params.Args,
		}
	}

	// Append data suffix if provided
	if len(params.DataSuffix) > 0 {
		calldata = append(calldata, params.DataSuffix...)
	}

	// Execute the call
	callParams := CallParameters{
		Account:              params.Account,
		To:                   &params.Address,
		Data:                 calldata,
		Value:                params.Value,
		Gas:                  params.Gas,
		GasPrice:             params.GasPrice,
		MaxFeePerGas:         params.MaxFeePerGas,
		MaxPriorityFeePerGas: params.MaxPriorityFeePerGas,
		Nonce:                params.Nonce,
		BlockNumber:          params.BlockNumber,
		BlockTag:             params.BlockTag,
		StateOverride:        params.StateOverride,
		BlockOverrides:       params.BlockOverrides,
		AccessList:           params.AccessList,
		Batch:                ptr(false), // Disable batching for simulate
	}

	callResult, err := Call(ctx, client, callParams)
	if err != nil {
		return nil, &SimulateContractError{
			Cause:        err,
			Address:      params.Address,
			FunctionName: params.FunctionName,
			Args:         params.Args,
		}
	}

	// Decode the result
	var result any
	if len(callResult.Data) > 0 {
		decoded, decodeErr := params.ABI.DecodeFunctionResult(params.FunctionName, callResult.Data)
		if decodeErr != nil {
			return nil, &SimulateContractError{
				Cause:        fmt.Errorf("failed to decode result: %w", decodeErr),
				Address:      params.Address,
				FunctionName: params.FunctionName,
				Args:         params.Args,
			}
		}

		// If single return value, unwrap it
		if len(decoded) == 1 {
			result = decoded[0]
		} else {
			result = decoded
		}
	}

	return &SimulateContractReturnType{
		Result: result,
		Request: SimulateContractRequest{
			ABI:                  params.ABI,
			Address:              params.Address,
			FunctionName:         params.FunctionName,
			Args:                 params.Args,
			Account:              params.Account,
			DataSuffix:           params.DataSuffix,
			Value:                params.Value,
			Gas:                  params.Gas,
			GasPrice:             params.GasPrice,
			MaxFeePerGas:         params.MaxFeePerGas,
			MaxPriorityFeePerGas: params.MaxPriorityFeePerGas,
			Nonce:                params.Nonce,
			AccessList:           params.AccessList,
			StateOverride:        params.StateOverride,
			BlockOverrides:       params.BlockOverrides,
			BlockNumber:          params.BlockNumber,
			BlockTag:             params.BlockTag,
		},
	}, nil
}

// SimulateContractError is returned when contract simulation fails.
type SimulateContractError struct {
	Cause        error
	Address      common.Address
	FunctionName string
	Args         []any
}

func (e *SimulateContractError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("contract simulation failed for %s.%s: %v", e.Address.Hex(), e.FunctionName, e.Cause)
	}
	return fmt.Sprintf("contract simulation failed for %s.%s", e.Address.Hex(), e.FunctionName)
}

func (e *SimulateContractError) Unwrap() error {
	return e.Cause
}
