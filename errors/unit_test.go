package errors

import (
	"strings"
	"testing"
)

func TestInvalidDecimalNumberError_Error(t *testing.T) {
	err := &InvalidDecimalNumberError{Value: "1.2.3"}
	got := err.Error()
	if !strings.Contains(got, "decimal number") || !strings.Contains(got, "1.2.3") {
		t.Errorf("Error() = %q, want to contain 'decimal number' and value", got)
	}
}
