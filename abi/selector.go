package abi

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// ComputeSelector computes the 4-byte function selector from a function signature.
// Example: ComputeSelector("transfer(address,uint256)") returns [4]byte{0xa9, 0x05, 0x9c, 0xbb}
func ComputeSelector(signature string) [4]byte {
	hash := crypto.Keccak256([]byte(signature))
	var selector [4]byte
	copy(selector[:], hash[:4])
	return selector
}

// ComputeSelectorHex computes the function selector and returns it as a hex string with 0x prefix.
func ComputeSelectorHex(signature string) string {
	selector := ComputeSelector(signature)
	return "0x" + hex.EncodeToString(selector[:])
}

// ComputeEventTopic computes the 32-byte topic hash for an event signature.
// Example: ComputeEventTopic("Transfer(address,address,uint256)") returns the keccak256 hash
func ComputeEventTopic(signature string) common.Hash {
	return crypto.Keccak256Hash([]byte(signature))
}

// ComputeEventTopicHex computes the event topic and returns it as a hex string with 0x prefix.
func ComputeEventTopicHex(signature string) string {
	topic := ComputeEventTopic(signature)
	return topic.Hex()
}

// SelectorToHex converts a 4-byte selector to a hex string with 0x prefix.
func SelectorToHex(selector [4]byte) string {
	return "0x" + hex.EncodeToString(selector[:])
}

// HexToSelector converts a hex string to a 4-byte selector.
// Accepts both "0x" prefixed and unprefixed strings.
func HexToSelector(hexStr string) ([4]byte, error) {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	hexStr = strings.TrimPrefix(hexStr, "0X")

	if len(hexStr) != 8 {
		return [4]byte{}, fmt.Errorf("invalid selector hex length: expected 8 characters, got %d", len(hexStr))
	}

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return [4]byte{}, fmt.Errorf("invalid hex string: %w", err)
	}

	var selector [4]byte
	copy(selector[:], bytes)
	return selector, nil
}

// MustHexToSelector converts a hex string to a selector, panicking on error.
func MustHexToSelector(hexStr string) [4]byte {
	selector, err := HexToSelector(hexStr)
	if err != nil {
		panic(err)
	}
	return selector
}

// BuildFunctionSignature builds a function signature from name and parameter types.
// Example: BuildFunctionSignature("transfer", []string{"address", "uint256"}) returns "transfer(address,uint256)"
func BuildFunctionSignature(name string, paramTypes []string) string {
	return fmt.Sprintf("%s(%s)", name, strings.Join(paramTypes, ","))
}

// BuildEventSignature builds an event signature from name and parameter types.
// This is the same as BuildFunctionSignature but provided for semantic clarity.
func BuildEventSignature(name string, paramTypes []string) string {
	return BuildFunctionSignature(name, paramTypes)
}

// ParseFunctionSignature parses a function signature and returns the name and parameter types.
// Example: ParseFunctionSignature("transfer(address,uint256)") returns ("transfer", []string{"address", "uint256"})
func ParseFunctionSignature(signature string) (string, []string, error) {
	// Find the opening parenthesis
	parenIdx := strings.Index(signature, "(")
	if parenIdx == -1 {
		return "", nil, fmt.Errorf("invalid signature: missing opening parenthesis")
	}

	// Check for closing parenthesis
	if !strings.HasSuffix(signature, ")") {
		return "", nil, fmt.Errorf("invalid signature: missing closing parenthesis")
	}

	name := signature[:parenIdx]
	if name == "" {
		return "", nil, fmt.Errorf("invalid signature: empty function name")
	}

	// Extract parameter types
	paramsStr := signature[parenIdx+1 : len(signature)-1]
	if paramsStr == "" {
		return name, nil, nil
	}

	// Handle nested parentheses (tuples)
	params := parseParameterTypes(paramsStr)
	return name, params, nil
}

// parseParameterTypes splits parameter types, handling nested parentheses for tuples.
func parseParameterTypes(paramsStr string) []string {
	if paramsStr == "" {
		return nil
	}

	var params []string
	var current strings.Builder
	depth := 0

	for _, c := range paramsStr {
		switch c {
		case '(':
			depth++
			current.WriteRune(c)
		case ')':
			depth--
			current.WriteRune(c)
		case ',':
			if depth == 0 {
				params = append(params, strings.TrimSpace(current.String()))
				current.Reset()
			} else {
				current.WriteRune(c)
			}
		default:
			current.WriteRune(c)
		}
	}

	if current.Len() > 0 {
		params = append(params, strings.TrimSpace(current.String()))
	}

	return params
}

// StandardSelectors contains common ERC function selectors for quick lookup.
var StandardSelectors = map[string][4]byte{
	// ERC20
	"name":         ComputeSelector("name()"),
	"symbol":       ComputeSelector("symbol()"),
	"decimals":     ComputeSelector("decimals()"),
	"totalSupply":  ComputeSelector("totalSupply()"),
	"balanceOf":    ComputeSelector("balanceOf(address)"),
	"transfer":     ComputeSelector("transfer(address,uint256)"),
	"transferFrom": ComputeSelector("transferFrom(address,address,uint256)"),
	"approve":      ComputeSelector("approve(address,uint256)"),
	"allowance":    ComputeSelector("allowance(address,address)"),

	// ERC721
	"ownerOf":           ComputeSelector("ownerOf(uint256)"),
	"safeTransferFrom":  ComputeSelector("safeTransferFrom(address,address,uint256)"),
	"getApproved":       ComputeSelector("getApproved(uint256)"),
	"setApprovalForAll": ComputeSelector("setApprovalForAll(address,bool)"),
	"isApprovedForAll":  ComputeSelector("isApprovedForAll(address,address)"),
}

// StandardEventTopics contains common ERC event topics for quick lookup.
var StandardEventTopics = map[string]common.Hash{
	// ERC20/ERC721 Transfer
	"Transfer": ComputeEventTopic("Transfer(address,address,uint256)"),
	// ERC20 Approval
	"Approval": ComputeEventTopic("Approval(address,address,uint256)"),
	// ERC721 ApprovalForAll
	"ApprovalForAll": ComputeEventTopic("ApprovalForAll(address,address,bool)"),
}

// IsStandardSelector checks if a selector matches a known standard function.
func IsStandardSelector(selector [4]byte) (string, bool) {
	for name, stdSelector := range StandardSelectors {
		if selector == stdSelector {
			return name, true
		}
	}
	return "", false
}

// IsStandardEventTopic checks if a topic matches a known standard event.
func IsStandardEventTopic(topic common.Hash) (string, bool) {
	for name, stdTopic := range StandardEventTopics {
		if topic == stdTopic {
			return name, true
		}
	}
	return "", false
}
