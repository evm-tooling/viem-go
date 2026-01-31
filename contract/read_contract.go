package contract

import (
	"context"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/types"
)

// ReadContractParams contains the parameters for a ReadContract call.
type ReadContractParams struct {
	// Address is the contract address to call.
	Address common.Address
	// ABI is the contract ABI as JSON bytes or a pre-parsed *abi.ABI.
	ABI any // []byte, string, or *abi.ABI
	// FunctionName is the name of the function to call.
	FunctionName string
	// Args are the function arguments.
	Args []any
	// BlockTag is the block to read from (default: latest).
	BlockTag client.BlockTag
	// From is the address to use as the caller (optional).
	From *common.Address
}

// ReadContract calls a contract function and returns the result as the specified type T.
// This is the primary API for making typed contract read calls.
//
// Example:
//
//	balance, err := contract.ReadContract[*big.Int](client, contract.ReadContractParams{
//	    Address:      tokenAddr,
//	    ABI:          erc20ABI,
//	    FunctionName: "balanceOf",
//	    Args:         []any{ownerAddr},
//	})
//
//	name, err := contract.ReadContract[string](client, contract.ReadContractParams{
//	    Address:      tokenAddr,
//	    ABI:          erc20ABI,
//	    FunctionName: "name",
//	})
func ReadContract[T any](c *client.PublicClient, params ReadContractParams) (T, error) {
	var zero T

	// Parse the ABI
	parsedABI, err := parseABIParam(params.ABI)
	if err != nil {
		return zero, fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Encode the call
	calldata, err := parsedABI.EncodeCall(params.FunctionName, params.Args...)
	if err != nil {
		return zero, fmt.Errorf("failed to encode call for %q: %w", params.FunctionName, err)
	}

	// Build call request
	callReq := types.CallRequest{
		From: params.From,
		To:   params.Address,
		Data: calldata,
	}

	// Make the call
	var result []byte
	if params.BlockTag != "" {
		result, err = c.Call(context.Background(), callReq, params.BlockTag)
	} else {
		result, err = c.Call(context.Background(), callReq)
	}
	if err != nil {
		return zero, fmt.Errorf("eth_call failed for %q: %w", params.FunctionName, err)
	}

	// Decode and convert the result to type T
	return decodeResultAs[T](parsedABI, params.FunctionName, result)
}

// ReadContractWithContext is like ReadContract but accepts a context.
func ReadContractWithContext[T any](ctx context.Context, c *client.PublicClient, params ReadContractParams) (T, error) {
	var zero T

	// Parse the ABI
	parsedABI, err := parseABIParam(params.ABI)
	if err != nil {
		return zero, fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Encode the call
	calldata, err := parsedABI.EncodeCall(params.FunctionName, params.Args...)
	if err != nil {
		return zero, fmt.Errorf("failed to encode call for %q: %w", params.FunctionName, err)
	}

	// Build call request
	callReq := types.CallRequest{
		From: params.From,
		To:   params.Address,
		Data: calldata,
	}

	// Make the call
	var result []byte
	if params.BlockTag != "" {
		result, err = c.Call(ctx, callReq, params.BlockTag)
	} else {
		result, err = c.Call(ctx, callReq)
	}
	if err != nil {
		return zero, fmt.Errorf("eth_call failed for %q: %w", params.FunctionName, err)
	}

	// Decode and convert the result to type T
	return decodeResultAs[T](parsedABI, params.FunctionName, result)
}

// parseABIParam parses the ABI parameter which can be []byte, string, or *abi.ABI.
func parseABIParam(abiParam any) (*abi.ABI, error) {
	switch v := abiParam.(type) {
	case *abi.ABI:
		return v, nil
	case []byte:
		return abi.Parse(v)
	case string:
		return abi.Parse([]byte(v))
	default:
		return nil, fmt.Errorf("ABI must be []byte, string, or *abi.ABI, got %T", abiParam)
	}
}

// decodeResultAs decodes the raw result bytes and converts to the target type T.
func decodeResultAs[T any](parsedABI *abi.ABI, functionName string, data []byte) (T, error) {
	var zero T
	targetType := reflect.TypeOf(zero)

	// Check if T is a struct (for multi-return values)
	if targetType != nil && targetType.Kind() == reflect.Struct {
		// Use UnpackIntoInterface for struct types
		result := new(T)
		err := parsedABI.DecodeFunctionResultInto(functionName, data, result)
		if err != nil {
			return zero, fmt.Errorf("failed to decode result into struct for %q: %w", functionName, err)
		}
		return *result, nil
	}

	// For non-struct types, decode and type assert
	decoded, err := parsedABI.DecodeFunctionResult(functionName, data)
	if err != nil {
		return zero, fmt.Errorf("failed to decode result for %q: %w", functionName, err)
	}

	if len(decoded) == 0 {
		return zero, fmt.Errorf("function %q returned no values", functionName)
	}

	// Try direct type assertion first
	if typed, ok := decoded[0].(T); ok {
		return typed, nil
	}

	// Handle common type conversions
	return convertToType[T](decoded[0], functionName)
}

// convertToType attempts to convert a value to the target type T.
func convertToType[T any](value any, functionName string) (T, error) {
	var zero T
	targetType := reflect.TypeOf(zero)

	// Handle pointer to big.Int specifically
	if targetType == reflect.TypeOf((*big.Int)(nil)) {
		switch v := value.(type) {
		case *big.Int:
			return any(v).(T), nil
		case int64:
			return any(big.NewInt(v)).(T), nil
		case uint64:
			return any(new(big.Int).SetUint64(v)).(T), nil
		}
	}

	// Handle uint8 (decimals)
	if targetType != nil && targetType.Kind() == reflect.Uint8 {
		switch v := value.(type) {
		case uint8:
			return any(v).(T), nil
		case *big.Int:
			return any(uint8(v.Uint64())).(T), nil
		}
	}

	// Handle uint64
	if targetType != nil && targetType.Kind() == reflect.Uint64 {
		switch v := value.(type) {
		case uint64:
			return any(v).(T), nil
		case *big.Int:
			return any(v.Uint64()).(T), nil
		}
	}

	// Handle common.Address
	if targetType == reflect.TypeOf(common.Address{}) {
		if addr, ok := value.(common.Address); ok {
			return any(addr).(T), nil
		}
	}

	// Handle bool
	if targetType != nil && targetType.Kind() == reflect.Bool {
		if b, ok := value.(bool); ok {
			return any(b).(T), nil
		}
	}

	// Handle string
	if targetType != nil && targetType.Kind() == reflect.String {
		if s, ok := value.(string); ok {
			return any(s).(T), nil
		}
	}

	// Handle []byte
	if targetType != nil && targetType.Kind() == reflect.Slice && targetType.Elem().Kind() == reflect.Uint8 {
		if b, ok := value.([]byte); ok {
			return any(b).(T), nil
		}
	}

	// Handle [32]byte
	if targetType != nil && targetType.Kind() == reflect.Array && targetType.Len() == 32 && targetType.Elem().Kind() == reflect.Uint8 {
		switch v := value.(type) {
		case [32]byte:
			return any(v).(T), nil
		case []byte:
			if len(v) == 32 {
				var arr [32]byte
				copy(arr[:], v)
				return any(arr).(T), nil
			}
		}
	}

	return zero, fmt.Errorf("function %q returned %T, cannot convert to %T", functionName, value, zero)
}

// MustReadContract is like ReadContract but panics on error.
func MustReadContract[T any](c *client.PublicClient, params ReadContractParams) T {
	result, err := ReadContract[T](c, params)
	if err != nil {
		panic(err)
	}
	return result
}
