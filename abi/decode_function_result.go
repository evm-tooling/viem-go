package abi

import (
	"fmt"
)

// DecodeFunctionResult decodes the return data from a function call.
// Returns a slice of decoded values matching the function's output parameters.
//
// Example:
//
//	result, err := abi.DecodeFunctionResult("balanceOf", returnData)
func (a *ABI) DecodeFunctionResult(functionName string, data []byte) ([]any, error) {
	m, ok := a.gethABI.Methods[functionName]
	if !ok {
		return nil, fmt.Errorf("function %q not found on ABI", functionName)
	}

	if len(data) == 0 {
		if len(m.Outputs) > 0 {
			return nil, fmt.Errorf("expected return data for function %q but got empty", functionName)
		}
		return nil, nil
	}

	unpacked, err := m.Outputs.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode function result for %q: %w", functionName, err)
	}

	return unpacked, nil
}

// DecodeFunctionResultInto decodes the return data into the provided struct or variables.
// The output parameter must be a pointer to a struct or a slice of pointers.
func (a *ABI) DecodeFunctionResultInto(functionName string, data []byte, output any) error {
	m, ok := a.gethABI.Methods[functionName]
	if !ok {
		return fmt.Errorf("function %q not found on ABI", functionName)
	}

	if len(data) == 0 {
		if len(m.Outputs) > 0 {
			return fmt.Errorf("expected return data for function %q but got empty", functionName)
		}
		return nil
	}

	return a.gethABI.UnpackIntoInterface(output, functionName, data)
}

// DecodeReturn is an alias for DecodeFunctionResult.
// Deprecated: Use DecodeFunctionResult instead.
func (a *ABI) DecodeReturn(method string, data []byte) ([]any, error) {
	return a.DecodeFunctionResult(method, data)
}

// Unpack is an alias for DecodeFunctionResult for compatibility with go-ethereum naming.
func (a *ABI) Unpack(method string, data []byte) ([]any, error) {
	return a.DecodeFunctionResult(method, data)
}

// DecodeReturnInto is an alias for DecodeFunctionResultInto.
// Deprecated: Use DecodeFunctionResultInto instead.
func (a *ABI) DecodeReturnInto(method string, data []byte, output any) error {
	return a.DecodeFunctionResultInto(method, data, output)
}

// UnpackIntoInterface is an alias for DecodeFunctionResultInto.
func (a *ABI) UnpackIntoInterface(method string, data []byte, output any) error {
	return a.DecodeFunctionResultInto(method, data, output)
}
