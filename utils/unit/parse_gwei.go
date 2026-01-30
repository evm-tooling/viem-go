package unit

import "math/big"

// ParseGwei converts a string representation of gwei to numerical wei.
//
// Example:
//
//	ParseGwei("420")
//	// big.Int representing 420000000000
//
//	ParseGwei("1.5")
//	// big.Int representing 1500000000
//
//	ParseGwei("0.1")
//	// big.Int representing 100000000
func ParseGwei(gwei string) (*big.Int, error) {
	return ParseUnits(gwei, GweiDecimals)
}

// MustParseGwei is like ParseGwei but panics on error.
func MustParseGwei(gwei string) *big.Int {
	result, err := ParseGwei(gwei)
	if err != nil {
		panic(err)
	}
	return result
}
