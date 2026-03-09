package errors

import (
	"strings"
	"testing"
)

func TestSliceOffsetOutOfBoundsError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *SliceOffsetOutOfBoundsError
		contains []string
	}{
		{"start", &SliceOffsetOutOfBoundsError{Offset: 5, Position: SlicePositionStart, Size: 5}, []string{"out-of-bounds", "starting", "5"}},
		{"end", &SliceOffsetOutOfBoundsError{Offset: 10, Position: SlicePositionEnd, Size: 5}, []string{"out-of-bounds", "ending", "10"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func TestSizeExceedsPaddingSizeError_Error(t *testing.T) {
	err := &SizeExceedsPaddingSizeError{Size: 40, TargetSize: 32, Type: "hex"}
	got := err.Error()
	if !strings.Contains(got, "exceeds padding") || !strings.Contains(got, "40") || !strings.Contains(got, "32") {
		t.Errorf("Error() = %q, want to contain 'exceeds padding' and sizes", got)
	}
}

func TestInvalidBytesLengthError_Error(t *testing.T) {
	err := &InvalidBytesLengthError{Size: 20, TargetSize: 32, Type: "bytes"}
	testErrorContains(t, err, []string{"expected to be", "32", "20"})
}
