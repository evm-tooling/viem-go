package abi

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
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
//
// For functions with multiple return values, the output should be a pointer to a struct
// where fields are mapped by position to the function outputs:
//
//	type ReservesResult struct {
//	    Reserve0       *big.Int
//	    Reserve1       *big.Int
//	    BlockTimestamp uint32
//	}
//	var result ReservesResult
//	err := abi.DecodeFunctionResultInto("getReserves", data, &result)
//
// Field names are matched case-insensitively to the ABI output names, or by position
// if names don't match.
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

	// Try the standard go-ethereum unpack first
	err := a.gethABI.UnpackIntoInterface(output, functionName, data)
	if err == nil {
		return nil
	}

	// If standard unpack fails, try our custom struct binding
	return a.decodeIntoStruct(functionName, data, output)
}

// decodeIntoStruct provides custom struct binding for multi-value returns.
func (a *ABI) decodeIntoStruct(functionName string, data []byte, output any) error {
	// First decode the raw values
	unpacked, err := a.DecodeFunctionResult(functionName, data)
	if err != nil {
		return err
	}

	// Get the output value and ensure it's a pointer to struct
	outVal := reflect.ValueOf(output)
	if outVal.Kind() != reflect.Ptr {
		return fmt.Errorf("output must be a pointer, got %T", output)
	}

	outVal = outVal.Elem()
	if outVal.Kind() != reflect.Struct {
		// For non-struct types, try direct assignment
		if len(unpacked) == 1 {
			return assignValue(outVal, unpacked[0])
		}
		return fmt.Errorf("output must be a struct for multiple return values, got %s", outVal.Kind())
	}

	outType := outVal.Type()
	numFields := outType.NumField()

	// Ensure we have enough values
	if len(unpacked) > numFields {
		return fmt.Errorf("struct has %d fields but function returns %d values", numFields, len(unpacked))
	}

	// Try to match by name first, then by position
	m := a.gethABI.Methods[functionName]
	used := make([]bool, len(unpacked))

	// First pass: match by name
	for i := 0; i < numFields; i++ {
		field := outType.Field(i)
		if !field.IsExported() {
			continue
		}

		// Check for abi tag
		abiName := field.Tag.Get("abi")
		if abiName == "" {
			abiName = field.Name
		}

		// Find matching output by name
		for j, outParam := range m.Outputs {
			if !used[j] && matchesName(abiName, outParam.Name) {
				if err := assignValue(outVal.Field(i), unpacked[j]); err != nil {
					return fmt.Errorf("failed to assign field %s: %w", field.Name, err)
				}
				used[j] = true
				break
			}
		}
	}

	// Second pass: assign remaining by position
	unpackedIdx := 0
	for i := 0; i < numFields && unpackedIdx < len(unpacked); i++ {
		field := outType.Field(i)
		if !field.IsExported() {
			continue
		}

		// Skip already used values
		for unpackedIdx < len(unpacked) && used[unpackedIdx] {
			unpackedIdx++
		}
		if unpackedIdx >= len(unpacked) {
			break
		}

		// Check if this field was already assigned by name
		fieldVal := outVal.Field(i)
		if !fieldVal.IsZero() {
			continue
		}

		if err := assignValue(fieldVal, unpacked[unpackedIdx]); err != nil {
			return fmt.Errorf("failed to assign field %s: %w", field.Name, err)
		}
		used[unpackedIdx] = true
		unpackedIdx++
	}

	return nil
}

// matchesName checks if two names match (case-insensitive).
func matchesName(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	return equalFold(a, b)
}

// equalFold is a simple case-insensitive string comparison.
func equalFold(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		c1, c2 := s[i], t[i]
		if c1 >= 'A' && c1 <= 'Z' {
			c1 += 'a' - 'A'
		}
		if c2 >= 'A' && c2 <= 'Z' {
			c2 += 'a' - 'A'
		}
		if c1 != c2 {
			return false
		}
	}
	return true
}

// assignValue assigns a value to a reflect.Value with type conversion.
func assignValue(dst reflect.Value, src any) error {
	if src == nil {
		return nil
	}

	srcVal := reflect.ValueOf(src)
	dstType := dst.Type()

	// Direct assignment if types match
	if srcVal.Type().AssignableTo(dstType) {
		dst.Set(srcVal)
		return nil
	}

	// Handle common conversions
	switch dstType.Kind() {
	case reflect.Ptr:
		// Handle *big.Int
		if dstType == reflect.TypeOf((*big.Int)(nil)) {
			switch v := src.(type) {
			case *big.Int:
				dst.Set(reflect.ValueOf(v))
				return nil
			case int64:
				dst.Set(reflect.ValueOf(big.NewInt(v)))
				return nil
			case uint64:
				dst.Set(reflect.ValueOf(new(big.Int).SetUint64(v)))
				return nil
			}
		}

	case reflect.Uint8:
		switch v := src.(type) {
		case uint8:
			dst.SetUint(uint64(v))
			return nil
		case *big.Int:
			dst.SetUint(v.Uint64())
			return nil
		}

	case reflect.Uint16:
		switch v := src.(type) {
		case uint16:
			dst.SetUint(uint64(v))
			return nil
		case *big.Int:
			dst.SetUint(v.Uint64())
			return nil
		}

	case reflect.Uint32:
		switch v := src.(type) {
		case uint32:
			dst.SetUint(uint64(v))
			return nil
		case *big.Int:
			dst.SetUint(v.Uint64())
			return nil
		}

	case reflect.Uint64:
		switch v := src.(type) {
		case uint64:
			dst.SetUint(v)
			return nil
		case *big.Int:
			dst.SetUint(v.Uint64())
			return nil
		}

	case reflect.Int64:
		switch v := src.(type) {
		case int64:
			dst.SetInt(v)
			return nil
		case *big.Int:
			dst.SetInt(v.Int64())
			return nil
		}

	case reflect.Bool:
		if b, ok := src.(bool); ok {
			dst.SetBool(b)
			return nil
		}

	case reflect.String:
		if s, ok := src.(string); ok {
			dst.SetString(s)
			return nil
		}

	case reflect.Slice:
		if dstType.Elem().Kind() == reflect.Uint8 {
			if b, ok := src.([]byte); ok {
				dst.SetBytes(b)
				return nil
			}
		}

	case reflect.Array:
		// Handle [N]byte types like common.Address, common.Hash
		if dstType == reflect.TypeOf(common.Address{}) {
			if addr, ok := src.(common.Address); ok {
				dst.Set(reflect.ValueOf(addr))
				return nil
			}
		}
		if dstType == reflect.TypeOf(common.Hash{}) {
			if hash, ok := src.(common.Hash); ok {
				dst.Set(reflect.ValueOf(hash))
				return nil
			}
		}
		// Handle [32]byte
		if dstType.Len() == 32 && dstType.Elem().Kind() == reflect.Uint8 {
			switch v := src.(type) {
			case [32]byte:
				dst.Set(reflect.ValueOf(v))
				return nil
			case []byte:
				if len(v) == 32 {
					var arr [32]byte
					copy(arr[:], v)
					dst.Set(reflect.ValueOf(arr))
					return nil
				}
			}
		}
	}

	// Try convertible types
	if srcVal.Type().ConvertibleTo(dstType) {
		dst.Set(srcVal.Convert(dstType))
		return nil
	}

	return fmt.Errorf("cannot assign %T to %s", src, dstType)
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
