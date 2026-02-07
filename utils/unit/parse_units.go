package unit

import (
	"errors"
	"math/big"
	"strconv"
	"strings"
)

// ErrInvalidDecimalNumber is returned when the value is not a valid decimal number.
var ErrInvalidDecimalNumber = errors.New("invalid decimal number")

// Cached powers of 10 for common decimal values.
// Avoids repeated big.Int allocations for the most common units (6, 8, 9, 18).
var powersOf10 [78]*big.Int

func init() {
	powersOf10[0] = big.NewInt(1)
	for i := 1; i < len(powersOf10); i++ {
		powersOf10[i] = new(big.Int).Mul(powersOf10[i-1], big.NewInt(10))
	}
}

// isValidDecimal checks if a string is a valid decimal number without using regex.
// Matches the pattern: ^(-?)([0-9]*)\.?([0-9]*)$ with at least one digit.
func isValidDecimal(s string) bool {
	if len(s) == 0 {
		return false
	}

	i := 0
	if s[0] == '-' {
		i++
	}

	hasDot := false
	hasDigit := false

	for ; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= '0' && c <= '9':
			hasDigit = true
		case c == '.' && !hasDot:
			hasDot = true
		default:
			return false
		}
	}

	// Must have at least one digit, or just "." alone is invalid
	return hasDigit
}

// ParseUnits multiplies a string representation of a number by a given exponent
// of base 10 (10^decimals).
//
// Example:
//
//	ParseUnits("420", 9)
//	// big.Int representing 420000000000
//
//	ParseUnits("1", 18)
//	// big.Int representing 1000000000000000000
//
//	ParseUnits("1.5", 18)
//	// big.Int representing 1500000000000000000
func ParseUnits(value string, decimals int) (*big.Int, error) {
	if !isValidDecimal(value) {
		return nil, ErrInvalidDecimalNumber
	}

	// Split into integer and fraction parts
	integer, fraction := splitDecimal(value)

	negative := false
	if len(integer) > 0 && integer[0] == '-' {
		negative = true
		integer = integer[1:]
	}

	// Trim trailing zeros from fraction
	fraction = trimRight(fraction, '0')

	// Handle rounding when fraction is larger than decimals
	if decimals == 0 {
		// Round the fraction
		if len(fraction) > 0 {
			// Check if first digit >= 5
			if fraction[0] >= '5' {
				integer = incrementString(integer)
			}
		}
		fraction = ""
	} else if len(fraction) > decimals {
		// Round off if fraction is larger than decimals
		left := fraction[:decimals-1]
		unit := fraction[decimals-1]
		right := fraction[decimals:]

		// Calculate rounded value
		roundVal, _ := strconv.ParseFloat(string(unit)+"."+right, 64)
		rounded := int(roundVal + 0.5)

		if rounded > 9 {
			// Carry over
			fraction = incrementString(left) + "0"
			// Pad to ensure proper length
			for len(fraction) < len(left)+1 {
				fraction = "0" + fraction
			}
		} else {
			fraction = left + strconv.Itoa(rounded)
		}

		// Check if we need to carry to integer
		if len(fraction) > decimals {
			fraction = fraction[1:]
			integer = incrementString(integer)
		}

		fraction = fraction[:decimals]
	} else {
		// Pad fraction with trailing zeros using a single allocation
		if pad := decimals - len(fraction); pad > 0 {
			fraction = fraction + strings.Repeat("0", pad)
		}
	}

	// Handle empty integer
	if integer == "" {
		integer = "0"
	}

	// Combine integer and fraction, then parse
	combined := integer + fraction

	// Fast path: if the combined value fits in uint64, avoid big.Int.SetString
	if !negative && len(combined) <= 19 {
		val, err := strconv.ParseUint(combined, 10, 64)
		if err == nil {
			return new(big.Int).SetUint64(val), nil
		}
	}

	result, ok := new(big.Int).SetString(combined, 10)
	if !ok {
		return nil, ErrInvalidDecimalNumber
	}

	if negative {
		result.Neg(result)
	}

	return result, nil
}

// splitDecimal splits a decimal string into integer and fraction parts.
// Avoids strings.Split which allocates a slice.
func splitDecimal(s string) (integer, fraction string) {
	dot := strings.IndexByte(s, '.')
	if dot < 0 {
		return s, ""
	}
	return s[:dot], s[dot+1:]
}

// trimRight trims trailing bytes from a string without allocating if nothing to trim.
func trimRight(s string, b byte) string {
	i := len(s)
	for i > 0 && s[i-1] == b {
		i--
	}
	return s[:i]
}

// incrementString increments a decimal number string by 1.
// Avoids big.Int allocation for simple carry operations.
func incrementString(s string) string {
	if s == "" || s == "0" {
		return "1"
	}

	// Fast path: try uint64
	if len(s) <= 18 {
		val, err := strconv.ParseUint(s, 10, 64)
		if err == nil {
			return strconv.FormatUint(val+1, 10)
		}
	}

	// Fallback to big.Int for huge numbers
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return "1"
	}
	v.Add(v, big.NewInt(1))
	return v.String()
}

// MustParseUnits is like ParseUnits but panics on error.
func MustParseUnits(value string, decimals int) *big.Int {
	result, err := ParseUnits(value, decimals)
	if err != nil {
		panic(err)
	}
	return result
}
