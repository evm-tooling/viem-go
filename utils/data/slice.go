package data

import (
	"errors"
	"fmt"
	"strings"
)

// ErrSliceOffsetOutOfBounds is returned when a slice offset is out of bounds.
var ErrSliceOffsetOutOfBounds = errors.New("slice offset out of bounds")

// SliceOptions configures slice behavior.
type SliceOptions struct {
	// Strict validates that the slice produces the expected size.
	Strict bool
}

// Slice returns a section of a byte slice or hex string given start/end byte offsets.
//
// Example:
//
//	result := SliceBytes([]byte{0x01, 0x02, 0x03, 0x04}, 1, 3)
//	// []byte{0x02, 0x03}
//
//	result := SliceHex("0x01020304", 1, 3)
//	// "0x0203"
func Slice(value any, start, end int) any {
	switch v := value.(type) {
	case []byte:
		result, _ := SliceBytes(v, start, end)
		return result
	case string:
		result, _ := SliceHex(v, start, end)
		return result
	default:
		return nil
	}
}

// SliceBytes returns a section of a byte slice given start/end byte offsets.
//
// Example:
//
//	result, _ := SliceBytes([]byte{0x01, 0x02, 0x03, 0x04}, 1, 3)
//	// []byte{0x02, 0x03}
//
//	result, _ := SliceBytes([]byte{0x01, 0x02, 0x03}, 0, 2)
//	// []byte{0x01, 0x02}
func SliceBytes(value []byte, start, end int, opts ...SliceOptions) ([]byte, error) {
	strict := false
	if len(opts) > 0 {
		strict = opts[0].Strict
	}

	// Validate start offset
	if start > 0 && start > len(value)-1 {
		return nil, fmt.Errorf("%w: start offset %d exceeds size %d", ErrSliceOffsetOutOfBounds, start, len(value))
	}

	// Handle negative or default values
	if start < 0 {
		start = 0
	}
	if end <= 0 || end > len(value) {
		end = len(value)
	}

	result := value[start:end]

	// Strict mode validates the result size
	if strict && len(result) != end-start {
		return nil, fmt.Errorf("%w: end offset %d produces unexpected size", ErrSliceOffsetOutOfBounds, end)
	}

	return result, nil
}

// SliceHex returns a section of a hex string given start/end byte offsets.
//
// Example:
//
//	result, _ := SliceHex("0x01020304", 1, 3)
//	// "0x0203"
//
//	result, _ := SliceHex("0x01020304", 0, 2)
//	// "0x0102"
func SliceHex(value string, start, end int, opts ...SliceOptions) (string, error) {
	strict := false
	if len(opts) > 0 {
		strict = opts[0].Strict
	}

	// Remove 0x prefix
	h := strings.TrimPrefix(value, "0x")
	h = strings.TrimPrefix(h, "0X")

	byteLen := len(h) / 2

	// Validate start offset
	if start > 0 && start > byteLen-1 {
		return "", fmt.Errorf("%w: start offset %d exceeds size %d", ErrSliceOffsetOutOfBounds, start, byteLen)
	}

	// Handle negative or default values
	if start < 0 {
		start = 0
	}
	if end <= 0 || end > byteLen {
		end = byteLen
	}

	// Convert byte offsets to hex character offsets
	result := "0x" + h[start*2:end*2]

	// Strict mode validates the result size
	if strict {
		resultSize := SizeHex(result)
		if resultSize != end-start {
			return "", fmt.Errorf("%w: end offset %d produces unexpected size", ErrSliceOffsetOutOfBounds, end)
		}
	}

	return result, nil
}

// SliceBytesStart returns a slice from start to the end of the value.
func SliceBytesStart(value []byte, start int) ([]byte, error) {
	return SliceBytes(value, start, len(value))
}

// SliceHexStart returns a slice from start to the end of the hex string.
func SliceHexStart(value string, start int) (string, error) {
	h := strings.TrimPrefix(value, "0x")
	h = strings.TrimPrefix(h, "0X")
	return SliceHex(value, start, len(h)/2)
}
