package utils

import "regexp"

// ABI type matching regular expressions.
// These are used for parsing and validating Solidity ABI types.

// ArrayRegex matches array types like "uint256[]" or "bytes32[10]"
// Captures: group 1 = base type, group 2 = array size (empty for dynamic arrays)
//
// Examples:
//   - "uint256[]" -> ["uint256[]", "uint256", ""]
//   - "bytes32[10]" -> ["bytes32[10]", "bytes32", "10"]
var ArrayRegex = regexp.MustCompile(`^(.*)\[([0-9]*)\]$`)

// BytesRegex matches fixed-size bytes types like "bytes1" through "bytes32"
// Also matches "bytes" (dynamic bytes)
//
// Examples:
//   - "bytes" -> matches
//   - "bytes1" -> matches
//   - "bytes32" -> matches
//   - "bytes33" -> does not match
var BytesRegex = regexp.MustCompile(`^bytes([1-9]|1[0-9]|2[0-9]|3[0-2])?$`)

// IntegerRegex matches integer types like "uint256", "int128", "uint", "int"
// Captures: group 1 = "uint" or "int", group 2 = bit size (optional)
//
// Valid bit sizes: 8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96,
// 104, 112, 120, 128, 136, 144, 152, 160, 168, 176, 184, 192,
// 200, 208, 216, 224, 232, 240, 248, 256
//
// Examples:
//   - "uint256" -> ["uint256", "uint", "256"]
//   - "int" -> ["int", "int", ""]
//   - "uint8" -> ["uint8", "uint", "8"]
var IntegerRegex = regexp.MustCompile(`^(u?int)(8|16|24|32|40|48|56|64|72|80|88|96|104|112|120|128|136|144|152|160|168|176|184|192|200|208|216|224|232|240|248|256)?$`)

// IsArrayType checks if a type string represents an array type.
func IsArrayType(typ string) bool {
	return ArrayRegex.MatchString(typ)
}

// IsBytesType checks if a type string represents a bytes type.
func IsBytesType(typ string) bool {
	return BytesRegex.MatchString(typ)
}

// IsIntegerType checks if a type string represents an integer type.
func IsIntegerType(typ string) bool {
	return IntegerRegex.MatchString(typ)
}

// ParseArrayType parses an array type and returns the base type and size.
// Returns empty strings if not an array type.
// Size is empty string for dynamic arrays.
func ParseArrayType(typ string) (baseType string, size string) {
	matches := ArrayRegex.FindStringSubmatch(typ)
	if len(matches) < 3 {
		return "", ""
	}
	return matches[1], matches[2]
}

// ParseIntegerType parses an integer type and returns whether it's unsigned and the bit size.
// Returns (false, 0) if not an integer type.
// Default bit size is 256 if not specified.
func ParseIntegerType(typ string) (unsigned bool, bitSize int) {
	matches := IntegerRegex.FindStringSubmatch(typ)
	if len(matches) < 2 {
		return false, 0
	}

	unsigned = matches[1] == "uint"

	if len(matches) > 2 && matches[2] != "" {
		// Parse bit size
		size := 0
		for _, c := range matches[2] {
			size = size*10 + int(c-'0')
		}
		return unsigned, size
	}

	// Default to 256 bits
	return unsigned, 256
}
