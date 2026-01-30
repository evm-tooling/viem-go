package unit

import (
	"errors"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

var decimalRegex = regexp.MustCompile(`^(-?)([0-9]*)\.?([0-9]*)$`)

// ErrInvalidDecimalNumber is returned when the value is not a valid decimal number.
var ErrInvalidDecimalNumber = errors.New("invalid decimal number")

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
	if !decimalRegex.MatchString(value) {
		return nil, ErrInvalidDecimalNumber
	}

	parts := strings.Split(value, ".")
	integer := parts[0]
	fraction := "0"
	if len(parts) > 1 {
		fraction = parts[1]
	}

	negative := strings.HasPrefix(integer, "-")
	if negative {
		integer = integer[1:]
	}

	// Trim trailing zeros from fraction
	fraction = strings.TrimRight(fraction, "0")
	if fraction == "" {
		fraction = "0"
	}

	// Handle rounding when fraction is larger than decimals
	if decimals == 0 {
		// Round the fraction
		if fraction != "0" {
			fracFloat, _ := strconv.ParseFloat("0."+fraction, 64)
			if fracFloat >= 0.5 {
				intVal, _ := new(big.Int).SetString(integer, 10)
				if intVal == nil {
					intVal = big.NewInt(0)
				}
				intVal.Add(intVal, big.NewInt(1))
				integer = intVal.String()
			}
		}
		fraction = ""
	} else if len(fraction) > decimals {
		// Round off if fraction is larger than decimals
		left := fraction[:decimals-1]
		unit := fraction[decimals-1 : decimals]
		right := fraction[decimals:]

		// Calculate rounded value
		roundVal, _ := strconv.ParseFloat(unit+"."+right, 64)
		rounded := int(roundVal + 0.5)

		if rounded > 9 {
			// Carry over
			leftVal, _ := new(big.Int).SetString(left, 10)
			if leftVal == nil {
				leftVal = big.NewInt(0)
			}
			leftVal.Add(leftVal, big.NewInt(1))
			fraction = leftVal.String() + "0"
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
			intVal, _ := new(big.Int).SetString(integer, 10)
			if intVal == nil {
				intVal = big.NewInt(0)
			}
			intVal.Add(intVal, big.NewInt(1))
			integer = intVal.String()
		}

		fraction = fraction[:decimals]
	} else {
		// Pad fraction with trailing zeros
		for len(fraction) < decimals {
			fraction = fraction + "0"
		}
	}

	// Handle empty integer
	if integer == "" {
		integer = "0"
	}

	// Combine integer and fraction
	combined := integer + fraction
	result, ok := new(big.Int).SetString(combined, 10)
	if !ok {
		return nil, ErrInvalidDecimalNumber
	}

	if negative {
		result.Neg(result)
	}

	return result, nil
}

// MustParseUnits is like ParseUnits but panics on error.
func MustParseUnits(value string, decimals int) *big.Int {
	result, err := ParseUnits(value, decimals)
	if err != nil {
		panic(err)
	}
	return result
}
