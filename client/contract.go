package client

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/types"
)

// ReadContractOptions contains options for reading from a contract.
// This mirrors viem's readContract API.
type ReadContractOptions struct {
	// Address is the contract address.
	Address common.Address
	// ABI is the contract ABI as a JSON string or []byte.
	ABI any
	// FunctionName is the name of the function to call.
	FunctionName string
	// Args are the function arguments.
	Args []any
	// From is the address to use as the caller (optional).
	From *common.Address
	// Block is the block tag to read from (default: latest).
	Block BlockTag
}

// PrepareContractWriteOptions contains options for preparing a contract write transaction.
type PrepareContractWriteOptions struct {
	// Address is the contract address.
	Address common.Address
	// ABI is the contract ABI as a JSON string or []byte.
	ABI any
	// FunctionName is the name of the function to call.
	FunctionName string
	// Args are the function arguments.
	Args []any
	// From is the address sending the transaction.
	From common.Address
	// Value is the amount of ETH to send with the transaction.
	Value *big.Int
	// Gas is the gas limit (optional, will be estimated if not provided).
	Gas uint64
	// GasPrice is the gas price for legacy transactions.
	GasPrice *big.Int
	// MaxFeePerGas is the max fee per gas for EIP-1559 transactions.
	MaxFeePerGas *big.Int
	// MaxPriorityFeePerGas is the max priority fee for EIP-1559 transactions.
	MaxPriorityFeePerGas *big.Int
	// Nonce is the transaction nonce (optional, will be fetched if not provided).
	Nonce *uint64
}

// SimulateContractOptions contains options for simulating a contract call.
type SimulateContractOptions struct {
	ReadContractOptions
	// Value is the amount of ETH to simulate sending.
	Value *big.Int
}

// EncodeFunctionDataOptions contains options for encoding function data.
type EncodeFunctionDataOptions struct {
	// ABI is the contract ABI as a JSON string or []byte.
	ABI any
	// FunctionName is the name of the function to encode.
	FunctionName string
	// Args are the function arguments.
	Args []any
}

// DecodeFunctionResultOptions contains options for decoding function results.
type DecodeFunctionResultOptions struct {
	// ABI is the contract ABI as a JSON string or []byte.
	ABI any
	// FunctionName is the name of the function that was called.
	FunctionName string
	// Data is the raw result data to decode.
	Data []byte
}

// ReadContract reads from a contract using viem-style API.
// This is equivalent to viem's readContract function.
//
// Example:
//
//	result, err := client.ReadContract(ctx, ReadContractOptions{
//	    Address:      common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
//	    ABI:          `[{"name":"balanceOf","type":"function","inputs":[{"name":"owner","type":"address"}],"outputs":[{"type":"uint256"}]}]`,
//	    FunctionName: "balanceOf",
//	    Args:         []any{common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")},
//	})
func (c *PublicClient) ReadContract(ctx context.Context, opts ReadContractOptions) ([]any, error) {
	// Parse ABI
	parsedABI, err := parseABIInput(opts.ABI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Encode the function call
	calldata, err := parsedABI.EncodeCall(opts.FunctionName, opts.Args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode call for %q: %w", opts.FunctionName, err)
	}

	// Build call request
	callReq := types.CallRequest{
		From: opts.From,
		To:   opts.Address,
		Data: calldata,
	}

	// Execute eth_call
	var result []byte
	if opts.Block != "" {
		result, err = c.Call(ctx, callReq, opts.Block)
	} else {
		result, err = c.Call(ctx, callReq)
	}
	if err != nil {
		return nil, err
	}

	// Decode the return values
	decoded, err := parsedABI.DecodeFunctionResult(opts.FunctionName, result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode result for %q: %w", opts.FunctionName, err)
	}

	return decoded, nil
}

// SimulateContract simulates a contract call without sending a transaction.
// This is useful for checking if a transaction would succeed and getting return values.
func (c *PublicClient) SimulateContract(ctx context.Context, opts SimulateContractOptions) ([]any, error) {
	// Parse ABI
	parsedABI, err := parseABIInput(opts.ABI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Encode the function call
	calldata, err := parsedABI.EncodeCall(opts.FunctionName, opts.Args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode call for %q: %w", opts.FunctionName, err)
	}

	// Build call request
	callReq := types.CallRequest{
		From:  opts.From,
		To:    opts.Address,
		Data:  calldata,
		Value: opts.Value,
	}

	// Execute eth_call
	var result []byte
	if opts.Block != "" {
		result, err = c.Call(ctx, callReq, opts.Block)
	} else {
		result, err = c.Call(ctx, callReq)
	}
	if err != nil {
		return nil, err
	}

	// Decode the return values
	decoded, err := parsedABI.DecodeFunctionResult(opts.FunctionName, result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode result for %q: %w", opts.FunctionName, err)
	}

	return decoded, nil
}

// PrepareContractWrite prepares a transaction for a contract write.
// Returns a Transaction that can be signed and sent via WalletClient.
func (c *PublicClient) PrepareContractWrite(ctx context.Context, opts PrepareContractWriteOptions) (*types.Transaction, error) {
	// Parse ABI
	parsedABI, err := parseABIInput(opts.ABI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Encode the function call
	calldata, err := parsedABI.EncodeCall(opts.FunctionName, opts.Args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode call for %q: %w", opts.FunctionName, err)
	}

	// Get nonce if not provided
	var nonce uint64
	if opts.Nonce != nil {
		nonce = *opts.Nonce
	} else {
		nonce, err = c.GetTransactionCount(ctx, opts.From, BlockTagPending)
		if err != nil {
			return nil, fmt.Errorf("failed to get nonce: %w", err)
		}
	}

	// Estimate gas if not provided
	gas := opts.Gas
	if gas == 0 {
		callReq := types.CallRequest{
			From:  &opts.From,
			To:    opts.Address,
			Data:  calldata,
			Value: opts.Value,
		}
		gas, err = c.EstimateGas(ctx, callReq)
		if err != nil {
			return nil, fmt.Errorf("failed to estimate gas: %w", err)
		}
		// Add 20% buffer
		gas = gas * 120 / 100
	}

	// Get gas price if not provided (for legacy transactions)
	gasPrice := opts.GasPrice
	if gasPrice == nil && opts.MaxFeePerGas == nil {
		gasPrice, err = c.GetGasPrice(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get gas price: %w", err)
		}
	}

	// Get chain ID if available
	var chainID *big.Int
	if c.chain != nil {
		chainID = big.NewInt(int64(c.chain.ID))
	}

	tx := &types.Transaction{
		From:                 opts.From,
		To:                   &opts.Address,
		Data:                 calldata,
		Value:                opts.Value,
		Nonce:                &nonce,
		Gas:                  gas,
		GasPrice:             gasPrice,
		MaxFeePerGas:         opts.MaxFeePerGas,
		MaxPriorityFeePerGas: opts.MaxPriorityFeePerGas,
		ChainID:              chainID,
	}

	return tx, nil
}

// EncodeFunctionData encodes function call data for a contract.
// This is useful when you want to build a transaction manually.
//
// Example:
//
//	data, err := EncodeFunctionData(EncodeFunctionDataOptions{
//	    ABI:          `[{"name":"transfer","type":"function","inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}]}]`,
//	    FunctionName: "transfer",
//	    Args:         []any{common.HexToAddress("0x..."), big.NewInt(1000000)},
//	})
func EncodeFunctionData(opts EncodeFunctionDataOptions) ([]byte, error) {
	parsedABI, err := parseABIInput(opts.ABI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	return parsedABI.EncodeCall(opts.FunctionName, opts.Args...)
}

// DecodeFunctionResult decodes the result of a contract call.
func DecodeFunctionResult(opts DecodeFunctionResultOptions) ([]any, error) {
	parsedABI, err := parseABIInput(opts.ABI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	return parsedABI.DecodeFunctionResult(opts.FunctionName, opts.Data)
}

// parseABIInput parses an ABI from various input types.
func parseABIInput(input any) (*abi.ABI, error) {
	switch v := input.(type) {
	case string:
		return abi.Parse([]byte(v))
	case []byte:
		return abi.Parse(v)
	case *abi.ABI:
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported ABI type: %T (expected string, []byte, or *abi.ABI)", input)
	}
}

// Suppress unused import warning for hexutil
var _ = hexutil.Encode
var _ = json.Marshal
