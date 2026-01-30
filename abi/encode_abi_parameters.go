package abi

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// AbiParam represents a parameter definition for encoding/decoding.
// This mirrors viem's AbiParameter type.
type AbiParam struct {
	Name       string     `json:"name,omitempty"`
	Type       string     `json:"type"`
	Components []AbiParam `json:"components,omitempty"`
}

// EncodeAbiParameters encodes values based on ABI parameter definitions.
// This is a standalone function that doesn't require a parsed ABI.
//
// Example:
//
//	encoded, err := EncodeAbiParameters(
//	    []AbiParam{{Type: "uint256"}, {Type: "bool"}},
//	    []any{big.NewInt(420), true},
//	)
func EncodeAbiParameters(params []AbiParam, values []any) ([]byte, error) {
	if len(params) != len(values) {
		return nil, fmt.Errorf("params/values length mismatch: expected %d, got %d", len(params), len(values))
	}

	if len(params) == 0 {
		return []byte{}, nil
	}

	// Build go-ethereum Arguments from our params
	args, err := paramsToArguments(params)
	if err != nil {
		return nil, fmt.Errorf("failed to parse parameters: %w", err)
	}

	// Convert values to the expected types
	convertedValues, err := convertValues(params, values)
	if err != nil {
		return nil, fmt.Errorf("failed to convert values: %w", err)
	}

	// Pack the values
	packed, err := args.Pack(convertedValues...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode parameters: %w", err)
	}

	return packed, nil
}

// paramsToArguments converts AbiParam slice to go-ethereum Arguments.
func paramsToArguments(params []AbiParam) (abi.Arguments, error) {
	args := make(abi.Arguments, len(params))
	for i, param := range params {
		typ, err := parseType(param)
		if err != nil {
			return nil, fmt.Errorf("failed to parse type %q: %w", param.Type, err)
		}
		args[i] = abi.Argument{
			Name: param.Name,
			Type: typ,
		}
	}
	return args, nil
}

// parseType parses an AbiParam into a go-ethereum Type.
func parseType(param AbiParam) (abi.Type, error) {
	// Handle tuple types
	if param.Type == "tuple" || strings.HasPrefix(param.Type, "tuple[") {
		return parseTupleType(param)
	}

	// Parse regular types
	return abi.NewType(param.Type, "", nil)
}

// parseTupleType parses a tuple type with components.
func parseTupleType(param AbiParam) (abi.Type, error) {
	if len(param.Components) == 0 {
		return abi.Type{}, fmt.Errorf("tuple type requires components")
	}

	// Build component types
	components := make([]abi.ArgumentMarshaling, len(param.Components))
	for i, comp := range param.Components {
		components[i] = abiParamToArgumentMarshaling(comp)
	}

	return abi.NewType(param.Type, "", components)
}

// abiParamToArgumentMarshaling converts AbiParam to abi.ArgumentMarshaling.
func abiParamToArgumentMarshaling(param AbiParam) abi.ArgumentMarshaling {
	var components []abi.ArgumentMarshaling
	if len(param.Components) > 0 {
		components = make([]abi.ArgumentMarshaling, len(param.Components))
		for i, comp := range param.Components {
			components[i] = abiParamToArgumentMarshaling(comp)
		}
	}

	return abi.ArgumentMarshaling{
		Name:       param.Name,
		Type:       param.Type,
		Components: components,
	}
}

// convertValues converts Go values to the types expected by go-ethereum ABI packing.
func convertValues(params []AbiParam, values []any) ([]any, error) {
	result := make([]any, len(values))
	for i, value := range values {
		converted, err := convertValue(params[i], value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert value at index %d: %w", i, err)
		}
		result[i] = converted
	}
	return result, nil
}

// convertValue converts a single value to the expected type.
func convertValue(param AbiParam, value any) (any, error) {
	baseType := param.Type

	// Handle arrays
	if strings.HasSuffix(baseType, "[]") || strings.Contains(baseType, "[") {
		return convertArrayValue(param, value)
	}

	// Handle tuples
	if baseType == "tuple" {
		return convertTupleValue(param, value)
	}

	// Handle basic types
	return convertBasicValue(baseType, value)
}

// convertBasicValue converts basic Solidity types.
func convertBasicValue(typ string, value any) (any, error) {
	switch {
	case typ == "address":
		return convertAddress(value)
	case typ == "bool":
		return convertBool(value)
	case strings.HasPrefix(typ, "uint") || strings.HasPrefix(typ, "int"):
		return convertIntegerForType(typ, value)
	case strings.HasPrefix(typ, "bytes"):
		return convertBytes(typ, value)
	case typ == "string":
		return convertString(value)
	default:
		return value, nil
	}
}

// convertAddress converts a value to common.Address.
func convertAddress(value any) (common.Address, error) {
	switch v := value.(type) {
	case common.Address:
		return v, nil
	case string:
		if !common.IsHexAddress(v) {
			return common.Address{}, fmt.Errorf("invalid address: %s", v)
		}
		return common.HexToAddress(v), nil
	case []byte:
		if len(v) != 20 {
			return common.Address{}, fmt.Errorf("invalid address length: %d", len(v))
		}
		return common.BytesToAddress(v), nil
	default:
		return common.Address{}, fmt.Errorf("cannot convert %T to address", value)
	}
}

// convertBool converts a value to bool.
func convertBool(value any) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", value)
	}
}

// convertIntegerForType converts a value to the appropriate integer type for go-ethereum.
// go-ethereum expects native types for small integers (uint8, uint16, etc.)
// and *big.Int for larger integers (uint256, int256).
func convertIntegerForType(typ string, value any) (any, error) {
	// First, get the value as a big.Int for manipulation
	bi, err := convertToBigInt(value)
	if err != nil {
		return nil, err
	}

	// For uint256/int256 or unspecified size, use *big.Int
	switch typ {
	case "uint256", "int256", "uint", "int":
		return bi, nil
	case "uint8":
		return uint8(bi.Uint64()), nil
	case "uint16":
		return uint16(bi.Uint64()), nil
	case "uint32":
		return uint32(bi.Uint64()), nil
	case "uint64":
		return bi.Uint64(), nil
	case "int8":
		return int8(bi.Int64()), nil
	case "int16":
		return int16(bi.Int64()), nil
	case "int32":
		return int32(bi.Int64()), nil
	case "int64":
		return bi.Int64(), nil
	default:
		// For other integer types (uint24, uint40, etc.), use *big.Int
		return bi, nil
	}
}

// convertToBigInt converts a value to *big.Int.
func convertToBigInt(value any) (*big.Int, error) {
	switch v := value.(type) {
	case *big.Int:
		return v, nil
	case int:
		return big.NewInt(int64(v)), nil
	case int8:
		return big.NewInt(int64(v)), nil
	case int16:
		return big.NewInt(int64(v)), nil
	case int32:
		return big.NewInt(int64(v)), nil
	case int64:
		return big.NewInt(v), nil
	case uint:
		return new(big.Int).SetUint64(uint64(v)), nil
	case uint8:
		return new(big.Int).SetUint64(uint64(v)), nil
	case uint16:
		return new(big.Int).SetUint64(uint64(v)), nil
	case uint32:
		return new(big.Int).SetUint64(uint64(v)), nil
	case uint64:
		return new(big.Int).SetUint64(v), nil
	default:
		return nil, fmt.Errorf("cannot convert %T to integer", value)
	}
}

// convertBytes converts a value to []byte or fixed bytes.
func convertBytes(typ string, value any) (any, error) {
	// For fixed bytes types, just pass through if it's already the right type
	if typ != "bytes" && strings.HasPrefix(typ, "bytes") {
		// Fixed bytes (bytes1 to bytes32) - pass through arrays
		switch v := value.(type) {
		case [1]byte, [2]byte, [3]byte, [4]byte, [5]byte, [6]byte, [7]byte, [8]byte,
			[9]byte, [10]byte, [11]byte, [12]byte, [13]byte, [14]byte, [15]byte, [16]byte,
			[17]byte, [18]byte, [19]byte, [20]byte, [21]byte, [22]byte, [23]byte, [24]byte,
			[25]byte, [26]byte, [27]byte, [28]byte, [29]byte, [30]byte, [31]byte, [32]byte:
			return v, nil
		case []byte:
			return v, nil
		case string:
			if strings.HasPrefix(v, "0x") || strings.HasPrefix(v, "0X") {
				return common.FromHex(v), nil
			}
			return []byte(v), nil
		default:
			return nil, fmt.Errorf("cannot convert %T to bytes", value)
		}
	}

	// Dynamic bytes
	switch v := value.(type) {
	case []byte:
		return v, nil
	case string:
		// Handle hex strings
		if strings.HasPrefix(v, "0x") || strings.HasPrefix(v, "0X") {
			return common.FromHex(v), nil
		}
		return []byte(v), nil
	default:
		return nil, fmt.Errorf("cannot convert %T to bytes", value)
	}
}

// convertString converts a value to string.
func convertString(value any) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		return "", fmt.Errorf("cannot convert %T to string", value)
	}
}

// convertArrayValue converts array values.
func convertArrayValue(param AbiParam, value any) (any, error) {
	// Get the element type
	elementType := getArrayElementType(param.Type)
	elementParam := AbiParam{
		Type:       elementType,
		Components: param.Components,
	}

	// Handle the array
	switch v := value.(type) {
	case []any:
		result := make([]any, len(v))
		for i, elem := range v {
			converted, err := convertValue(elementParam, elem)
			if err != nil {
				return nil, fmt.Errorf("failed to convert array element %d: %w", i, err)
			}
			result[i] = converted
		}
		return result, nil
	default:
		// For typed slices, return as-is (go-ethereum can handle them)
		return value, nil
	}
}

// getArrayElementType extracts the element type from an array type.
func getArrayElementType(typ string) string {
	// Handle dynamic arrays like "uint256[]"
	if strings.HasSuffix(typ, "[]") {
		return strings.TrimSuffix(typ, "[]")
	}
	// Handle fixed arrays like "uint256[3]"
	idx := strings.LastIndex(typ, "[")
	if idx > 0 {
		return typ[:idx]
	}
	return typ
}

// convertTupleValue converts tuple/struct values.
func convertTupleValue(param AbiParam, value any) (any, error) {
	switch v := value.(type) {
	case map[string]any:
		// Named tuple - convert to slice in component order
		result := make([]any, len(param.Components))
		for i, comp := range param.Components {
			val, ok := v[comp.Name]
			if !ok {
				return nil, fmt.Errorf("missing tuple field: %s", comp.Name)
			}
			converted, err := convertValue(comp, val)
			if err != nil {
				return nil, err
			}
			result[i] = converted
		}
		return result, nil
	case []any:
		// Unnamed tuple - convert each element
		if len(v) != len(param.Components) {
			return nil, fmt.Errorf("tuple length mismatch: expected %d, got %d", len(param.Components), len(v))
		}
		result := make([]any, len(v))
		for i, elem := range v {
			converted, err := convertValue(param.Components[i], elem)
			if err != nil {
				return nil, err
			}
			result[i] = converted
		}
		return result, nil
	default:
		// Might be a struct - return as-is and let go-ethereum handle it
		return value, nil
	}
}
