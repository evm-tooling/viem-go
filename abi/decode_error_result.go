package abi

import (
	"fmt"
	"math/big"
)

// DecodedErrorResult represents a decoded error.
type DecodedErrorResult struct {
	ErrorName string
	Args      []any
	Selector  [4]byte
	AbiItem   *Error // nil for standard errors like Error(string) and Panic(uint256)
}

// Standard error selectors
var (
	// Error(string) selector: 0x08c379a0
	errorSelector = [4]byte{0x08, 0xc3, 0x79, 0xa0}
	// Panic(uint256) selector: 0x4e487b71
	panicSelector = [4]byte{0x4e, 0x48, 0x7b, 0x71}
)

// DecodeErrorResult decodes an error return value.
// Handles standard errors (Error(string), Panic(uint256)) and custom errors from the ABI.
//
// Example:
//
//	result, err := abi.DecodeErrorResult(errorData)
//	fmt.Println(result.ErrorName, result.Args)
func (a *ABI) DecodeErrorResult(data []byte) (*DecodedErrorResult, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("cannot decode zero data")
	}

	if len(data) < 4 {
		return nil, fmt.Errorf("error data too short: expected at least 4 bytes, got %d", len(data))
	}

	var selector [4]byte
	copy(selector[:], data[:4])

	// Check for standard Error(string)
	if selector == errorSelector {
		return decodeStandardError(data)
	}

	// Check for Panic(uint256)
	if selector == panicSelector {
		return decodePanicError(data)
	}

	// Try to match against custom errors in ABI
	if a != nil {
		for _, e := range a.gethABI.Errors {
			var errSelector [4]byte
			copy(errSelector[:], e.ID[:4])
			if errSelector == selector {
				unpacked, err := e.Inputs.Unpack(data[4:])
				if err != nil {
					return nil, fmt.Errorf("failed to decode error %q: %w", e.Name, err)
				}

				// Get the error from our map
				abiError, _ := a.GetError(e.Name)

				return &DecodedErrorResult{
					ErrorName: e.Name,
					Args:      unpacked,
					Selector:  selector,
					AbiItem:   abiError,
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("unknown error selector: 0x%x", selector)
}

// decodeStandardError decodes the standard Error(string) revert.
func decodeStandardError(data []byte) (*DecodedErrorResult, error) {
	if len(data) < 4 {
		return &DecodedErrorResult{
			ErrorName: "Error",
			Selector:  errorSelector,
		}, nil
	}

	if len(data) == 4 {
		return &DecodedErrorResult{
			ErrorName: "Error",
			Selector:  errorSelector,
		}, nil
	}

	// Decode the string
	reason, err := decodeRevertString(data[4:])
	if err != nil {
		return &DecodedErrorResult{
			ErrorName: "Error",
			Selector:  errorSelector,
		}, nil
	}

	return &DecodedErrorResult{
		ErrorName: "Error",
		Args:      []any{reason},
		Selector:  errorSelector,
	}, nil
}

// decodePanicError decodes the Panic(uint256) error.
func decodePanicError(data []byte) (*DecodedErrorResult, error) {
	if len(data) < 36 {
		return &DecodedErrorResult{
			ErrorName: "Panic",
			Selector:  panicSelector,
		}, nil
	}

	code := new(big.Int).SetBytes(data[4:36])
	return &DecodedErrorResult{
		ErrorName: "Panic",
		Args:      []any{code},
		Selector:  panicSelector,
	}, nil
}

// decodeRevertString decodes an ABI-encoded string from error data.
func decodeRevertString(data []byte) (string, error) {
	if len(data) < 64 {
		return "", fmt.Errorf("data too short for string")
	}

	// First 32 bytes is the offset (should be 32 for a single string)
	// Next 32 bytes is the length
	length := new(big.Int).SetBytes(data[32:64]).Uint64()

	if uint64(len(data)) < 64+length {
		return "", fmt.Errorf("data too short for string content")
	}

	return string(data[64 : 64+length]), nil
}

// DecodeError is an alias for DecodeErrorResult that returns the error name and args separately.
// Deprecated: Use DecodeErrorResult instead.
func (a *ABI) DecodeError(data []byte) (string, []any, error) {
	result, err := a.DecodeErrorResult(data)
	if err != nil {
		return "", nil, err
	}
	return result.ErrorName, result.Args, nil
}

// DecodeErrorResultWithoutABI decodes standard errors without requiring an ABI.
// Only handles Error(string) and Panic(uint256).
func DecodeErrorResultWithoutABI(data []byte) (*DecodedErrorResult, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("cannot decode zero data")
	}

	if len(data) < 4 {
		return nil, fmt.Errorf("error data too short: expected at least 4 bytes, got %d", len(data))
	}

	var selector [4]byte
	copy(selector[:], data[:4])

	// Check for standard Error(string)
	if selector == errorSelector {
		return decodeStandardError(data)
	}

	// Check for Panic(uint256)
	if selector == panicSelector {
		return decodePanicError(data)
	}

	return nil, fmt.Errorf("unknown error selector: 0x%x (no ABI provided for custom errors)", selector)
}
