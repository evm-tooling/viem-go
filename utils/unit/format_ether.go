package unit

import "math/big"

// EtherDecimals is the number of decimals for ether (18).
const EtherDecimals = 18

// GweiToEtherDecimals is the number of decimals when converting gwei to ether (9).
const GweiToEtherDecimals = 9

// FormatEther converts numerical wei to a string representation of ether.
//
// Example:
//
//	FormatEther(big.NewInt(1000000000000000000))
//	// "1"
//
//	FormatEther(big.NewInt(1500000000000000000))
//	// "1.5"
//
//	FormatEther(big.NewInt(123456789000000000))
//	// "0.123456789"
func FormatEther(wei *big.Int) string {
	return FormatUnits(wei, EtherDecimals)
}

// FormatEtherFromGwei converts numerical gwei to a string representation of ether.
//
// Example:
//
//	FormatEtherFromGwei(big.NewInt(1000000000))
//	// "1"
func FormatEtherFromGwei(gwei *big.Int) string {
	return FormatUnits(gwei, GweiToEtherDecimals)
}

// FormatEtherInt64 is a convenience function that takes an int64 wei value.
func FormatEtherInt64(wei int64) string {
	return FormatEther(big.NewInt(wei))
}

// FormatEtherUint64 is a convenience function that takes a uint64 wei value.
func FormatEtherUint64(wei uint64) string {
	return FormatEther(new(big.Int).SetUint64(wei))
}

// FormatEtherString parses a string wei value and formats it as ether.
func FormatEtherString(wei string) string {
	v, ok := new(big.Int).SetString(wei, 10)
	if !ok {
		return "0"
	}
	return FormatEther(v)
}
