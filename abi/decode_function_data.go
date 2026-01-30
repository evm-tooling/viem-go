package abi

import (
	"fmt"
)

// DecodedFunctionData represents decoded function calldata.
type DecodedFunctionData struct {
	FunctionName string
	Args         []any
	Selector     [4]byte
}

// DecodeFunctionData decodes calldata and returns the function name and arguments.
// The data must include the 4-byte function selector.
//
// Example:
//
//	result, err := abi.DecodeFunctionData(calldata)
//	fmt.Println(result.FunctionName, result.Args)
func (a *ABI) DecodeFunctionData(data []byte) (*DecodedFunctionData, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("data too short: expected at least 4 bytes, got %d", len(data))
	}

	var selector [4]byte
	copy(selector[:], data[:4])

	// Find matching method
	for _, m := range a.gethABI.Methods {
		var methodSelector [4]byte
		copy(methodSelector[:], m.ID)
		if methodSelector == selector {
			args, err := m.Inputs.Unpack(data[4:])
			if err != nil {
				return nil, fmt.Errorf("failed to decode args for function %q: %w", m.Name, err)
			}
			return &DecodedFunctionData{
				FunctionName: m.Name,
				Args:         args,
				Selector:     selector,
			}, nil
		}
	}

	return nil, fmt.Errorf("function with selector %x not found on ABI", selector)
}

// DecodeFunctionDataByName decodes function arguments using a known function name.
// The data must include the 4-byte function selector.
func (a *ABI) DecodeFunctionDataByName(functionName string, data []byte) ([]any, error) {
	m, ok := a.gethABI.Methods[functionName]
	if !ok {
		return nil, fmt.Errorf("function %q not found on ABI", functionName)
	}

	if len(data) < 4 {
		return nil, fmt.Errorf("data too short: expected at least 4 bytes, got %d", len(data))
	}

	// Skip the 4-byte selector
	unpacked, err := m.Inputs.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode args for function %q: %w", functionName, err)
	}

	return unpacked, nil
}

// DecodeFunctionArgsFromData decodes function arguments from raw data (without selector).
func (a *ABI) DecodeFunctionArgsFromData(functionName string, data []byte) ([]any, error) {
	m, ok := a.gethABI.Methods[functionName]
	if !ok {
		return nil, fmt.Errorf("function %q not found on ABI", functionName)
	}

	unpacked, err := m.Inputs.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode args for function %q: %w", functionName, err)
	}

	return unpacked, nil
}

// DecodeCalldata is an alias for DecodeFunctionData.
// Deprecated: Use DecodeFunctionData instead.
func (a *ABI) DecodeCalldata(data []byte) (string, []any, error) {
	result, err := a.DecodeFunctionData(data)
	if err != nil {
		return "", nil, err
	}
	return result.FunctionName, result.Args, nil
}

// DecodeArgs is an alias for DecodeFunctionDataByName.
// Deprecated: Use DecodeFunctionDataByName instead.
func (a *ABI) DecodeArgs(method string, data []byte) ([]any, error) {
	return a.DecodeFunctionDataByName(method, data)
}

// DecodeArgsFromData is an alias for DecodeFunctionArgsFromData.
// Deprecated: Use DecodeFunctionArgsFromData instead.
func (a *ABI) DecodeArgsFromData(method string, data []byte) ([]any, error) {
	return a.DecodeFunctionArgsFromData(method, data)
}
