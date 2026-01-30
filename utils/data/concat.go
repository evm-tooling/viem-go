package data

import (
	"strings"
)

// Concat concatenates multiple byte slices or hex strings.
// Automatically detects the type based on the first element.
//
// Example:
//
//	// Bytes
//	result := Concat([]byte{0x01}, []byte{0x02})
//	// []byte{0x01, 0x02}
//
//	// Hex strings
//	result := ConcatHex("0x01", "0x02")
//	// "0x0102"
func Concat(values ...[]byte) []byte {
	return ConcatBytes(values...)
}

// ConcatBytes concatenates multiple byte slices into one.
//
// Example:
//
//	result := ConcatBytes([]byte{0x01, 0x02}, []byte{0x03, 0x04})
//	// []byte{0x01, 0x02, 0x03, 0x04}
func ConcatBytes(values ...[]byte) []byte {
	// Calculate total length
	length := 0
	for _, arr := range values {
		length += len(arr)
	}

	// Create result array and copy values
	result := make([]byte, length)
	offset := 0
	for _, arr := range values {
		copy(result[offset:], arr)
		offset += len(arr)
	}

	return result
}

// ConcatHex concatenates multiple hex strings into one.
// All inputs should have 0x prefix.
//
// Example:
//
//	result := ConcatHex("0x0102", "0x0304")
//	// "0x01020304"
func ConcatHex(values ...string) string {
	var builder strings.Builder
	builder.WriteString("0x")

	for _, hex := range values {
		// Remove 0x prefix
		h := strings.TrimPrefix(hex, "0x")
		h = strings.TrimPrefix(h, "0X")
		builder.WriteString(h)
	}

	return builder.String()
}
