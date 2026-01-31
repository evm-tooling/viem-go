package contract

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/types"
)

// WriteOptions contains options for write operations.
type WriteOptions struct {
	// From is the address sending the transaction.
	From common.Address
	// Value is the amount of ETH to send with the transaction.
	Value *big.Int
	// Gas is the gas limit for the transaction (0 = estimate).
	Gas uint64
	// GasPrice is the gas price for legacy transactions.
	GasPrice *big.Int
	// MaxFeePerGas is the max fee per gas for EIP-1559 transactions.
	MaxFeePerGas *big.Int
	// MaxPriorityFeePerGas is the max priority fee for EIP-1559 transactions.
	MaxPriorityFeePerGas *big.Int
	// Nonce is the transaction nonce (nil = auto).
	Nonce *uint64
}

// EstimateGas estimates the gas required for a contract method call.
func (c *Contract) EstimateGas(ctx context.Context, opts WriteOptions, method string, args ...any) (uint64, error) {
	// Validate method exists
	_, err := c.abi.GetFunction(method)
	if err != nil {
		return 0, err
	}

	// Encode the call
	calldata, err := c.abi.EncodeCall(method, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to encode call for %q: %w", method, err)
	}

	// Build call request
	callReq := types.CallRequest{
		From:  &opts.From,
		To:    c.address,
		Data:  calldata,
		Value: opts.Value,
	}

	// Estimate gas
	gas, err := c.client.EstimateGas(ctx, callReq)
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas for %q: %w", method, err)
	}

	return gas, nil
}

// PrepareTransaction prepares a transaction for signing without sending it.
// Returns the populated Transaction struct. Use this with a signer to send transactions.
func (c *Contract) PrepareTransaction(ctx context.Context, opts WriteOptions, method string, args ...any) (*types.Transaction, error) {
	// Validate method exists
	_, err := c.abi.GetFunction(method)
	if err != nil {
		return nil, err
	}

	// Encode the call
	calldata, err := c.abi.EncodeCall(method, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode call for %q: %w", method, err)
	}

	// Get nonce if not provided
	var nonce uint64
	if opts.Nonce != nil {
		nonce = *opts.Nonce
	} else {
		nonce, err = c.client.GetTransactionCount(ctx, opts.From, client.BlockTagPending)
		if err != nil {
			return nil, fmt.Errorf("failed to get nonce: %w", err)
		}
	}

	// Estimate gas if not provided
	gas := opts.Gas
	if gas == 0 {
		callReq := types.CallRequest{
			From:  &opts.From,
			To:    c.address,
			Data:  calldata,
			Value: opts.Value,
		}
		gas, err = c.client.EstimateGas(ctx, callReq)
		if err != nil {
			return nil, fmt.Errorf("failed to estimate gas: %w", err)
		}
		// Add 20% buffer
		gas = gas * 120 / 100
	}

	// Get gas price if not provided (for legacy transactions)
	gasPrice := opts.GasPrice
	if gasPrice == nil && opts.MaxFeePerGas == nil {
		gasPrice, err = c.client.GetGasPrice(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get gas price: %w", err)
		}
	}

	// Get chain ID if available
	var chainID *big.Int
	if c.client.Chain() != nil {
		chainID = big.NewInt(int64(c.client.Chain().ID))
	}

	tx := &types.Transaction{
		From:                 opts.From,
		To:                   &c.address,
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

// Calldata returns the encoded calldata for a method call without sending.
// Useful for building transactions manually or for multicall.
func (c *Contract) Calldata(method string, args ...any) ([]byte, error) {
	return c.abi.EncodeCall(method, args...)
}

// Deploy deploys a new contract with the given bytecode and constructor arguments.
// Returns the prepared transaction. Use a WalletClient to sign and send.
func PrepareDeployTransaction(ctx context.Context, cl *client.PublicClient, contractABI []byte, bytecode []byte, opts WriteOptions, args ...any) (*types.Transaction, error) {
	parsedABI, err := parseABIForDeploy(contractABI)
	if err != nil {
		return nil, err
	}

	// Encode constructor arguments if any
	var data []byte
	if len(args) > 0 {
		constructorArgs, encodeErr := parsedABI.EncodeConstructor(args...)
		if encodeErr != nil {
			return nil, fmt.Errorf("failed to encode constructor arguments: %w", encodeErr)
		}
		data = append(bytecode, constructorArgs...)
	} else {
		data = bytecode
	}

	// Get nonce if not provided
	var nonce uint64
	if opts.Nonce != nil {
		nonce = *opts.Nonce
	} else {
		nonce, err = cl.GetTransactionCount(ctx, opts.From, client.BlockTagPending)
		if err != nil {
			return nil, fmt.Errorf("failed to get nonce: %w", err)
		}
	}

	// Get gas price if not provided
	gasPrice := opts.GasPrice
	if gasPrice == nil && opts.MaxFeePerGas == nil {
		gasPrice, err = cl.GetGasPrice(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get gas price: %w", err)
		}
	}

	// Get chain ID if available
	var chainID *big.Int
	if cl.Chain() != nil {
		chainID = big.NewInt(int64(cl.Chain().ID))
	}

	tx := &types.Transaction{
		From:                 opts.From,
		To:                   nil, // nil for contract creation
		Data:                 data,
		Value:                opts.Value,
		Nonce:                &nonce,
		Gas:                  opts.Gas,
		GasPrice:             gasPrice,
		MaxFeePerGas:         opts.MaxFeePerGas,
		MaxPriorityFeePerGas: opts.MaxPriorityFeePerGas,
		ChainID:              chainID,
	}

	return tx, nil
}

// parseABIForDeploy is a helper to parse ABI for deployment.
func parseABIForDeploy(contractABI []byte) (*abi.ABI, error) {
	return abi.Parse(contractABI)
}
