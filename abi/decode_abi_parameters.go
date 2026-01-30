package abi

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// DecodeAbiParameters decodes ABI-encoded data based on parameter definitions.
// This is a standalone function that doesn't require a parsed ABI.
//
// Example:
//
//	decoded, err := DecodeAbiParameters(
//	    []AbiParam{{Type: "uint256"}, {Type: "bool"}},
//	    hexData,
//	)
func DecodeAbiParameters(params []AbiParam, data []byte) ([]any, error) {
	if len(params) == 0 {
		return []any{}, nil
	}

	if len(data) == 0 && len(params) > 0 {
		return nil, fmt.Errorf("cannot decode zero data with non-empty params")
	}

	if len(data) > 0 && len(data) < 32 {
		return nil, fmt.Errorf("data too small: expected at least 32 bytes, got %d", len(data))
	}

	// Build go-ethereum Arguments from our params
	args, err := paramsToArguments(params)
	if err != nil {
		return nil, fmt.Errorf("failed to parse parameters: %w", err)
	}

	// Unpack the data
	unpacked, err := args.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode parameters: %w", err)
	}

	// Convert unpacked values to standard Go types
	result := make([]any, len(unpacked))
	for i, val := range unpacked {
		result[i] = normalizeDecodedValue(params[i], val)
	}

	return result, nil
}

// normalizeDecodedValue converts go-ethereum decoded values to more convenient Go types.
func normalizeDecodedValue(param AbiParam, value any) any {
	if value == nil {
		return nil
	}

	baseType := param.Type

	// Handle arrays
	if strings.HasSuffix(baseType, "[]") || (strings.Contains(baseType, "[") && !strings.HasPrefix(baseType, "bytes")) {
		return normalizeArrayValue(param, value)
	}

	// Handle tuples
	if baseType == "tuple" {
		return normalizeTupleValue(param, value)
	}

	// Handle basic types
	return normalizeBasicValue(baseType, value)
}

// normalizeBasicValue normalizes basic decoded values.
func normalizeBasicValue(typ string, value any) any {
	switch {
	case typ == "address":
		if addr, ok := value.(common.Address); ok {
			return addr
		}
		return value

	case typ == "bool":
		if b, ok := value.(bool); ok {
			return b
		}
		return value

	case strings.HasPrefix(typ, "uint") || strings.HasPrefix(typ, "int"):
		return normalizeIntegerValue(typ, value)

	case strings.HasPrefix(typ, "bytes"):
		return normalizeBytesValue(typ, value)

	case typ == "string":
		if s, ok := value.(string); ok {
			return s
		}
		return value

	default:
		return value
	}
}

// normalizeIntegerValue converts integer values.
// For small integers (<=48 bits), returns int/int64
// For larger integers, returns *big.Int
func normalizeIntegerValue(typ string, value any) any {
	// Extract bit size from type
	size := 256 // default
	if typ != "uint" && typ != "int" {
		var prefix string
		if strings.HasPrefix(typ, "uint") {
			prefix = "uint"
		} else {
			prefix = "int"
		}
		fmt.Sscanf(typ[len(prefix):], "%d", &size)
	}

	// Get the value as *big.Int
	var bi *big.Int
	switch v := value.(type) {
	case *big.Int:
		bi = v
	case int64:
		bi = big.NewInt(v)
	case uint64:
		bi = new(big.Int).SetUint64(v)
	case int32:
		bi = big.NewInt(int64(v))
	case uint32:
		bi = new(big.Int).SetUint64(uint64(v))
	case int16:
		bi = big.NewInt(int64(v))
	case uint16:
		bi = new(big.Int).SetUint64(uint64(v))
	case int8:
		bi = big.NewInt(int64(v))
	case uint8:
		bi = new(big.Int).SetUint64(uint64(v))
	default:
		return value
	}

	// For small integers, return native types
	if size <= 48 {
		if strings.HasPrefix(typ, "uint") {
			if bi.IsUint64() && bi.Uint64() <= uint64(1<<48-1) {
				return bi.Int64() // Use int64 for consistency
			}
		} else {
			if bi.IsInt64() {
				return bi.Int64()
			}
		}
	}

	return bi
}

// normalizeBytesValue normalizes bytes values.
func normalizeBytesValue(typ string, value any) any {
	// For fixed bytes (bytes1 to bytes32), return as hex string for viem compatibility
	if typ != "bytes" {
		switch v := value.(type) {
		case [1]byte:
			return common.Bytes2Hex(v[:])
		case [2]byte:
			return "0x" + common.Bytes2Hex(v[:])
		case [3]byte:
			return "0x" + common.Bytes2Hex(v[:])
		case [4]byte:
			return "0x" + common.Bytes2Hex(v[:])
		case [8]byte:
			return "0x" + common.Bytes2Hex(v[:])
		case [16]byte:
			return "0x" + common.Bytes2Hex(v[:])
		case [20]byte:
			return "0x" + common.Bytes2Hex(v[:])
		case [32]byte:
			return "0x" + common.Bytes2Hex(v[:])
		case []byte:
			return "0x" + common.Bytes2Hex(v)
		}
	}

	// For dynamic bytes, return as []byte
	if b, ok := value.([]byte); ok {
		return b
	}

	return value
}

// normalizeArrayValue normalizes array/slice values.
func normalizeArrayValue(param AbiParam, value any) any {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return value
	}

	// Get element type
	elementType := getArrayElementType(param.Type)
	elementParam := AbiParam{
		Type:       elementType,
		Components: param.Components,
	}

	// Convert each element
	result := make([]any, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		result[i] = normalizeDecodedValue(elementParam, rv.Index(i).Interface())
	}

	return result
}

// normalizeTupleValue normalizes tuple/struct values.
func normalizeTupleValue(param AbiParam, value any) any {
	rv := reflect.ValueOf(value)

	// Handle struct types (go-ethereum returns structs for tuples)
	if rv.Kind() == reflect.Struct {
		return normalizeStructToMap(param, rv)
	}

	// Handle slice/array (unnamed tuples)
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		result := make([]any, rv.Len())
		for i := 0; i < rv.Len() && i < len(param.Components); i++ {
			result[i] = normalizeDecodedValue(param.Components[i], rv.Index(i).Interface())
		}
		return result
	}

	return value
}

// normalizeStructToMap converts a struct to a map with component names.
func normalizeStructToMap(param AbiParam, rv reflect.Value) any {
	// Check if all components have names
	hasNames := true
	for _, comp := range param.Components {
		if comp.Name == "" {
			hasNames = false
			break
		}
	}

	if hasNames {
		// Return as map
		result := make(map[string]any)
		for i := 0; i < rv.NumField() && i < len(param.Components); i++ {
			comp := param.Components[i]
			fieldVal := rv.Field(i).Interface()
			result[comp.Name] = normalizeDecodedValue(comp, fieldVal)
		}
		return result
	}

	// Return as array for unnamed tuples
	result := make([]any, rv.NumField())
	for i := 0; i < rv.NumField() && i < len(param.Components); i++ {
		result[i] = normalizeDecodedValue(param.Components[i], rv.Field(i).Interface())
	}
	return result
}

// DecodeAbiParametersInto decodes ABI-encoded data into a provided struct or slice.
func DecodeAbiParametersInto(params []AbiParam, data []byte, output any) error {
	if len(params) == 0 {
		return nil
	}

	if len(data) == 0 && len(params) > 0 {
		return fmt.Errorf("cannot decode zero data with non-empty params")
	}

	// Build go-ethereum Arguments from our params
	args, err := paramsToArguments(params)
	if err != nil {
		return fmt.Errorf("failed to parse parameters: %w", err)
	}

	// Unpack the data
	unpacked, err := args.Unpack(data)
	if err != nil {
		return fmt.Errorf("failed to decode parameters: %w", err)
	}

	// Copy to output using reflection
	rv := reflect.ValueOf(output)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("output must be a pointer")
	}
	rv = rv.Elem()

	switch rv.Kind() {
	case reflect.Slice:
		if rv.Len() < len(unpacked) {
			rv.Set(reflect.MakeSlice(rv.Type(), len(unpacked), len(unpacked)))
		}
		for i, v := range unpacked {
			rv.Index(i).Set(reflect.ValueOf(v))
		}
	case reflect.Struct:
		for i := 0; i < rv.NumField() && i < len(unpacked); i++ {
			rv.Field(i).Set(reflect.ValueOf(unpacked[i]))
		}
	default:
		if len(unpacked) == 1 {
			rv.Set(reflect.ValueOf(unpacked[0]))
		}
	}

	return nil
}

// DecodeWithSelector decodes data that starts with a 4-byte selector.
// Returns the selector and decoded parameters.
func DecodeWithSelector(params []AbiParam, data []byte) ([4]byte, []any, error) {
	if len(data) < 4 {
		return [4]byte{}, nil, fmt.Errorf("data too short: expected at least 4 bytes, got %d", len(data))
	}

	var selector [4]byte
	copy(selector[:], data[:4])

	decoded, err := DecodeAbiParameters(params, data[4:])
	if err != nil {
		return selector, nil, err
	}

	return selector, decoded, nil
}
