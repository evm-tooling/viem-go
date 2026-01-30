package data

import (
	"regexp"
	"strings"
)

var strictHexRegex = regexp.MustCompile(`^0x[0-9a-fA-F]*$`)

// IsHexOptions configures hex validation behavior.
type IsHexOptions struct {
	// Strict validates that the string contains only valid hex characters.
	// Default is true.
	Strict bool
}

// IsHex checks if a value is a valid hex string.
// By default (strict=true), validates that the string contains only valid hex characters.
// With strict=false, only checks for 0x prefix.
//
// Example:
//
//	IsHex("0x0102")                              // true
//	IsHex("0x")                                  // true (empty hex)
//	IsHex("0xgg", IsHexOptions{Strict: true})   // false (invalid chars)
//	IsHex("0xgg", IsHexOptions{Strict: false})  // true (has prefix)
//	IsHex("0102")                               // false (no prefix)
func IsHex(value string, opts ...IsHexOptions) bool {
	if value == "" {
		return false
	}

	strict := true
	if len(opts) > 0 {
		strict = opts[0].Strict
	}

	if strict {
		return strictHexRegex.MatchString(value)
	}

	return strings.HasPrefix(value, "0x") || strings.HasPrefix(value, "0X")
}

// IsHexString is an alias for IsHex with strict=true.
func IsHexString(value string) bool {
	return IsHex(value, IsHexOptions{Strict: true})
}
