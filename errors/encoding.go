package errors

import "fmt"

// IntegerOutOfRangeError is returned when a number is not in safe integer range.
type IntegerOutOfRangeError struct {
	Value  string
	Min    string
	Max    string
	Signed bool
	Size   int
}

func (e *IntegerOutOfRangeError) Error() string {
	rangeStr := ""
	if e.Max != "" {
		rangeStr = fmt.Sprintf("(%s to %s)", e.Min, e.Max)
	} else {
		rangeStr = fmt.Sprintf("(above %s)", e.Min)
	}
	sizeStr := ""
	if e.Size > 0 {
		signedStr := "unsigned"
		if e.Signed {
			signedStr = "signed"
		}
		sizeStr = fmt.Sprintf("%d-bit %s ", e.Size*8, signedStr)
	}
	return fmt.Sprintf("Number \"%s\" is not in safe %sinteger range %s", e.Value, sizeStr, rangeStr)
}

// InvalidBytesBooleanError is returned when bytes value is not a valid boolean.
type InvalidBytesBooleanError struct {
	Bytes string
}

func (e *InvalidBytesBooleanError) Error() string {
	return fmt.Sprintf("Bytes value \"%s\" is not a valid boolean. The bytes array must contain a single byte of either a 0 or 1 value.", e.Bytes)
}

// InvalidHexBooleanError is returned when hex value is not a valid boolean.
type InvalidHexBooleanError struct {
	Hex string
}

func (e *InvalidHexBooleanError) Error() string {
	return fmt.Sprintf("Hex value \"%s\" is not a valid boolean. The hex value must be \"0x0\" (false) or \"0x1\" (true).", e.Hex)
}

// InvalidHexValueError is returned when hex value has odd length.
type InvalidHexValueError struct {
	Value  string
	Length int
}

func (e *InvalidHexValueError) Error() string {
	length := e.Length
	if length == 0 && e.Value != "" {
		length = len(e.Value)
	}
	return fmt.Sprintf("Hex value \"%s\" is an odd length (%d). It must be an even length.", e.Value, length)
}

// SizeOverflowError is returned when size exceeds maximum.
type SizeOverflowError struct {
	GivenSize int
	MaxSize   int
}

func (e *SizeOverflowError) Error() string {
	return fmt.Sprintf("Size cannot exceed %d bytes. Given size: %d bytes.", e.MaxSize, e.GivenSize)
}
