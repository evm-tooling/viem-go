package unit

import "math/big"

// ParseEther converts a string representation of ether to numerical wei.
//
// Example:
//
//	ParseEther("420")
//	// big.Int representing 420000000000000000000
//
//	ParseEther("1.5")
//	// big.Int representing 1500000000000000000
//
//	ParseEther("0.1")
//	// big.Int representing 100000000000000000
func ParseEther(ether string) (*big.Int, error) {
	return ParseUnits(ether, EtherDecimals)
}

// MustParseEther is like ParseEther but panics on error.
func MustParseEther(ether string) *big.Int {
	result, err := ParseEther(ether)
	if err != nil {
		panic(err)
	}
	return result
}

// ParseEtherToGwei converts a string representation of ether to numerical gwei.
//
// Example:
//
//	ParseEtherToGwei("1")
//	// big.Int representing 1000000000
func ParseEtherToGwei(ether string) (*big.Int, error) {
	return ParseUnits(ether, GweiToEtherDecimals)
}

// MustParseEtherToGwei is like ParseEtherToGwei but panics on error.
func MustParseEtherToGwei(ether string) *big.Int {
	result, err := ParseEtherToGwei(ether)
	if err != nil {
		panic(err)
	}
	return result
}
