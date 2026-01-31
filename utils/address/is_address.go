package address

import (
	"regexp"
	"strings"
)

var (
	// addressRegex matches a valid Ethereum address format
	addressRegex = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
)

// IsAddressOptions configures address validation behavior.
type IsAddressOptions struct {
	// Strict enables checksum validation. Default is true.
	Strict bool
}

// IsAddress checks if a string is a valid Ethereum address.
// By default (strict=true), it validates the checksum if the address contains mixed case.
//
// Example:
//
//	isAddress("0xa5cc3c03994db5b0d9a5eedd10cabab0813678ac") // true (lowercase)
//	isAddress("0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC") // true (valid checksum)
//	isAddress("0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC", IsAddressOptions{Strict: false}) // true
func IsAddress(address string, opts ...IsAddressOptions) bool {
	strict := true
	if len(opts) > 0 {
		strict = opts[0].Strict
	}

	// Check basic format
	if !addressRegex.MatchString(address) {
		return false
	}

	// If all lowercase, it's valid (no checksum to verify)
	if strings.ToLower(address) == address {
		return true
	}

	// If strict mode and contains uppercase, verify checksum
	if strict {
		checksummed := ChecksumAddress(address)
		return checksummed == Address(address)
	}

	return true
}
