package public

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/types"
)

// EstimateContractGasParameters contains the parameters for the
// EstimateContractGas action.
//
// This mirrors viem's EstimateContractGasParameters type in a Go-friendly form.
type EstimateContractGasParameters struct {
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

	// Gas is an optional gas limit hint for the estimation request.
	Gas *uint64

	// GasPrice is the legacy gas price.
	GasPrice *big.Int

	// MaxFeePerGas is the max fee per gas (EIP-1559).
	MaxFeePerGas *big.Int

	// MaxPriorityFeePerGas is the max priority fee per gas (EIP-1559).
	MaxPriorityFeePerGas *big.Int

	// Nonce is the transaction nonce.
	Nonce *uint64

	// AccessList is the EIP-2930 access list.
	AccessList types.AccessList

	// StateOverride contains state overrides for the estimation.
	StateOverride types.StateOverride

	// BlockNumber is the block number to estimate at.
	// Mutually exclusive with BlockTag.
	BlockNumber *uint64

	// BlockTag is the block tag to estimate at (e.g., "latest", "pending").
	// Mutually exclusive with BlockNumber.
	BlockTag BlockTag

	// DataSuffix is optional data to append to the end of the calldata.
	// Useful for adding a "domain" tag.
	DataSuffix []byte
}

// EstimateContractGasReturnType is the return type for the EstimateContractGas
// action. It represents the gas estimate in units of gas.
type EstimateContractGasReturnType = uint64

// EstimateContractGas estimates the gas required to successfully execute a
// contract write function call.
//
// Internally, this uses the EstimateGas action with ABI-encoded calldata,
// mirroring viem's `estimateContractGas` implementation.
func EstimateContractGas(
	ctx context.Context,
	client Client,
	params EstimateContractGasParameters,
) (EstimateContractGasReturnType, error) {
	if params.ABI == nil {
		return 0, fmt.Errorf("ABI is required")
	}
	if params.FunctionName == "" {
		return 0, fmt.Errorf("functionName is required")
	}

	// Encode function data.
	calldata, err := params.ABI.EncodeFunctionData(params.FunctionName, params.Args...)
	if err != nil {
		return 0, fmt.Errorf("failed to encode function data: %w", err)
	}

	// Append data suffix if provided.
	if len(params.DataSuffix) > 0 {
		calldata = append(calldata, params.DataSuffix...)
	}

	// Delegate to EstimateGas with the encoded calldata.
	gas, err := EstimateGas(ctx, client, EstimateGasParameters{
		Account:              params.Account,
		To:                   &params.Address,
		Data:                 calldata,
		Value:                params.Value,
		Gas:                  params.Gas,
		GasPrice:             params.GasPrice,
		MaxFeePerGas:         params.MaxFeePerGas,
		MaxPriorityFeePerGas: params.MaxPriorityFeePerGas,
		MaxFeePerBlobGas:     nil, // can be set by caller via EstimateGasParameters if needed
		Nonce:                params.Nonce,
		AccessList:           params.AccessList,
		StateOverride:        params.StateOverride,
		BlockNumber:          params.BlockNumber,
		BlockTag:             params.BlockTag,
	})
	if err != nil {
		return 0, fmt.Errorf("contract gas estimation failed for %s.%s: %w", params.Address.Hex(), params.FunctionName, err)
	}

	return gas, nil
}
