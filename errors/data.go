package errors

import "fmt"

// SlicePosition is the position for slice offset (start or end).
type SlicePosition string

const (
	SlicePositionStart SlicePosition = "start"
	SlicePositionEnd   SlicePosition = "end"
)

// SliceOffsetOutOfBoundsError is returned when slice offset is out of bounds.
type SliceOffsetOutOfBoundsError struct {
	Offset   int
	Position SlicePosition
	Size     int
}

func (e *SliceOffsetOutOfBoundsError) Error() string {
	posStr := "ending"
	if e.Position == SlicePositionStart {
		posStr = "starting"
	}
	return fmt.Sprintf("Slice %s at offset \"%d\" is out-of-bounds (size: %d).", posStr, e.Offset, e.Size)
}

// SizeExceedsPaddingSizeError is returned when size exceeds padding size.
type SizeExceedsPaddingSizeError struct {
	Size       int
	TargetSize int
	Type       string
}

func (e *SizeExceedsPaddingSizeError) Error() string {
	typeStr := "Hex"
	if e.Type == "bytes" {
		typeStr = "Bytes"
	}
	return fmt.Sprintf("%s size (%d) exceeds padding size (%d).", typeStr, e.Size, e.TargetSize)
}

// InvalidBytesLengthError is returned when bytes/hex length is invalid.
type InvalidBytesLengthError struct {
	Size       int
	TargetSize int
	Type       string
}

func (e *InvalidBytesLengthError) Error() string {
	typeStr := "Hex"
	if e.Type == "bytes" {
		typeStr = "Bytes"
	}
	return fmt.Sprintf("%s is expected to be %d %s long, but is %d %s long.", typeStr, e.TargetSize, e.Type, e.Size, e.Type)
}
