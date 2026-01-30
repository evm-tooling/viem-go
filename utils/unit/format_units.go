package unit

import (
	"math/big"
	"strings"
)

// FormatUnits divides a number by a given exponent of base 10 (10^decimals),
// and formats it into a string representation of the number.
//
// Example:
//
//	FormatUnits(big.NewInt(420000000000), 9)
//	// "420"
//
//	FormatUnits(big.NewInt(1000000000000000000), 18)
//	// "1"
//
//	FormatUnits(big.NewInt(123456789), 6)
//	// "123.456789"
func FormatUnits(value *big.Int, decimals int) string {
	if value == nil {
		return "0"
	}

	display := value.String()

	negative := strings.HasPrefix(display, "-")
	if negative {
		display = display[1:]
	}

	// Pad with leading zeros if necessary
	for len(display) < decimals {
		display = "0" + display
	}

	// Split into integer and fraction parts
	splitPoint := len(display) - decimals
	integer := display[:splitPoint]
	fraction := display[splitPoint:]

	// Remove trailing zeros from fraction
	fraction = strings.TrimRight(fraction, "0")

	// Build result
	if integer == "" {
		integer = "0"
	}

	result := ""
	if negative {
		result = "-"
	}
	result += integer

	if fraction != "" {
		result += "." + fraction
	}

	return result
}

// FormatUnitsInt64 is a convenience function that takes an int64 value.
func FormatUnitsInt64(value int64, decimals int) string {
	return FormatUnits(big.NewInt(value), decimals)
}

// FormatUnitsUint64 is a convenience function that takes a uint64 value.
func FormatUnitsUint64(value uint64, decimals int) string {
	return FormatUnits(new(big.Int).SetUint64(value), decimals)
}

// FormatUnitsString parses a string value and formats it.
func FormatUnitsString(value string, decimals int) string {
	v, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return "0"
	}
	return FormatUnits(v, decimals)
}
