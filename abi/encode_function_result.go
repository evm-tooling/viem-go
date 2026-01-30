package abi

import (
	"fmt"
)

// EncodeFunctionResult encodes return values for a function.
// This is the inverse of DecodeFunctionResult.
//
// Example:
//
//	encoded, err := abi.EncodeFunctionResult("balanceOf", balance)
func (a *ABI) EncodeFunctionResult(functionName string, values ...any) ([]byte, error) {
	m, ok := a.gethABI.Methods[functionName]
	if !ok {
		return nil, fmt.Errorf("function %q not found on ABI", functionName)
	}

	if len(m.Outputs) == 0 {
		if len(values) > 0 {
			return nil, fmt.Errorf("function %q has no outputs but %d values provided", functionName, len(values))
		}
		return []byte{}, nil
	}

	if len(values) != len(m.Outputs) {
		return nil, fmt.Errorf("function %q expects %d outputs but %d values provided", functionName, len(m.Outputs), len(values))
	}

	packed, err := m.Outputs.Pack(values...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function result for %q: %w", functionName, err)
	}

	return packed, nil
}
