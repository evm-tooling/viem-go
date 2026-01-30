package address

import (
	"fmt"
	"strings"
)

// IsAddressEqual compares two addresses for equality (case-insensitive).
// Returns an error if either address is invalid.
//
// Example:
//
//	isAddressEqual(
//	  "0xa5cc3c03994db5b0d9a5eedd10cabab0813678ac",
//	  "0xa5cC3c03994DB5b0d9A5EEdD10CabaB0813678AC"
//	) // true
func IsAddressEqual(a, b string) (bool, error) {
	if !IsAddress(a, IsAddressOptions{Strict: false}) {
		return false, fmt.Errorf("%w: %s", ErrInvalidAddress, a)
	}
	if !IsAddress(b, IsAddressOptions{Strict: false}) {
		return false, fmt.Errorf("%w: %s", ErrInvalidAddress, b)
	}
	return strings.EqualFold(a, b), nil
}
