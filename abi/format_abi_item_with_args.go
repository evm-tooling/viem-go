package abi

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// FormatAbiItemWithArgsOptions configures the output format.
type FormatAbiItemWithArgsOptions struct {
	// IncludeFunctionName includes the function/event name in the output.
	// Default is true.
	IncludeFunctionName bool
	// IncludeName includes parameter names in the output.
	// Default is false.
	IncludeName bool
}

// FormatAbiItemWithArgs formats an ABI item with its actual argument values.
// This is useful for logging and debugging.
//
// Example:
//
//	formatted := FormatAbiItemWithArgs(function, args, nil)
//	// "transfer(0x1234..., 1000000000000000000)"
//
//	formatted := FormatAbiItemWithArgs(function, args, &FormatAbiItemWithArgsOptions{IncludeName: true})
//	// "transfer(to: 0x1234..., amount: 1000000000000000000)"
func FormatAbiItemWithArgs(item AbiItem, args []any, opts *FormatAbiItemWithArgsOptions) string {
	if opts == nil {
		opts = &FormatAbiItemWithArgsOptions{
			IncludeFunctionName: true,
			IncludeName:         false,
		}
	}

	var name string
	var inputs []Parameter

	switch v := item.(type) {
	case Function:
		name = v.Name
		inputs = v.Inputs
	case *Function:
		name = v.Name
		inputs = v.Inputs
	case Event:
		name = v.Name
		inputs = v.Inputs
	case *Event:
		name = v.Name
		inputs = v.Inputs
	case Error:
		name = v.Name
		inputs = v.Inputs
	case *Error:
		name = v.Name
		inputs = v.Inputs
	default:
		return ""
	}

	if inputs == nil {
		if opts.IncludeFunctionName {
			return name + "()"
		}
		return "()"
	}

	// Format each argument
	formattedArgs := make([]string, len(inputs))
	for i, input := range inputs {
		var argStr string
		if i < len(args) {
			argStr = formatArgValue(args[i])
		} else {
			argStr = "undefined"
		}

		if opts.IncludeName && input.Name != "" {
			formattedArgs[i] = input.Name + ": " + argStr
		} else {
			formattedArgs[i] = argStr
		}
	}

	if opts.IncludeFunctionName {
		return name + "(" + strings.Join(formattedArgs, ", ") + ")"
	}
	return "(" + strings.Join(formattedArgs, ", ") + ")"
}

// FormatFunctionCallWithArgs is a convenience function for formatting function calls.
func (a *ABI) FormatFunctionCallWithArgs(functionName string, args []any, opts *FormatAbiItemWithArgsOptions) (string, error) {
	fn, err := a.GetFunction(functionName)
	if err != nil {
		return "", err
	}
	return FormatAbiItemWithArgs(fn, args, opts), nil
}

// formatArgValue formats a single argument value for display.
func formatArgValue(arg any) string {
	if arg == nil {
		return "null"
	}

	switch v := arg.(type) {
	case common.Address:
		return v.Hex()
	case *common.Address:
		if v == nil {
			return "null"
		}
		return v.Hex()
	case common.Hash:
		return v.Hex()
	case *big.Int:
		if v == nil {
			return "null"
		}
		return v.String()
	case bool:
		if v {
			return "true"
		}
		return "false"
	case string:
		// Check if it's an address
		if common.IsHexAddress(v) {
			return v
		}
		// Quote strings
		return fmt.Sprintf("%q", v)
	case []byte:
		if len(v) == 0 {
			return "0x"
		}
		return "0x" + common.Bytes2Hex(v)
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case [1]byte, [2]byte, [3]byte, [4]byte, [8]byte, [16]byte, [20]byte, [32]byte:
		return formatFixedBytes(v)
	default:
		// For complex types (arrays, structs), use JSON
		return formatComplexValue(v)
	}
}

// formatFixedBytes formats fixed-size byte arrays.
func formatFixedBytes(v any) string {
	switch b := v.(type) {
	case [1]byte:
		return "0x" + common.Bytes2Hex(b[:])
	case [2]byte:
		return "0x" + common.Bytes2Hex(b[:])
	case [3]byte:
		return "0x" + common.Bytes2Hex(b[:])
	case [4]byte:
		return "0x" + common.Bytes2Hex(b[:])
	case [8]byte:
		return "0x" + common.Bytes2Hex(b[:])
	case [16]byte:
		return "0x" + common.Bytes2Hex(b[:])
	case [20]byte:
		return "0x" + common.Bytes2Hex(b[:])
	case [32]byte:
		return "0x" + common.Bytes2Hex(b[:])
	default:
		return fmt.Sprintf("%v", v)
	}
}

// formatComplexValue formats arrays, slices, and structs.
func formatComplexValue(v any) string {
	// Try JSON marshaling for complex types
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}
