package signature

import (
	"fmt"
	"math/big"
	"regexp"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// HashTypedData computes the EIP-712 hash of typed data.
// https://eips.ethereum.org/EIPS/eip-712
//
// Example:
//
//	hash, err := HashTypedData(TypedDataDefinition{
//		Domain: TypedDataDomain{
//			Name:    "Ether Mail",
//			Version: "1",
//			ChainId: big.NewInt(1),
//			VerifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
//		},
//		Types: map[string][]TypedDataField{
//			"Person": {
//				{Name: "name", Type: "string"},
//				{Name: "wallet", Type: "address"},
//			},
//			"Mail": {
//				{Name: "from", Type: "Person"},
//				{Name: "to", Type: "Person"},
//				{Name: "contents", Type: "string"},
//			},
//		},
//		PrimaryType: "Mail",
//		Message: map[string]any{
//			"from": map[string]any{
//				"name":   "Cow",
//				"wallet": "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
//			},
//			"to": map[string]any{
//				"name":   "Bob",
//				"wallet": "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB",
//			},
//			"contents": "Hello, Bob!",
//		},
//	})
func HashTypedData(data TypedDataDefinition) (string, error) {
	// Get domain types
	domainTypes := getTypesForEIP712Domain(data.Domain)

	// Create full types map with EIP712Domain
	types := make(map[string][]TypedDataField)
	for k, v := range data.Types {
		types[k] = v
	}
	types["EIP712Domain"] = domainTypes

	// Build the parts
	parts := []string{"0x1901"}

	// Hash the domain
	domainHash, err := hashStruct("EIP712Domain", domainToMessage(data.Domain), types)
	if err != nil {
		return "", fmt.Errorf("failed to hash domain: %w", err)
	}
	parts = append(parts, strings.TrimPrefix(domainHash, "0x"))

	// Hash the message if not EIP712Domain
	if data.PrimaryType != "EIP712Domain" {
		messageHash, err := hashStruct(data.PrimaryType, data.Message, types)
		if err != nil {
			return "", fmt.Errorf("failed to hash message: %w", err)
		}
		parts = append(parts, strings.TrimPrefix(messageHash, "0x"))
	}

	// Concatenate and hash
	combined := concatHex(parts...)
	return keccak256Hex(combined), nil
}

// HashDomain hashes just the EIP-712 domain.
func HashDomain(domain TypedDataDomain) (string, error) {
	types := map[string][]TypedDataField{
		"EIP712Domain": getTypesForEIP712Domain(domain),
	}
	return hashStruct("EIP712Domain", domainToMessage(domain), types)
}

// hashStruct hashes a struct according to EIP-712.
func hashStruct(primaryType string, data map[string]any, types map[string][]TypedDataField) (string, error) {
	encoded, err := encodeData(primaryType, data, types)
	if err != nil {
		return "", err
	}
	return keccak256Hex(encoded), nil
}

// encodeData encodes data according to EIP-712.
func encodeData(primaryType string, data map[string]any, types map[string][]TypedDataField) ([]byte, error) {
	// Start with the type hash
	typeHash := hashType(primaryType, types)

	// Build encoded values
	encodedValues := [][]byte{hexToBytes(typeHash)}

	fields := types[primaryType]
	for _, field := range fields {
		value := data[field.Name]
		encoded, err := encodeField(field.Type, value, types)
		if err != nil {
			return nil, fmt.Errorf("failed to encode field %s: %w", field.Name, err)
		}
		encodedValues = append(encodedValues, encoded)
	}

	// Concatenate all encoded values
	var result []byte
	for _, v := range encodedValues {
		result = append(result, v...)
	}

	return result, nil
}

// hashType computes the type hash for a struct.
func hashType(primaryType string, types map[string][]TypedDataField) string {
	encoded := encodeType(primaryType, types)
	return keccak256Hex(encoded)
}

// EncodeType encodes the type string for a struct.
func EncodeType(primaryType string, types map[string][]TypedDataField) string {
	return encodeType(primaryType, types)
}

// encodeType encodes the type string for a struct.
func encodeType(primaryType string, types map[string][]TypedDataField) string {
	// Find all dependencies
	deps := findTypeDependencies(primaryType, types, make(map[string]bool))
	delete(deps, primaryType)

	// Sort dependencies alphabetically
	sortedDeps := make([]string, 0, len(deps))
	for dep := range deps {
		sortedDeps = append(sortedDeps, dep)
	}
	sort.Strings(sortedDeps)

	// Primary type comes first, then sorted dependencies
	allTypes := append([]string{primaryType}, sortedDeps...)

	var result strings.Builder
	for _, typeName := range allTypes {
		fields := types[typeName]
		result.WriteString(typeName)
		result.WriteString("(")
		for i, field := range fields {
			if i > 0 {
				result.WriteString(",")
			}
			result.WriteString(field.Type)
			result.WriteString(" ")
			result.WriteString(field.Name)
		}
		result.WriteString(")")
	}

	return result.String()
}

// findTypeDependencies finds all type dependencies recursively.
func findTypeDependencies(primaryType string, types map[string][]TypedDataField, results map[string]bool) map[string]bool {
	// Extract base type (handle arrays like "Person[]")
	baseType := extractBaseType(primaryType)

	if results[baseType] || types[baseType] == nil {
		return results
	}

	results[baseType] = true

	for _, field := range types[baseType] {
		findTypeDependencies(field.Type, types, results)
	}

	return results
}

// extractBaseType extracts the base type from a type string (e.g., "Person[]" -> "Person").
func extractBaseType(typ string) string {
	re := regexp.MustCompile(`^\w*`)
	match := re.FindString(typ)
	return match
}

// encodeField encodes a single field value according to EIP-712.
func encodeField(fieldType string, value any, types map[string][]TypedDataField) ([]byte, error) {
	// Check if it's a custom type
	if types[fieldType] != nil {
		data, ok := value.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected map for type %s", fieldType)
		}
		encoded, err := encodeData(fieldType, data, types)
		if err != nil {
			return nil, err
		}
		return hexToBytes(keccak256Hex(encoded)), nil
	}

	// Handle bytes
	if fieldType == "bytes" {
		return hexToBytes(keccak256Hex(value)), nil
	}

	// Handle string
	if fieldType == "string" {
		str, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("expected string for type string")
		}
		return hexToBytes(keccak256Hex(str)), nil
	}

	// Handle arrays
	if strings.HasSuffix(fieldType, "]") {
		elementType := fieldType[:strings.LastIndex(fieldType, "[")]
		arr, ok := value.([]any)
		if !ok {
			return nil, fmt.Errorf("expected array for type %s", fieldType)
		}

		var encodedElements []byte
		for _, elem := range arr {
			encoded, err := encodeField(elementType, elem, types)
			if err != nil {
				return nil, err
			}
			encodedElements = append(encodedElements, encoded...)
		}
		return hexToBytes(keccak256Hex(encodedElements)), nil
	}

	// Handle basic types using ABI encoding
	return encodeAbiValue(fieldType, value)
}

// encodeAbiValue encodes a value using ABI encoding.
func encodeAbiValue(fieldType string, value any) ([]byte, error) {
	abiType, err := abi.NewType(fieldType, "", nil)
	if err != nil {
		return nil, fmt.Errorf("invalid type %s: %w", fieldType, err)
	}

	args := abi.Arguments{{Type: abiType}}

	// Convert value to expected type
	converted, err := convertTypedDataValue(fieldType, value)
	if err != nil {
		return nil, err
	}

	return args.Pack(converted)
}

// convertTypedDataValue converts a value to the expected Go type for ABI encoding.
func convertTypedDataValue(fieldType string, value any) (any, error) {
	switch {
	case fieldType == "address":
		switch v := value.(type) {
		case string:
			return common.HexToAddress(v), nil
		case common.Address:
			return v, nil
		default:
			return nil, fmt.Errorf("cannot convert %T to address", value)
		}

	case fieldType == "bool":
		switch v := value.(type) {
		case bool:
			return v, nil
		default:
			return nil, fmt.Errorf("cannot convert %T to bool", value)
		}

	case strings.HasPrefix(fieldType, "uint") || strings.HasPrefix(fieldType, "int"):
		switch v := value.(type) {
		case *big.Int:
			return v, nil
		case int:
			return big.NewInt(int64(v)), nil
		case int64:
			return big.NewInt(v), nil
		case float64:
			return big.NewInt(int64(v)), nil
		case string:
			n := new(big.Int)
			n.SetString(v, 10)
			return n, nil
		default:
			return nil, fmt.Errorf("cannot convert %T to integer", value)
		}

	case strings.HasPrefix(fieldType, "bytes"):
		switch v := value.(type) {
		case string:
			return hexToBytes(v), nil
		case []byte:
			return v, nil
		default:
			return nil, fmt.Errorf("cannot convert %T to bytes", value)
		}

	default:
		return value, nil
	}
}

// getTypesForEIP712Domain returns the types for the EIP712Domain struct.
func getTypesForEIP712Domain(domain TypedDataDomain) []TypedDataField {
	var types []TypedDataField

	if domain.Name != "" {
		types = append(types, TypedDataField{Name: "name", Type: "string"})
	}
	if domain.Version != "" {
		types = append(types, TypedDataField{Name: "version", Type: "string"})
	}
	if domain.ChainId != nil {
		types = append(types, TypedDataField{Name: "chainId", Type: "uint256"})
	}
	if domain.VerifyingContract != "" {
		types = append(types, TypedDataField{Name: "verifyingContract", Type: "address"})
	}
	if domain.Salt != "" {
		types = append(types, TypedDataField{Name: "salt", Type: "bytes32"})
	}

	return types
}

// domainToMessage converts a TypedDataDomain to a message map.
func domainToMessage(domain TypedDataDomain) map[string]any {
	message := make(map[string]any)

	if domain.Name != "" {
		message["name"] = domain.Name
	}
	if domain.Version != "" {
		message["version"] = domain.Version
	}
	if domain.ChainId != nil {
		message["chainId"] = domain.ChainId
	}
	if domain.VerifyingContract != "" {
		message["verifyingContract"] = domain.VerifyingContract
	}
	if domain.Salt != "" {
		message["salt"] = domain.Salt
	}

	return message
}
