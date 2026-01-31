package contract

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/types"
)

// BoundContract is a contract instance that supports typed method calls.
// It wraps a Contract and provides type-safe Call methods.
//
// Example:
//
//	// Define method descriptors
//	var ERC20 = struct {
//	    Name      contract.Fn[string]
//	    Symbol    contract.Fn[string]
//	    Decimals  contract.Fn[uint8]
//	    BalanceOf contract.Fn1[common.Address, *big.Int]
//	    Allowance contract.Fn2[common.Address, common.Address, *big.Int]
//	}{
//	    Name:      contract.Fn[string]{Name: "name"},
//	    Symbol:    contract.Fn[string]{Name: "symbol"},
//	    Decimals:  contract.Fn[uint8]{Name: "decimals"},
//	    BalanceOf: contract.Fn1[common.Address, *big.Int]{Name: "balanceOf"},
//	    Allowance: contract.Fn2[common.Address, common.Address, *big.Int]{Name: "allowance"},
//	}
//
//	// Bind the contract
//	token, err := contract.Bind(tokenAddr, erc20ABI, client)
//
//	// Make type-safe calls
//	name, err := contract.Call(token, ctx, ERC20.Name)
//	balance, err := contract.Call1(token, ctx, ERC20.BalanceOf, ownerAddr)
//	allowance, err := contract.Call2(token, ctx, ERC20.Allowance, owner, spender)
type BoundContract struct {
	// Contract is the underlying contract instance.
	// It is exported to allow access to lower-level methods.
	*Contract
}

// Bind creates a new BoundContract from an address, ABI, and client.
// This is the primary way to create a contract for typed calls.
func Bind(address common.Address, abiJSON []byte, c *client.PublicClient) (*BoundContract, error) {
	cont, err := NewContract(address, abiJSON, c)
	if err != nil {
		return nil, err
	}
	return &BoundContract{Contract: cont}, nil
}

// BindWithABI creates a BoundContract using a pre-parsed ABI string.
func BindWithABI(address common.Address, abiStr string, c *client.PublicClient) (*BoundContract, error) {
	return Bind(address, []byte(abiStr), c)
}

// MustBind creates a BoundContract, panicking on error.
func MustBind(address common.Address, abiJSON []byte, c *client.PublicClient) *BoundContract {
	bc, err := Bind(address, abiJSON, c)
	if err != nil {
		panic(err)
	}
	return bc
}

// =============================================================================
// Typed Call Functions (zero arguments)
// =============================================================================

// Call executes a typed zero-argument function and returns the result.
// The return type is inferred from the Fn type parameter.
//
// Example:
//
//	var Name = contract.Fn[string]{Name: "name"}
//	name, err := contract.Call(token, ctx, Name)
func Call[TReturn any](bc *BoundContract, ctx context.Context, fn Fn[TReturn]) (TReturn, error) {
	return callTyped[TReturn](bc, ctx, fn.Name)
}

// CallWithOptions executes a typed zero-argument function with options.
func CallWithOptions[TReturn any](bc *BoundContract, ctx context.Context, opts ReadOptions, fn Fn[TReturn]) (TReturn, error) {
	return callTypedWithOptions[TReturn](bc, ctx, opts, fn.Name)
}

// =============================================================================
// Typed Call Functions (one argument)
// =============================================================================

// Call1 executes a typed single-argument function and returns the result.
// Both the argument type and return type are enforced at compile time.
//
// Example:
//
//	var BalanceOf = contract.Fn1[common.Address, *big.Int]{Name: "balanceOf"}
//	balance, err := contract.Call1(token, ctx, BalanceOf, ownerAddr)
func Call1[TArg, TReturn any](bc *BoundContract, ctx context.Context, fn Fn1[TArg, TReturn], arg TArg) (TReturn, error) {
	return callTyped[TReturn](bc, ctx, fn.Name, arg)
}

// Call1WithOptions executes a typed single-argument function with options.
func Call1WithOptions[TArg, TReturn any](bc *BoundContract, ctx context.Context, opts ReadOptions, fn Fn1[TArg, TReturn], arg TArg) (TReturn, error) {
	return callTypedWithOptions[TReturn](bc, ctx, opts, fn.Name, arg)
}

// =============================================================================
// Typed Call Functions (two arguments)
// =============================================================================

// Call2 executes a typed two-argument function and returns the result.
//
// Example:
//
//	var Allowance = contract.Fn2[common.Address, common.Address, *big.Int]{Name: "allowance"}
//	allowance, err := contract.Call2(token, ctx, Allowance, owner, spender)
func Call2[TArg1, TArg2, TReturn any](bc *BoundContract, ctx context.Context, fn Fn2[TArg1, TArg2, TReturn], arg1 TArg1, arg2 TArg2) (TReturn, error) {
	return callTyped[TReturn](bc, ctx, fn.Name, arg1, arg2)
}

// Call2WithOptions executes a typed two-argument function with options.
func Call2WithOptions[TArg1, TArg2, TReturn any](bc *BoundContract, ctx context.Context, opts ReadOptions, fn Fn2[TArg1, TArg2, TReturn], arg1 TArg1, arg2 TArg2) (TReturn, error) {
	return callTypedWithOptions[TReturn](bc, ctx, opts, fn.Name, arg1, arg2)
}

// =============================================================================
// Typed Call Functions (three arguments)
// =============================================================================

// Call3 executes a typed three-argument function and returns the result.
func Call3[TArg1, TArg2, TArg3, TReturn any](bc *BoundContract, ctx context.Context, fn Fn3[TArg1, TArg2, TArg3, TReturn], arg1 TArg1, arg2 TArg2, arg3 TArg3) (TReturn, error) {
	return callTyped[TReturn](bc, ctx, fn.Name, arg1, arg2, arg3)
}

// Call3WithOptions executes a typed three-argument function with options.
func Call3WithOptions[TArg1, TArg2, TArg3, TReturn any](bc *BoundContract, ctx context.Context, opts ReadOptions, fn Fn3[TArg1, TArg2, TArg3, TReturn], arg1 TArg1, arg2 TArg2, arg3 TArg3) (TReturn, error) {
	return callTypedWithOptions[TReturn](bc, ctx, opts, fn.Name, arg1, arg2, arg3)
}

// =============================================================================
// Typed Call Functions (four arguments)
// =============================================================================

// Call4 executes a typed four-argument function and returns the result.
func Call4[TArg1, TArg2, TArg3, TArg4, TReturn any](bc *BoundContract, ctx context.Context, fn Fn4[TArg1, TArg2, TArg3, TArg4, TReturn], arg1 TArg1, arg2 TArg2, arg3 TArg3, arg4 TArg4) (TReturn, error) {
	return callTyped[TReturn](bc, ctx, fn.Name, arg1, arg2, arg3, arg4)
}

// Call4WithOptions executes a typed four-argument function with options.
func Call4WithOptions[TArg1, TArg2, TArg3, TArg4, TReturn any](bc *BoundContract, ctx context.Context, opts ReadOptions, fn Fn4[TArg1, TArg2, TArg3, TArg4, TReturn], arg1 TArg1, arg2 TArg2, arg3 TArg3, arg4 TArg4) (TReturn, error) {
	return callTypedWithOptions[TReturn](bc, ctx, opts, fn.Name, arg1, arg2, arg3, arg4)
}

// =============================================================================
// Internal Implementation
// =============================================================================

// callTyped is the internal implementation for typed calls.
func callTyped[TReturn any](bc *BoundContract, ctx context.Context, methodName string, args ...any) (TReturn, error) {
	var zero TReturn

	result, err := bc.Read(ctx, methodName, args...)
	if err != nil {
		return zero, err
	}

	if len(result) == 0 {
		return zero, fmt.Errorf("method %q returned no values", methodName)
	}

	// Try direct type assertion
	if typed, ok := result[0].(TReturn); ok {
		return typed, nil
	}

	// Try type conversion
	converted, err := convertResult[TReturn](result[0], methodName)
	if err != nil {
		return zero, err
	}

	return converted, nil
}

// callTypedWithOptions is the internal implementation for typed calls with options.
func callTypedWithOptions[TReturn any](bc *BoundContract, ctx context.Context, opts ReadOptions, methodName string, args ...any) (TReturn, error) {
	var zero TReturn

	result, err := bc.ReadWithOptions(ctx, opts, methodName, args...)
	if err != nil {
		return zero, err
	}

	if len(result) == 0 {
		return zero, fmt.Errorf("method %q returned no values", methodName)
	}

	// Try direct type assertion
	if typed, ok := result[0].(TReturn); ok {
		return typed, nil
	}

	// Try type conversion
	converted, err := convertResult[TReturn](result[0], methodName)
	if err != nil {
		return zero, err
	}

	return converted, nil
}

// =============================================================================
// Write Method Helpers
// =============================================================================

// PrepareWrite prepares a transaction for a zero-argument write method.
func PrepareWrite(bc *BoundContract, ctx context.Context, opts WriteOptions, fn FnWrite) (*Transaction, error) {
	return bc.PrepareTransaction(ctx, opts, fn.Name)
}

// PrepareWrite1 prepares a transaction for a single-argument write method.
func PrepareWrite1[TArg any](bc *BoundContract, ctx context.Context, opts WriteOptions, fn FnWrite1[TArg], arg TArg) (*Transaction, error) {
	return bc.PrepareTransaction(ctx, opts, fn.Name, arg)
}

// PrepareWrite2 prepares a transaction for a two-argument write method.
func PrepareWrite2[TArg1, TArg2 any](bc *BoundContract, ctx context.Context, opts WriteOptions, fn FnWrite2[TArg1, TArg2], arg1 TArg1, arg2 TArg2) (*Transaction, error) {
	return bc.PrepareTransaction(ctx, opts, fn.Name, arg1, arg2)
}

// PrepareWrite3 prepares a transaction for a three-argument write method.
func PrepareWrite3[TArg1, TArg2, TArg3 any](bc *BoundContract, ctx context.Context, opts WriteOptions, fn FnWrite3[TArg1, TArg2, TArg3], arg1 TArg1, arg2 TArg2, arg3 TArg3) (*Transaction, error) {
	return bc.PrepareTransaction(ctx, opts, fn.Name, arg1, arg2, arg3)
}

// Transaction is an alias for types.Transaction for convenience.
type Transaction = types.Transaction
