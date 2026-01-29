package contract

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/client"
	"github.com/ethereum/go-ethereum/common"
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

// Write calls a contract method that modifies state.
// Returns the transaction hash.
func (c *Contract) Write(ctx context.Context, opts WriteOptions, method string, args ...any) (common.Hash, error) {
	// Validate method exists
	_, err := c.abi.GetFunction(method)
	if err != nil {
		return common.Hash{}, err
	}

	// Encode the call
	calldata, err := c.abi.EncodeCall(method, args...)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to encode call for %q: %w", method, err)
	}

	// Build transaction
	tx := client.Transaction{
		From:                 opts.From,
		To:                   &c.address,
		Data:                 calldata,
		Value:                opts.Value,
		Gas:                  opts.Gas,
		GasPrice:             opts.GasPrice,
		MaxFeePerGas:         opts.MaxFeePerGas,
		MaxPriorityFeePerGas: opts.MaxPriorityFeePerGas,
		Nonce:                opts.Nonce,
	}

	// Send the transaction
	hash, err := c.client.SendTransaction(ctx, tx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to send transaction for %q: %w", method, err)
	}

	return hash, nil
}

// WriteAndWait calls a contract method and waits for the transaction to be mined.
// Returns the transaction receipt.
func (c *Contract) WriteAndWait(ctx context.Context, opts WriteOptions, method string, args ...any) (*client.Receipt, error) {
	hash, err := c.Write(ctx, opts, method, args...)
	if err != nil {
		return nil, err
	}

	return c.client.WaitForTransaction(ctx, hash)
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
	callReq := client.CallRequest{
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
// Returns the populated Transaction struct.
func (c *Contract) PrepareTransaction(ctx context.Context, opts WriteOptions, method string, args ...any) (*client.Transaction, error) {
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
		nonce, err = c.client.GetTransactionCount(ctx, opts.From, client.BlockPending)
		if err != nil {
			return nil, fmt.Errorf("failed to get nonce: %w", err)
		}
	}

	// Estimate gas if not provided
	gas := opts.Gas
	if gas == 0 {
		callReq := client.CallRequest{
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

	tx := &client.Transaction{
		From:                 opts.From,
		To:                   &c.address,
		Data:                 calldata,
		Value:                opts.Value,
		Nonce:                &nonce,
		Gas:                  gas,
		GasPrice:             gasPrice,
		MaxFeePerGas:         opts.MaxFeePerGas,
		MaxPriorityFeePerGas: opts.MaxPriorityFeePerGas,
		ChainID:              c.client.ChainID(),
	}

	return tx, nil
}

// Calldata returns the encoded calldata for a method call without sending.
// Useful for building transactions manually or for multicall.
func (c *Contract) Calldata(method string, args ...any) ([]byte, error) {
	return c.abi.EncodeCall(method, args...)
}

// Deploy deploys a new contract with the given bytecode and constructor arguments.
// Returns the transaction hash.
func Deploy(ctx context.Context, cl *client.Client, contractABI []byte, bytecode []byte, opts WriteOptions, args ...any) (common.Hash, error) {
	parsedABI, err := parseABIForDeploy(contractABI)
	if err != nil {
		return common.Hash{}, err
	}

	// Encode constructor arguments if any
	var data []byte
	if len(args) > 0 {
		constructorArgs, err := parsedABI.EncodeConstructor(args...)
		if err != nil {
			return common.Hash{}, fmt.Errorf("failed to encode constructor arguments: %w", err)
		}
		data = append(bytecode, constructorArgs...)
	} else {
		data = bytecode
	}

	// Build deployment transaction (To is nil for contract creation)
	tx := client.Transaction{
		From:                 opts.From,
		To:                   nil,
		Data:                 data,
		Value:                opts.Value,
		Gas:                  opts.Gas,
		GasPrice:             opts.GasPrice,
		MaxFeePerGas:         opts.MaxFeePerGas,
		MaxPriorityFeePerGas: opts.MaxPriorityFeePerGas,
		Nonce:                opts.Nonce,
	}

	return cl.SendTransaction(ctx, tx)
}

// DeployAndWait deploys a contract and waits for it to be mined.
// Returns the contract address and receipt.
func DeployAndWait(ctx context.Context, cl *client.Client, contractABI []byte, bytecode []byte, opts WriteOptions, args ...any) (common.Address, *client.Receipt, error) {
	hash, err := Deploy(ctx, cl, contractABI, bytecode, opts, args...)
	if err != nil {
		return common.Address{}, nil, err
	}

	receipt, err := cl.WaitForTransaction(ctx, hash)
	if err != nil {
		return common.Address{}, nil, err
	}

	if receipt.ContractAddress == nil {
		return common.Address{}, receipt, fmt.Errorf("contract deployment failed: no contract address in receipt")
	}

	return *receipt.ContractAddress, receipt, nil
}

// parseABIForDeploy is a helper to parse ABI for deployment.
func parseABIForDeploy(contractABI []byte) (*abi.ABI, error) {
	return abi.Parse(contractABI)
}
