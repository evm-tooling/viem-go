package abi

import (
	"fmt"
)

// EncodeFunctionData encodes a function call with the given function name and arguments.
// Returns the full calldata including the 4-byte function selector.
//
// Example:
//
//	calldata, err := abi.EncodeFunctionData("transfer", to, amount)
func (a *ABI) EncodeFunctionData(functionName string, args ...any) ([]byte, error) {
	m, ok := a.gethABI.Methods[functionName]
	if !ok {
		return nil, fmt.Errorf("function %q not found on ABI", functionName)
	}

	packed, err := a.gethABI.Pack(functionName, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function data for %q: %w", functionName, err)
	}

	// gethABI.Pack already includes the selector
	_ = m // method found
	return packed, nil
}

// EncodeFunctionDataWithSelector encodes function arguments and prepends the given selector.
// Useful when you have a custom selector or are encoding for a different function.
func (a *ABI) EncodeFunctionDataWithSelector(selector [4]byte, functionName string, args ...any) ([]byte, error) {
	m, ok := a.gethABI.Methods[functionName]
	if !ok {
		return nil, fmt.Errorf("function %q not found on ABI", functionName)
	}

	// Pack arguments only (without selector)
	packed, err := m.Inputs.Pack(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode arguments for function %q: %w", functionName, err)
	}

	// Prepend custom selector
	result := make([]byte, 4+len(packed))
	copy(result[:4], selector[:])
	copy(result[4:], packed)
	return result, nil
}

// EncodeFunctionArgs encodes only the function arguments without the selector.
// Useful for encoding constructor arguments or when you need raw argument encoding.
func (a *ABI) EncodeFunctionArgs(functionName string, args ...any) ([]byte, error) {
	m, ok := a.gethABI.Methods[functionName]
	if !ok {
		return nil, fmt.Errorf("function %q not found on ABI", functionName)
	}

	packed, err := m.Inputs.Pack(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode arguments for function %q: %w", functionName, err)
	}

	return packed, nil
}

// EncodeConstructor encodes constructor arguments.
func (a *ABI) EncodeConstructor(args ...any) ([]byte, error) {
	if a.gethABI.Constructor.Inputs == nil {
		if len(args) > 0 {
			return nil, fmt.Errorf("constructor takes no arguments but %d provided", len(args))
		}
		return nil, nil
	}

	packed, err := a.gethABI.Constructor.Inputs.Pack(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode constructor arguments: %w", err)
	}

	return packed, nil
}

// Pack is an alias for EncodeFunctionData for compatibility with go-ethereum naming.
func (a *ABI) Pack(functionName string, args ...any) ([]byte, error) {
	return a.EncodeFunctionData(functionName, args...)
}

// EncodeCall is an alias for EncodeFunctionData.
// Deprecated: Use EncodeFunctionData instead.
func (a *ABI) EncodeCall(method string, args ...any) ([]byte, error) {
	return a.EncodeFunctionData(method, args...)
}

// EncodeArgs is an alias for EncodeFunctionArgs.
// Deprecated: Use EncodeFunctionArgs instead.
func (a *ABI) EncodeArgs(method string, args ...any) ([]byte, error) {
	return a.EncodeFunctionArgs(method, args...)
}

// EncodeCallWithSelector is an alias for EncodeFunctionDataWithSelector.
// Deprecated: Use EncodeFunctionDataWithSelector instead.
func (a *ABI) EncodeCallWithSelector(selector [4]byte, method string, args ...any) ([]byte, error) {
	return a.EncodeFunctionDataWithSelector(selector, method, args...)
}
