package contract

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/types"
)

// =============================================================================
// Typed Function Descriptors (Fn, Fn1, Fn2, Fn3, Fn4)
// =============================================================================

// Fn represents a zero-argument function returning TReturn.
// Use this to define typed method descriptors for contract calls.
//
// Example:
//
//	var Name = contract.Fn[string]{Name: "name"}
//	name, err := contract.Call(bc, ctx, Name)
type Fn[TReturn any] struct {
	Name string
}

// Fn1 represents a single-argument function.
// The type parameters specify the argument type and return type.
//
// Example:
//
//	var BalanceOf = contract.Fn1[common.Address, *big.Int]{Name: "balanceOf"}
//	balance, err := contract.Call1(bc, ctx, BalanceOf, ownerAddr)
type Fn1[TArg, TReturn any] struct {
	Name string
}

// Fn2 represents a two-argument function.
//
// Example:
//
//	var Allowance = contract.Fn2[common.Address, common.Address, *big.Int]{Name: "allowance"}
//	allowance, err := contract.Call2(bc, ctx, Allowance, owner, spender)
type Fn2[TArg1, TArg2, TReturn any] struct {
	Name string
}

// Fn3 represents a three-argument function.
type Fn3[TArg1, TArg2, TArg3, TReturn any] struct {
	Name string
}

// Fn4 represents a four-argument function.
type Fn4[TArg1, TArg2, TArg3, TArg4, TReturn any] struct {
	Name string
}

// FnWrite represents a write function (state-changing).
type FnWrite struct {
	Name string
}

// FnWrite1 represents a single-argument write function.
type FnWrite1[TArg any] struct {
	Name string
}

// FnWrite2 represents a two-argument write function.
type FnWrite2[TArg1, TArg2 any] struct {
	Name string
}

// FnWrite3 represents a three-argument write function.
type FnWrite3[TArg1, TArg2, TArg3 any] struct {
	Name string
}

// =============================================================================
// Legacy Types (kept for backwards compatibility)
// =============================================================================

// ReadMethod is a typed method descriptor for read-only contract calls.
// The type parameter TReturn specifies the return type of the method.
// Deprecated: Use Fn, Fn1, Fn2 etc. for better type safety on arguments.
type ReadMethod[TReturn any] struct {
	Name string
}

// WriteMethod is a method descriptor for state-changing contract calls.
// Deprecated: Use FnWrite, FnWrite1, FnWrite2 etc. for better type safety.
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
		// Try type conversion
		converted, convErr := convertResult[TReturn](result[0], method.Name)
		if convErr != nil {
			return zero, fmt.Errorf("method %q returned %T, expected %T", method.Name, result[0], zero)
		}
		return converted, nil
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
		converted, convErr := convertResult[TReturn](result[0], method.Name)
		if convErr != nil {
			return zero, fmt.Errorf("method %q returned %T, expected %T", method.Name, result[0], zero)
		}
		return converted, nil
	}

	return typed, nil
}

// PrepareWriteTyped prepares a transaction for a typed write method.
// Returns the transaction for signing. Use with a WalletClient to send.
func PrepareWriteTyped(c *Contract, ctx context.Context, opts WriteOptions, method WriteMethod, args ...any) (*types.Transaction, error) {
	return c.PrepareTransaction(ctx, opts, method.Name, args...)
}

// =============================================================================
// TypedContract
// =============================================================================

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
func NewTypedContract[T any](address common.Address, abiJSON []byte, c *client.PublicClient, methods T) (*TypedContract[T], error) {
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

// =============================================================================
// Convenience Type Aliases
// =============================================================================

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

// =============================================================================
// Helper Functions
// =============================================================================

// convertResult attempts to convert a value to the target type TReturn.
func convertResult[TReturn any](value any, methodName string) (TReturn, error) {
	var zero TReturn

	// Try direct conversion for common types
	switch v := value.(type) {
	case *big.Int:
		if result, ok := any(v).(TReturn); ok {
			return result, nil
		}
		// Handle uint8, uint64 conversions from *big.Int
		switch any(zero).(type) {
		case uint8:
			return any(uint8(v.Uint64())).(TReturn), nil
		case uint64:
			return any(v.Uint64()).(TReturn), nil
		}
	case int64:
		if _, ok := any(zero).(*big.Int); ok {
			return any(big.NewInt(v)).(TReturn), nil
		}
	case uint64:
		if _, ok := any(zero).(*big.Int); ok {
			return any(new(big.Int).SetUint64(v)).(TReturn), nil
		}
	}

	return zero, fmt.Errorf("cannot convert %T to %T for method %q", value, zero, methodName)
}
