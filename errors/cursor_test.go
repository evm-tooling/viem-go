package errors

import (
	"strings"
	"testing"
)

func TestNegativeOffsetError_Error(t *testing.T) {
	err := &NegativeOffsetError{Offset: -5}
	testErrorContains(t, err, []string{"-5", "negative"})
}

func TestPositionOutOfBoundsError_Error(t *testing.T) {
	err := &PositionOutOfBoundsError{Length: 10, Position: 15}
	testErrorContains(t, err, []string{"out of bounds", "15", "10"})
}

func TestRecursiveReadLimitExceededError_Error(t *testing.T) {
	err := &RecursiveReadLimitExceededError{Count: 10, Limit: 5}
	got := err.Error()
	if !strings.Contains(got, "exceeded") || !strings.Contains(got, "10") || !strings.Contains(got, "5") {
		t.Errorf("Error() = %q, want to contain 'exceeded' and counts", got)
	}
}
