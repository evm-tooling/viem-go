package data

import "strings"

// TrimDirection specifies the trim direction.
type TrimDirection string

const (
	TrimLeft  TrimDirection = "left"
	TrimRight TrimDirection = "right"
)

// TrimOptions configures trim behavior.
type TrimOptions struct {
	// Dir specifies the trim direction. Default is "left".
	Dir TrimDirection
}

// Trim removes leading or trailing zero bytes/characters from a byte slice or hex string.
//
// Example:
//
//	result := TrimBytes([]byte{0x00, 0x00, 0x01, 0x02})
//	// []byte{0x01, 0x02}
//
//	result := TrimHex("0x00000102")
//	// "0x0102"
func Trim(value any, opts ...TrimOptions) any {
	dir := TrimLeft
	if len(opts) > 0 && opts[0].Dir != "" {
		dir = opts[0].Dir
	}

	switch v := value.(type) {
	case []byte:
		return TrimBytes(v, dir)
	case string:
		return TrimHex(v, dir)
	default:
		return value
	}
}

// TrimBytes removes leading or trailing zero bytes from a byte slice.
//
// Example:
//
//	result := TrimBytes([]byte{0x00, 0x00, 0x01, 0x02}, TrimLeft)
//	// []byte{0x01, 0x02}
//
//	result := TrimBytes([]byte{0x01, 0x02, 0x00, 0x00}, TrimRight)
//	// []byte{0x01, 0x02}
func TrimBytes(bytes []byte, dir TrimDirection) []byte {
	if len(bytes) == 0 {
		return bytes
	}

	sliceLength := 0
	for i := 0; i < len(bytes)-1; i++ {
		var idx int
		if dir == TrimLeft {
			idx = i
		} else {
			idx = len(bytes) - i - 1
		}

		if bytes[idx] == 0 {
			sliceLength++
		} else {
			break
		}
	}

	if dir == TrimLeft {
		return bytes[sliceLength:]
	}
	return bytes[:len(bytes)-sliceLength]
}

// TrimHex removes leading or trailing zeros from a hex string.
//
// Example:
//
//	result := TrimHex("0x00000102", TrimLeft)
//	// "0x0102"
//
//	result := TrimHex("0x01020000", TrimRight)
//	// "0x0102"
func TrimHex(hex string, dir TrimDirection) string {
	// Remove 0x prefix
	h := strings.TrimPrefix(hex, "0x")
	h = strings.TrimPrefix(h, "0X")

	if len(h) == 0 {
		return "0x"
	}

	sliceLength := 0
	for i := 0; i < len(h)-1; i++ {
		var idx int
		if dir == TrimLeft {
			idx = i
		} else {
			idx = len(h) - i - 1
		}

		if h[idx] == '0' {
			sliceLength++
		} else {
			break
		}
	}

	var result string
	if dir == TrimLeft {
		result = h[sliceLength:]
	} else {
		result = h[:len(h)-sliceLength]
	}

	// Ensure even length for proper hex encoding
	if len(result) == 1 && dir == TrimRight {
		result = result + "0"
	}
	if len(result)%2 == 1 {
		result = "0" + result
	}

	return "0x" + result
}

// TrimLeftBytes removes leading zero bytes.
func TrimLeftBytes(bytes []byte) []byte {
	return TrimBytes(bytes, TrimLeft)
}

// TrimRightBytes removes trailing zero bytes.
func TrimRightBytes(bytes []byte) []byte {
	return TrimBytes(bytes, TrimRight)
}

// TrimLeftHex removes leading zeros from a hex string.
func TrimLeftHex(hex string) string {
	return TrimHex(hex, TrimLeft)
}

// TrimRightHex removes trailing zeros from a hex string.
func TrimRightHex(hex string) string {
	return TrimHex(hex, TrimRight)
}
