package contract

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ChefBingbong/viem-go/client"
	"github.com/ethereum/go-ethereum/common"
)

// ReadMethod is a typed method descriptor for read-only contract calls.
// The type parameter TReturn specifies the return type of the method.
type ReadMethod[TReturn any] struct {
	Name string
}

// WriteMethod is a method descriptor for state-changing contract calls.
type WriteMethod struct {
	Name string
}

// ReadTyped calls a contract method with a typed return value.
// The return type is inferred from the ReadMethod type parameter.
//
// Example:
//
//	var BalanceOf = contract.ReadMethod[*big.Int]{Name: "balanceOf"}
//	balance, err := contract.ReadTyped(c, ctx, BalanceOf, ownerAddress)
func ReadTyped[TReturn any](c *Contract, ctx context.Context, method ReadMethod[TReturn], args ...any) (TReturn, error) {
	var zero TReturn

	result, err := c.Read(ctx, method.Name, args...)
	if err != nil {
		return zero, err
	}

	if len(result) == 0 {
		return zero, fmt.Errorf("method %q returned no values", method.Name)
	}

	// Type assert the result
	typed, ok := result[0].(TReturn)
	if !ok {
		return zero, fmt.Errorf("method %q returned %T, expected %T", method.Name, result[0], zero)
	}

	return typed, nil
}

// ReadTypedWithOptions calls a contract method with options and a typed return value.
func ReadTypedWithOptions[TReturn any](c *Contract, ctx context.Context, opts ReadOptions, method ReadMethod[TReturn], args ...any) (TReturn, error) {
	var zero TReturn

	result, err := c.ReadWithOptions(ctx, opts, method.Name, args...)
	if err != nil {
		return zero, err
	}

	if len(result) == 0 {
		return zero, fmt.Errorf("method %q returned no values", method.Name)
	}

	typed, ok := result[0].(TReturn)
	if !ok {
		return zero, fmt.Errorf("method %q returned %T, expected %T", method.Name, result[0], zero)
	}

	return typed, nil
}

// WriteTyped sends a transaction to a contract method.
// Returns the transaction hash.
//
// Example:
//
//	var Transfer = contract.WriteMethod{Name: "transfer"}
//	txHash, err := contract.WriteTyped(c, ctx, opts, Transfer, to, amount)
func WriteTyped(c *Contract, ctx context.Context, opts WriteOptions, method WriteMethod, args ...any) (common.Hash, error) {
	return c.Write(ctx, opts, method.Name, args...)
}

// WriteTypedAndWait sends a transaction and waits for it to be mined.
func WriteTypedAndWait(c *Contract, ctx context.Context, opts WriteOptions, method WriteMethod, args ...any) (*client.Receipt, error) {
	return c.WriteAndWait(ctx, opts, method.Name, args...)
}

// TypedContract is a generic wrapper that embeds a Contract and a method template.
// The template type T should be a struct containing ReadMethod and WriteMethod fields.
//
// Example:
//
//	type ERC20Methods struct {
//	    Name       contract.ReadMethod[string]
//	    Symbol     contract.ReadMethod[string]
//	    Decimals   contract.ReadMethod[uint8]
//	    BalanceOf  contract.ReadMethod[*big.Int]
//	    Transfer   contract.WriteMethod
//	}
//
//	var ERC20 = ERC20Methods{
//	    Name:      contract.ReadMethod[string]{Name: "name"},
//	    Symbol:    contract.ReadMethod[string]{Name: "symbol"},
//	    Decimals:  contract.ReadMethod[uint8]{Name: "decimals"},
//	    BalanceOf: contract.ReadMethod[*big.Int]{Name: "balanceOf"},
//	    Transfer:  contract.WriteMethod{Name: "transfer"},
//	}
//
//	token := contract.NewTypedContract(address, abiJSON, client, ERC20)
//	balance, err := token.Read(ctx, token.Methods.BalanceOf, owner)
type TypedContract[T any] struct {
	*Contract
	Methods T
}

// NewTypedContract creates a new TypedContract with the given methods template.
func NewTypedContract[T any](address common.Address, abiJSON []byte, c *client.Client, methods T) (*TypedContract[T], error) {
	cont, err := NewContract(address, abiJSON, c)
	if err != nil {
		return nil, err
	}

	return &TypedContract[T]{
		Contract: cont,
		Methods:  methods,
	}, nil
}

// Read calls a typed read method on this contract.
func (tc *TypedContract[T]) Read(ctx context.Context, method ReadMethod[any], args ...any) (any, error) {
	return ReadTyped(tc.Contract, ctx, method, args...)
}

// Convenience type aliases for common return types
type (
	// ReadBigInt is a method that returns *big.Int
	ReadBigInt = ReadMethod[*big.Int]
	// ReadAddress is a method that returns common.Address
	ReadAddress = ReadMethod[common.Address]
	// ReadBool is a method that returns bool
	ReadBool = ReadMethod[bool]
	// ReadString is a method that returns string
	ReadString = ReadMethod[string]
	// ReadBytes is a method that returns []byte
	ReadBytes = ReadMethod[[]byte]
	// ReadUint8 is a method that returns uint8
	ReadUint8 = ReadMethod[uint8]
	// ReadUint64 is a method that returns uint64
	ReadUint64 = ReadMethod[uint64]
)
