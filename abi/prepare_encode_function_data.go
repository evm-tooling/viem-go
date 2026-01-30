package abi

import (
	"fmt"
)

// PreparedFunctionData contains the resolved function for encoding.
type PreparedFunctionData struct {
	// Abi contains just the matched function
	Abi []Function
	// FunctionSelector is the 4-byte function selector
	FunctionSelector [4]byte
	// FunctionName is the resolved function name
	FunctionName string
}

// PrepareEncodeFunctionData resolves a function from the ABI and returns
// the function selector. This is a preparation step that can be used before
// the actual encoding, useful for caching or validation.
//
// Example:
//
//	prepared, err := abi.PrepareEncodeFunctionData("transfer", args)
//	if err != nil {
//	    // handle error
//	}
//	// prepared.FunctionSelector is the 4-byte selector
//	// prepared.Abi contains the resolved function
func (a *ABI) PrepareEncodeFunctionData(functionName string, args ...any) (*PreparedFunctionData, error) {
	// Use GetAbiItem to resolve the function (handles overloading)
	item, err := a.GetAbiItem(functionName, &GetAbiItemOptions{Args: args})
	if err != nil {
		return nil, err
	}

	fn, ok := item.(Function)
	if !ok {
		if fnPtr, ok := item.(*Function); ok {
			fn = *fnPtr
		} else {
			return nil, fmt.Errorf("%q is not a function", functionName)
		}
	}

	return &PreparedFunctionData{
		Abi:              []Function{fn},
		FunctionSelector: fn.Selector,
		FunctionName:     fn.Name,
	}, nil
}

// PrepareEncodeFunctionDataBySelector resolves a function from the ABI by its selector.
func (a *ABI) PrepareEncodeFunctionDataBySelector(selector [4]byte) (*PreparedFunctionData, error) {
	fn, err := a.GetFunctionBySelector(selector)
	if err != nil {
		return nil, err
	}

	return &PreparedFunctionData{
		Abi:              []Function{*fn},
		FunctionSelector: fn.Selector,
		FunctionName:     fn.Name,
	}, nil
}

// EncodeWithPrepared encodes function data using a pre-resolved function.
// This is more efficient when encoding the same function multiple times.
func (a *ABI) EncodeWithPrepared(prepared *PreparedFunctionData, args ...any) ([]byte, error) {
	if prepared == nil || len(prepared.Abi) == 0 {
		return nil, fmt.Errorf("invalid prepared function data")
	}

	return a.EncodeFunctionData(prepared.FunctionName, args...)
}
