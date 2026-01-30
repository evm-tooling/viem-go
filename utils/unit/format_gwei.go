package unit

import "math/big"

// GweiDecimals is the number of decimals for gwei (9).
const GweiDecimals = 9

// FormatGwei converts numerical wei to a string representation of gwei.
//
// Example:
//
//	FormatGwei(big.NewInt(1000000000))
//	// "1"
//
//	FormatGwei(big.NewInt(1500000000))
//	// "1.5"
//
//	FormatGwei(big.NewInt(123456789))
//	// "0.123456789"
func FormatGwei(wei *big.Int) string {
	return FormatUnits(wei, GweiDecimals)
}

// FormatGweiInt64 is a convenience function that takes an int64 wei value.
func FormatGweiInt64(wei int64) string {
	return FormatGwei(big.NewInt(wei))
}

// FormatGweiUint64 is a convenience function that takes a uint64 wei value.
func FormatGweiUint64(wei uint64) string {
	return FormatGwei(new(big.Int).SetUint64(wei))
}

// FormatGweiString parses a string wei value and formats it as gwei.
func FormatGweiString(wei string) string {
	v, ok := new(big.Int).SetString(wei, 10)
	if !ok {
		return "0"
	}
	return FormatGwei(v)
}
