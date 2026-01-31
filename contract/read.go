package contract

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/types"
)

// ReadOptions contains options for read operations.
type ReadOptions struct {
	// From is the address to use as the caller (optional).
	From *common.Address
	// Block is the block tag to read from (default: latest).
	Block client.BlockTag
}

// Read calls a contract method and returns the decoded return values.
// This performs an eth_call and does not modify state.
func (c *Contract) Read(ctx context.Context, method string, args ...any) ([]any, error) {
	return c.ReadWithOptions(ctx, ReadOptions{}, method, args...)
}

// ReadWithOptions calls a contract method with custom options.
func (c *Contract) ReadWithOptions(ctx context.Context, opts ReadOptions, method string, args ...any) ([]any, error) {
	// Validate method exists
	fn, err := c.abi.GetFunction(method)
	if err != nil {
		return nil, err
	}

	// Check if method is read-only (pure/view)
	// Note: We allow calling non-view functions via eth_call for simulation
	_ = fn.IsReadOnly() // Currently we don't enforce read-only, eth_call works for simulation

	// Encode the call
	calldata, err := c.abi.EncodeCall(method, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode call for %q: %w", method, err)
	}

	// Build call request
	callReq := types.CallRequest{
		From: opts.From,
		To:   c.address,
		Data: calldata,
	}

	// Make the call
	var result []byte
	if opts.Block != "" {
		result, err = c.client.Call(ctx, callReq, opts.Block)
	} else {
		result, err = c.client.Call(ctx, callReq)
	}
	if err != nil {
		return nil, fmt.Errorf("eth_call failed for %q: %w", method, err)
	}

	// Decode the return value
	decoded, err := c.abi.DecodeReturn(method, result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode return for %q: %w", method, err)
	}

	return decoded, nil
}

// ReadBigInt calls a method and returns the result as *big.Int.
// Useful for methods that return uint256 or int256.
func (c *Contract) ReadBigInt(ctx context.Context, method string, args ...any) (*big.Int, error) {
	result, err := c.Read(ctx, method, args...)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("method %q returned no values", method)
	}

	switch v := result[0].(type) {
	case *big.Int:
		return v, nil
	case int64:
		return big.NewInt(v), nil
	case uint64:
		return new(big.Int).SetUint64(v), nil
	default:
		return nil, fmt.Errorf("method %q returned %T, expected *big.Int", method, result[0])
	}
}

// ReadAddress calls a method and returns the result as common.Address.
func (c *Contract) ReadAddress(ctx context.Context, method string, args ...any) (common.Address, error) {
	result, err := c.Read(ctx, method, args...)
	if err != nil {
		return common.Address{}, err
	}

	if len(result) == 0 {
		return common.Address{}, fmt.Errorf("method %q returned no values", method)
	}

	addr, ok := result[0].(common.Address)
	if !ok {
		return common.Address{}, fmt.Errorf("method %q returned %T, expected common.Address", method, result[0])
	}

	return addr, nil
}

// ReadBool calls a method and returns the result as bool.
func (c *Contract) ReadBool(ctx context.Context, method string, args ...any) (bool, error) {
	result, err := c.Read(ctx, method, args...)
	if err != nil {
		return false, err
	}

	if len(result) == 0 {
		return false, fmt.Errorf("method %q returned no values", method)
	}

	b, ok := result[0].(bool)
	if !ok {
		return false, fmt.Errorf("method %q returned %T, expected bool", method, result[0])
	}

	return b, nil
}

// ReadString calls a method and returns the result as string.
func (c *Contract) ReadString(ctx context.Context, method string, args ...any) (string, error) {
	result, err := c.Read(ctx, method, args...)
	if err != nil {
		return "", err
	}

	if len(result) == 0 {
		return "", fmt.Errorf("method %q returned no values", method)
	}

	s, ok := result[0].(string)
	if !ok {
		return "", fmt.Errorf("method %q returned %T, expected string", method, result[0])
	}

	return s, nil
}

// ReadBytes calls a method and returns the result as []byte.
func (c *Contract) ReadBytes(ctx context.Context, method string, args ...any) ([]byte, error) {
	result, err := c.Read(ctx, method, args...)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("method %q returned no values", method)
	}

	b, ok := result[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("method %q returned %T, expected []byte", method, result[0])
	}

	return b, nil
}

// ReadBytes32 calls a method and returns the result as [32]byte.
func (c *Contract) ReadBytes32(ctx context.Context, method string, args ...any) ([32]byte, error) {
	result, err := c.Read(ctx, method, args...)
	if err != nil {
		return [32]byte{}, err
	}

	if len(result) == 0 {
		return [32]byte{}, fmt.Errorf("method %q returned no values", method)
	}

	switch v := result[0].(type) {
	case [32]byte:
		return v, nil
	case []byte:
		if len(v) != 32 {
			return [32]byte{}, fmt.Errorf("method %q returned %d bytes, expected 32", method, len(v))
		}
		var arr [32]byte
		copy(arr[:], v)
		return arr, nil
	default:
		return [32]byte{}, fmt.Errorf("method %q returned %T, expected [32]byte", method, result[0])
	}
}

// ReadUint8 calls a method and returns the result as uint8.
func (c *Contract) ReadUint8(ctx context.Context, method string, args ...any) (uint8, error) {
	result, err := c.Read(ctx, method, args...)
	if err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, fmt.Errorf("method %q returned no values", method)
	}

	switch v := result[0].(type) {
	case uint8:
		return v, nil
	case *big.Int:
		return uint8(v.Uint64()), nil
	default:
		return 0, fmt.Errorf("method %q returned %T, expected uint8", method, result[0])
	}
}

// ReadUint64 calls a method and returns the result as uint64.
func (c *Contract) ReadUint64(ctx context.Context, method string, args ...any) (uint64, error) {
	result, err := c.Read(ctx, method, args...)
	if err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, fmt.Errorf("method %q returned no values", method)
	}

	switch v := result[0].(type) {
	case uint64:
		return v, nil
	case *big.Int:
		return v.Uint64(), nil
	default:
		return 0, fmt.Errorf("method %q returned %T, expected uint64", method, result[0])
	}
}

// ReadInto decodes the return value into the provided struct or variables.
func (c *Contract) ReadInto(ctx context.Context, output any, method string, args ...any) error {
	// Encode the call
	calldata, err := c.abi.EncodeCall(method, args...)
	if err != nil {
		return fmt.Errorf("failed to encode call for %q: %w", method, err)
	}

	// Build call request
	callReq := types.CallRequest{
		To:   c.address,
		Data: calldata,
	}

	// Make the call
	result, err := c.client.Call(ctx, callReq)
	if err != nil {
		return fmt.Errorf("eth_call failed for %q: %w", method, err)
	}

	// Decode into the provided output
	return c.abi.DecodeReturnInto(method, result, output)
}

// Simulate simulates a transaction without actually executing it.
// Returns the return value of the function call.
func (c *Contract) Simulate(ctx context.Context, opts ReadOptions, method string, args ...any) ([]any, error) {
	return c.ReadWithOptions(ctx, opts, method, args...)
}
