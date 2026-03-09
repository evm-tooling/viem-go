package errors

import (
	"strings"
	"testing"
)

func TestIntegerOutOfRangeError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *IntegerOutOfRangeError
		contains []string
	}{
		{"basic", &IntegerOutOfRangeError{Value: "99999999999999999999", Min: "0"}, []string{"not in safe", "99999999999999999999"}},
		{"with range", &IntegerOutOfRangeError{Value: "300", Min: "0", Max: "255"}, []string{"not in safe", "0 to 255"}},
		{"with size", &IntegerOutOfRangeError{Value: "x", Min: "0", Size: 4, Signed: false}, []string{"32-bit", "unsigned"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { testErrorContains(t, tt.err, tt.contains) })
	}
}

func TestInvalidBytesBooleanError_Error(t *testing.T) {
	err := &InvalidBytesBooleanError{Bytes: "0x02"}
	testErrorContains(t, err, []string{"not a valid boolean", "0 or 1"})
}

func TestInvalidHexBooleanError_Error(t *testing.T) {
	err := &InvalidHexBooleanError{Hex: "0x2"}
	testErrorContains(t, err, []string{"0x0", "0x1"})
}

func TestInvalidHexValueError_Error(t *testing.T) {
	err := &InvalidHexValueError{Value: "0x123", Length: 5}
	got := err.Error()
	if !strings.Contains(got, "odd length") || !strings.Contains(got, "5") {
		t.Errorf("Error() = %q, want to contain 'odd length' and length", got)
	}
}

func TestSizeOverflowError_Error(t *testing.T) {
	err := &SizeOverflowError{GivenSize: 100, MaxSize: 32}
	testErrorContains(t, err, []string{"cannot exceed", "32", "100"})
}
