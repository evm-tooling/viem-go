package data

import (
	"errors"
	"fmt"
	"strings"
)

// ErrSizeExceedsPaddingSize is returned when the value exceeds the target padding size.
var ErrSizeExceedsPaddingSize = errors.New("size exceeds padding size")

// PadDirection specifies the padding direction.
type PadDirection string

const (
	PadLeft  PadDirection = "left"
	PadRight PadDirection = "right"
)

// PadOptions configures padding behavior.
type PadOptions struct {
	// Dir specifies the padding direction. Default is "left".
	Dir PadDirection
	// Size is the target size in bytes. Default is 32.
	Size int
}

// Pad pads a byte slice to the specified size.
// Default padding is on the left to 32 bytes.
//
// Example:
//
//	result, _ := Pad([]byte{0x01})
//	// []byte{0x00, 0x00, ..., 0x01} (32 bytes, left-padded)
//
//	result, _ := Pad([]byte{0x01}, PadOptions{Dir: PadRight, Size: 4})
//	// []byte{0x01, 0x00, 0x00, 0x00}
func Pad(bytes []byte, opts ...PadOptions) ([]byte, error) {
	dir := PadLeft
	size := 32
	if len(opts) > 0 {
		if opts[0].Dir != "" {
			dir = opts[0].Dir
		}
		if opts[0].Size > 0 {
			size = opts[0].Size
		}
	}

	return PadBytes(bytes, dir, size)
}

// PadBytes pads a byte slice to the specified size in the given direction.
//
// Example:
//
//	result, _ := PadBytes([]byte{0x01, 0x02}, PadLeft, 4)
//	// []byte{0x00, 0x00, 0x01, 0x02}
//
//	result, _ := PadBytes([]byte{0x01, 0x02}, PadRight, 4)
//	// []byte{0x01, 0x02, 0x00, 0x00}
func PadBytes(bytes []byte, dir PadDirection, size int) ([]byte, error) {
	if len(bytes) > size {
		return nil, fmt.Errorf("%w: size %d exceeds target size %d", ErrSizeExceedsPaddingSize, len(bytes), size)
	}

	result := make([]byte, size)

	if dir == PadRight {
		// Pad on right: [data, zeros]
		copy(result, bytes)
	} else {
		// Pad on left: [zeros, data]
		copy(result[size-len(bytes):], bytes)
	}

	return result, nil
}

// PadHex pads a hex string to the specified byte size.
//
// Example:
//
//	result, _ := PadHex("0x01", PadLeft, 4)
//	// "0x00000001"
//
//	result, _ := PadHex("0x01", PadRight, 4)
//	// "0x01000000"
func PadHex(hex string, dir PadDirection, size int) (string, error) {
	// Remove 0x prefix
	h := strings.TrimPrefix(hex, "0x")
	h = strings.TrimPrefix(h, "0X")

	// Check size (each byte is 2 hex chars)
	if len(h) > size*2 {
		return "", fmt.Errorf("%w: size %d exceeds target size %d", ErrSizeExceedsPaddingSize, (len(h)+1)/2, size)
	}

	targetLen := size * 2
	if dir == PadRight {
		// Pad on right
		for len(h) < targetLen {
			h = h + "0"
		}
	} else {
		// Pad on left
		for len(h) < targetLen {
			h = "0" + h
		}
	}

	return "0x" + h, nil
}

// PadLeftBytes is a convenience function that pads bytes on the left.
func PadLeftBytes(bytes []byte, size int) ([]byte, error) {
	return PadBytes(bytes, PadLeft, size)
}

// PadLeftHex is a convenience function that pads hex on the left.
func PadLeftHex(hex string, size int) (string, error) {
	return PadHex(hex, PadLeft, size)
}

// PadRightBytes is a convenience function that pads bytes on the right.
func PadRightBytes(bytes []byte, size int) ([]byte, error) {
	return PadBytes(bytes, PadRight, size)
}

// PadRightHex is a convenience function that pads hex on the right.
func PadRightHex(hex string, size int) (string, error) {
	return PadHex(hex, PadRight, size)
}
