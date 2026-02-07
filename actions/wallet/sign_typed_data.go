package wallet

import (
	"context"
	"fmt"
	"strings"

	json "github.com/goccy/go-json"

	"github.com/ChefBingbong/viem-go/utils/signature"
)

// SignTypedDataParameters contains the parameters for the SignTypedData action.
// This mirrors viem's SignTypedDataParameters type.
type SignTypedDataParameters struct {
	// Account is the account to sign with. If nil, uses the client's account.
	Account Account

	// Domain contains the EIP-712 domain parameters.
	Domain signature.TypedDataDomain

	// Types contains the type definitions (excluding EIP712Domain, which is auto-generated).
	Types map[string][]signature.TypedDataField

	// PrimaryType is the primary type being signed.
	PrimaryType string

	// Message is the structured message to sign.
	Message map[string]any
}

// SignTypedDataReturnType is the return type for the SignTypedData action (hex string).
type SignTypedDataReturnType = string

// SignTypedData signs typed data and calculates an Ethereum-specific signature in EIP-712 format:
// sign(keccak256("\x19\x01" ‖ domainSeparator ‖ hashStruct(message)))
//
// - For local accounts (implementing TypedDataSignableAccount), signs locally without an RPC call.
// - For JSON-RPC accounts, delegates to the `eth_signTypedData_v4` RPC method.
//
// This is equivalent to viem's `signTypedData` action.
//
// JSON-RPC Method: eth_signTypedData_v4
//
// Example:
//
//	sig, err := wallet.SignTypedData(ctx, client, wallet.SignTypedDataParameters{
//	    Domain: signature.TypedDataDomain{
//	        Name:              "Ether Mail",
//	        Version:           "1",
//	        ChainId:           big.NewInt(1),
//	        VerifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
//	    },
//	    Types: map[string][]signature.TypedDataField{
//	        "Person": {
//	            {Name: "name", Type: "string"},
//	            {Name: "wallet", Type: "address"},
//	        },
//	        "Mail": {
//	            {Name: "from", Type: "Person"},
//	            {Name: "to", Type: "Person"},
//	            {Name: "contents", Type: "string"},
//	        },
//	    },
//	    PrimaryType: "Mail",
//	    Message: map[string]any{
//	        "from": map[string]any{"name": "Cow", "wallet": "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826"},
//	        "to":   map[string]any{"name": "Bob", "wallet": "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB"},
//	        "contents": "Hello, Bob!",
//	    },
//	})
func SignTypedData(ctx context.Context, client Client, params SignTypedDataParameters) (SignTypedDataReturnType, error) {
	// Resolve account: param > client
	account := params.Account
	if account == nil {
		account = client.Account()
	}
	if account == nil {
		return "", &AccountNotFoundError{DocsPath: "/docs/actions/wallet/signTypedData"}
	}

	// Build the complete types map with EIP712Domain (mirrors viem's getTypesForEIP712Domain)
	types := make(map[string][]signature.TypedDataField)
	for k, v := range params.Types {
		types[k] = v
	}
	types["EIP712Domain"] = getTypesForEIP712Domain(params.Domain)

	// Build the typed data definition for local signing and validation
	typedData := signature.TypedDataDefinition{
		Domain:      params.Domain,
		Types:       types,
		PrimaryType: params.PrimaryType,
		Message:     params.Message,
	}

	// Validate typed data at runtime (mirrors viem's validateTypedData)
	if err := validateTypedData(typedData); err != nil {
		return "", fmt.Errorf("typed data validation failed: %w", err)
	}

	// If the account can sign typed data locally, use it directly
	if signable, ok := account.(TypedDataSignableAccount); ok {
		return signable.SignTypedData(typedData)
	}

	// Otherwise, serialize and send via eth_signTypedData_v4 RPC
	serialized := serializeTypedData(typedData)

	resp, err := client.Request(ctx, "eth_signTypedData_v4", account.Address().Hex(), serialized)
	if err != nil {
		return "", fmt.Errorf("eth_signTypedData_v4 failed: %w", err)
	}

	var hexResult string
	if unmarshalErr := json.Unmarshal(resp.Result, &hexResult); unmarshalErr != nil {
		return "", fmt.Errorf("failed to unmarshal signature: %w", unmarshalErr)
	}

	return hexResult, nil
}

// getTypesForEIP712Domain returns the EIP712Domain type fields based on which
// domain fields are populated. This mirrors viem's getTypesForEIP712Domain.
func getTypesForEIP712Domain(domain signature.TypedDataDomain) []signature.TypedDataField {
	var fields []signature.TypedDataField

	if domain.Name != "" {
		fields = append(fields, signature.TypedDataField{Name: "name", Type: "string"})
	}
	if domain.Version != "" {
		fields = append(fields, signature.TypedDataField{Name: "version", Type: "string"})
	}
	if domain.ChainId != nil {
		fields = append(fields, signature.TypedDataField{Name: "chainId", Type: "uint256"})
	}
	if domain.VerifyingContract != "" {
		fields = append(fields, signature.TypedDataField{Name: "verifyingContract", Type: "address"})
	}
	if domain.Salt != "" {
		fields = append(fields, signature.TypedDataField{Name: "salt", Type: "bytes32"})
	}

	return fields
}

// validateTypedData performs runtime validation on typed data, checking addresses,
// byte ranges, integer ranges, etc. This mirrors viem's validateTypedData.
func validateTypedData(data signature.TypedDataDefinition) error {
	// Validate the domain
	if data.Domain.VerifyingContract != "" {
		if !isValidHexAddress(data.Domain.VerifyingContract) {
			return fmt.Errorf("invalid verifying contract address: %s", data.Domain.VerifyingContract)
		}
	}

	// Validate the primary type exists in types
	if data.PrimaryType != "EIP712Domain" {
		if _, ok := data.Types[data.PrimaryType]; !ok {
			return fmt.Errorf("invalid primary type: %s (not found in types)", data.PrimaryType)
		}
	}

	// Validate message data against type definitions
	if data.PrimaryType != "EIP712Domain" && data.Message != nil {
		if err := validateStruct(data.Types[data.PrimaryType], data.Message, data.Types); err != nil {
			return err
		}
	}

	return nil
}

// validateStruct validates a struct's data against its type definition.
func validateStruct(fields []signature.TypedDataField, data map[string]any, types map[string][]signature.TypedDataField) error {
	for _, field := range fields {
		value, exists := data[field.Name]
		if !exists {
			continue
		}

		// Validate address fields
		if field.Type == "address" {
			if str, ok := value.(string); ok {
				if !isValidHexAddress(str) {
					return fmt.Errorf("invalid address for field %s: %s", field.Name, str)
				}
			}
		}

		// Recursively validate nested custom types
		if _, isCustomType := types[field.Type]; isCustomType {
			if nested, ok := value.(map[string]any); ok {
				if err := validateStruct(types[field.Type], nested, types); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// serializeTypedData serializes typed data for JSON-RPC transmission.
// This mirrors viem's serializeTypedData which normalizes addresses to lowercase
// and produces a JSON string.
func serializeTypedData(data signature.TypedDataDefinition) string {
	// Normalize domain: lowercase addresses
	domain := make(map[string]any)
	if data.Types["EIP712Domain"] != nil && (data.Domain.Name != "" || data.Domain.Version != "" ||
		data.Domain.ChainId != nil || data.Domain.VerifyingContract != "" || data.Domain.Salt != "") {
		domainFields := data.Types["EIP712Domain"]
		domainData := domainToMap(data.Domain)
		for _, f := range domainFields {
			if val, ok := domainData[f.Name]; ok {
				if f.Type == "address" {
					if str, ok := val.(string); ok {
						val = strings.ToLower(str)
					}
				}
				domain[f.Name] = val
			}
		}
	}

	// Normalize message: lowercase addresses
	var message any
	if data.PrimaryType != "EIP712Domain" {
		normalizedMsg := normalizeData(data.Types[data.PrimaryType], data.Message, data.Types)
		message = normalizedMsg
	}

	// Build the serializable structure
	payload := map[string]any{
		"domain":      domain,
		"message":     message,
		"primaryType": data.PrimaryType,
		"types":       data.Types,
	}

	// JSON serialize (mirrors viem's stringify)
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}

// normalizeData normalizes data by lowercasing address fields.
// This mirrors viem's normalizeData in serializeTypedData.
func normalizeData(fields []signature.TypedDataField, data map[string]any, types map[string][]signature.TypedDataField) map[string]any {
	if data == nil {
		return nil
	}
	normalized := make(map[string]any)
	for k, v := range data {
		normalized[k] = v
	}
	for _, field := range fields {
		if field.Type == "address" {
			if str, ok := normalized[field.Name].(string); ok {
				normalized[field.Name] = strings.ToLower(str)
			}
		}
		// Recursively normalize nested custom types
		if nestedFields, isCustom := types[field.Type]; isCustom {
			if nestedData, ok := normalized[field.Name].(map[string]any); ok {
				normalized[field.Name] = normalizeData(nestedFields, nestedData, types)
			}
		}
	}
	return normalized
}

// domainToMap converts a TypedDataDomain to a map.
func domainToMap(domain signature.TypedDataDomain) map[string]any {
	m := make(map[string]any)
	if domain.Name != "" {
		m["name"] = domain.Name
	}
	if domain.Version != "" {
		m["version"] = domain.Version
	}
	if domain.ChainId != nil {
		m["chainId"] = domain.ChainId
	}
	if domain.VerifyingContract != "" {
		m["verifyingContract"] = domain.VerifyingContract
	}
	if domain.Salt != "" {
		m["salt"] = domain.Salt
	}
	return m
}

// isValidHexAddress checks if a string is a valid hex-encoded Ethereum address.
func isValidHexAddress(addr string) bool {
	if len(addr) != 42 {
		return false
	}
	if addr[:2] != "0x" && addr[:2] != "0X" {
		return false
	}
	for _, c := range addr[2:] {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return false
		}
	}
	return true
}
