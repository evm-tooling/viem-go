package abi

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

var (
	integerRegex = regexp.MustCompile(`^(u?int)(\d*)$`)
	bytesRegex   = regexp.MustCompile(`^bytes(\d+)$`)
	arrayRegex   = regexp.MustCompile(`^(.+)\[\d*\]$`)
)

// EncodePacked performs non-padded ABI encoding, similar to Solidity's abi.encodePacked.
// Unlike standard ABI encoding, packed encoding does not pad values to 32 bytes.
//
// Example:
//
//	encoded, err := EncodePacked(
//	    []string{"address", "uint256"},
//	    []any{"0x14dC79964da2C08b23698B3D3cc7Ca32193d9955", big.NewInt(420)},
//	)
func EncodePacked(types []string, values []any) ([]byte, error) {
	if len(types) != len(values) {
		return nil, fmt.Errorf("types/values length mismatch: expected %d, got %d", len(types), len(values))
	}

	var result []byte
	for i, typ := range types {
		encoded, err := encodePackedValue(typ, values[i], false)
		if err != nil {
			return nil, fmt.Errorf("failed to encode value at index %d: %w", i, err)
		}
		result = append(result, encoded...)
	}

	return result, nil
}

// encodePackedValue encodes a single value in packed format.
func encodePackedValue(typ string, value any, isArrayElement bool) ([]byte, error) {
	// Handle address
	if typ == "address" {
		return encodePackedAddress(value, isArrayElement)
	}

	// Handle string
	if typ == "string" {
		return encodePackedString(value)
	}

	// Handle bytes (dynamic)
	if typ == "bytes" {
		return encodePackedDynamicBytes(value)
	}

	// Handle bool
	if typ == "bool" {
		return encodePackedBool(value, isArrayElement)
	}

	// Handle integers (uint*, int*)
	if match := integerRegex.FindStringSubmatch(typ); match != nil {
		return encodePackedInteger(match, value, isArrayElement)
	}

	// Handle fixed bytes (bytes1 to bytes32)
	if match := bytesRegex.FindStringSubmatch(typ); match != nil {
		return encodePackedFixedBytes(match, value, isArrayElement)
	}

	// Handle arrays
	if match := arrayRegex.FindStringSubmatch(typ); match != nil {
		return encodePackedArray(match[1], value)
	}

	return nil, fmt.Errorf("unsupported packed encoding type: %s", typ)
}

// encodePackedAddress encodes an address (20 bytes, or 32 if array element).
func encodePackedAddress(value any, isArrayElement bool) ([]byte, error) {
	var addr common.Address

	switch v := value.(type) {
	case common.Address:
		addr = v
	case string:
		if !common.IsHexAddress(v) {
			return nil, fmt.Errorf("invalid address: %s", v)
		}
		addr = common.HexToAddress(v)
	case []byte:
		if len(v) != 20 {
			return nil, fmt.Errorf("invalid address length: %d", len(v))
		}
		addr = common.BytesToAddress(v)
	default:
		return nil, fmt.Errorf("cannot convert %T to address", value)
	}

	if isArrayElement {
		// Pad to 32 bytes for array elements
		result := make([]byte, 32)
		copy(result[12:], addr.Bytes())
		return result, nil
	}

	return addr.Bytes(), nil
}

// encodePackedString encodes a string (no padding).
func encodePackedString(value any) ([]byte, error) {
	switch v := value.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to string", value)
	}
}

// encodePackedDynamicBytes encodes dynamic bytes (no padding).
func encodePackedDynamicBytes(value any) ([]byte, error) {
	switch v := value.(type) {
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

// encodePackedBool encodes a boolean (1 byte, or 32 if array element).
func encodePackedBool(value any, isArrayElement bool) ([]byte, error) {
	var b bool

	switch v := value.(type) {
	case bool:
		b = v
	default:
		return nil, fmt.Errorf("cannot convert %T to bool", value)
	}

	if isArrayElement {
		result := make([]byte, 32)
		if b {
			result[31] = 1
		}
		return result, nil
	}

	if b {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

// encodePackedInteger encodes an integer (variable size based on type).
func encodePackedInteger(match []string, value any, isArrayElement bool) ([]byte, error) {
	baseType := match[1]
	bitsStr := match[2]
	bits := 256
	if bitsStr != "" {
		var err error
		bits, err = strconv.Atoi(bitsStr)
		if err != nil {
			return nil, fmt.Errorf("invalid integer size: %s", bitsStr)
		}
	}

	size := bits / 8
	signed := baseType == "int"

	// Convert value to *big.Int
	var bi *big.Int
	switch v := value.(type) {
	case *big.Int:
		bi = v
	case int:
		bi = big.NewInt(int64(v))
	case int8:
		bi = big.NewInt(int64(v))
	case int16:
		bi = big.NewInt(int64(v))
	case int32:
		bi = big.NewInt(int64(v))
	case int64:
		bi = big.NewInt(v)
	case uint:
		bi = new(big.Int).SetUint64(uint64(v))
	case uint8:
		bi = new(big.Int).SetUint64(uint64(v))
	case uint16:
		bi = new(big.Int).SetUint64(uint64(v))
	case uint32:
		bi = new(big.Int).SetUint64(uint64(v))
	case uint64:
		bi = new(big.Int).SetUint64(v)
	default:
		return nil, fmt.Errorf("cannot convert %T to integer", value)
	}

	if isArrayElement {
		size = 32
	}

	// Handle signed integers with two's complement
	if signed && bi.Sign() < 0 {
		// Two's complement for negative numbers
		twosComplement := new(big.Int).Add(new(big.Int).Lsh(big.NewInt(1), uint(size*8)), bi)
		b := twosComplement.Bytes()
		if len(b) > size {
			b = b[len(b)-size:]
		}
		result := make([]byte, size)
		copy(result[size-len(b):], b)
		return result, nil
	}

	b := bi.Bytes()
	if len(b) > size {
		return nil, fmt.Errorf("integer overflow for %s", match[0])
	}

	result := make([]byte, size)
	copy(result[size-len(b):], b)
	return result, nil
}

// encodePackedFixedBytes encodes fixed bytes (bytesN).
func encodePackedFixedBytes(match []string, value any, isArrayElement bool) ([]byte, error) {
	size, err := strconv.Atoi(match[1])
	if err != nil {
		return nil, fmt.Errorf("invalid bytes size: %s", match[1])
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		if strings.HasPrefix(v, "0x") || strings.HasPrefix(v, "0X") {
			data = common.FromHex(v)
		} else {
			data = []byte(v)
		}
	default:
		return nil, fmt.Errorf("cannot convert %T to bytes", value)
	}

	if len(data) != size {
		return nil, fmt.Errorf("bytes size mismatch: expected %d, got %d", size, len(data))
	}

	if isArrayElement {
		// Pad to 32 bytes (right-padded) for array elements
		result := make([]byte, 32)
		copy(result, data)
		return result, nil
	}

	return data, nil
}

// encodePackedArray encodes an array of values.
func encodePackedArray(elementType string, value any) ([]byte, error) {
	// Handle the array based on type
	var result []byte

	switch v := value.(type) {
	case []any:
		for i, elem := range v {
			encoded, err := encodePackedValue(elementType, elem, true)
			if err != nil {
				return nil, fmt.Errorf("failed to encode array element %d: %w", i, err)
			}
			result = append(result, encoded...)
		}
	case []string:
		for i, elem := range v {
			encoded, err := encodePackedValue(elementType, elem, true)
			if err != nil {
				return nil, fmt.Errorf("failed to encode array element %d: %w", i, err)
			}
			result = append(result, encoded...)
		}
	case []*big.Int:
		for i, elem := range v {
			encoded, err := encodePackedValue(elementType, elem, true)
			if err != nil {
				return nil, fmt.Errorf("failed to encode array element %d: %w", i, err)
			}
			result = append(result, encoded...)
		}
	case []common.Address:
		for i, elem := range v {
			encoded, err := encodePackedValue(elementType, elem, true)
			if err != nil {
				return nil, fmt.Errorf("failed to encode array element %d: %w", i, err)
			}
			result = append(result, encoded...)
		}
	case []bool:
		for i, elem := range v {
			encoded, err := encodePackedValue(elementType, elem, true)
			if err != nil {
				return nil, fmt.Errorf("failed to encode array element %d: %w", i, err)
			}
			result = append(result, encoded...)
		}
	case [][]byte:
		for i, elem := range v {
			encoded, err := encodePackedValue(elementType, elem, true)
			if err != nil {
				return nil, fmt.Errorf("failed to encode array element %d: %w", i, err)
			}
			result = append(result, encoded...)
		}
	default:
		return nil, fmt.Errorf("unsupported array type: %T", value)
	}

	if len(result) == 0 {
		return []byte{}, nil
	}

	return result, nil
}
