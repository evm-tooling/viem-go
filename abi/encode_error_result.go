package abi

import (
	"fmt"
)

// EncodeErrorResult encodes an error with the given name and arguments.
// Returns the encoded error data including the 4-byte error selector.
//
// Example:
//
//	encoded, err := abi.EncodeErrorResult("InsufficientBalance", balance, required)
func (a *ABI) EncodeErrorResult(errorName string, args ...any) ([]byte, error) {
	e, ok := a.gethABI.Errors[errorName]
	if !ok {
		return nil, fmt.Errorf("error %q not found on ABI", errorName)
	}

	// Check if args were provided but error has no inputs
	if len(args) > 0 && len(e.Inputs) == 0 {
		return nil, fmt.Errorf("error %q takes no arguments but %d provided", errorName, len(args))
	}

	// For errors with no inputs, just return the selector
	if len(e.Inputs) == 0 {
		result := make([]byte, 4)
		copy(result[:4], e.ID[:4])
		return result, nil
	}

	packed, err := e.Inputs.Pack(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode error %q: %w", errorName, err)
	}

	// Prepend error selector
	result := make([]byte, 4+len(packed))
	copy(result[:4], e.ID[:4])
	copy(result[4:], packed)
	return result, nil
}

// EncodeError is an alias for EncodeErrorResult.
// Deprecated: Use EncodeErrorResult instead.
func (a *ABI) EncodeError(name string, args ...any) ([]byte, error) {
	return a.EncodeErrorResult(name, args...)
}
