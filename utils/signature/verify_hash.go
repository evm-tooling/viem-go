package signature

import (
	"strings"
)

// VerifyHash verifies that a hash was signed by the provided address.
//
// Note: Only supports Externally Owned Accounts. Does not support Contract Accounts.
//
// Example:
//
//	valid, err := VerifyHash(
//		"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
//		"0xd9eba16ed0ecae432b71fe008c98cc872bb4cc214d3220a36f365326cf807d68",
//		"0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db81c",
//	)
//	// true
func VerifyHash(address, hash string, signature any) (bool, error) {
	recoveredAddress, err := RecoverAddress(hash, signature)
	if err != nil {
		return false, err
	}

	return isAddressEqual(address, recoveredAddress), nil
}

// isAddressEqual compares two addresses case-insensitively.
func isAddressEqual(a, b string) bool {
	return strings.EqualFold(a, b)
}
