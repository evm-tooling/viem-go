package data

import "strings"

// Size returns the size of a hex string or byte slice in bytes.
//
// For hex strings, it calculates the byte length from the hex characters.
// For byte slices, it returns the length directly.
//
// Example:
//
//	Size("0x0102")       // 2
//	Size("0x")           // 0
//	Size([]byte{1, 2})   // 2
func Size(value any) int {
	switch v := value.(type) {
	case string:
		return SizeHex(v)
	case []byte:
		return len(v)
	default:
		return 0
	}
}

// SizeHex returns the byte size of a hex string.
// The hex string should have a 0x prefix.
//
// Example:
//
//	SizeHex("0x0102")   // 2
//	SizeHex("0x")       // 0
//	SizeHex("0x1")      // 1 (odd length, rounds up)
func SizeHex(hex string) int {
	// Remove 0x prefix
	h := strings.TrimPrefix(hex, "0x")
	h = strings.TrimPrefix(h, "0X")

	// Each byte is 2 hex characters, round up for odd lengths
	return (len(h) + 1) / 2
}

// SizeBytes returns the byte size of a byte slice.
// This is simply the length of the slice.
//
// Example:
//
//	SizeBytes([]byte{1, 2, 3})  // 3
func SizeBytes(bytes []byte) int {
	return len(bytes)
}
