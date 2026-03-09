package errors

import (
	"strings"
	"testing"
)

func TestEstimateGasExecutionError_Error(t *testing.T) {
	err := &EstimateGasExecutionError{Cause: &BaseError{ShortMessage: "revert"}}
	got := err.Error()
	if !strings.Contains(got, "revert") {
		t.Errorf("Error() = %q, want to contain cause message", got)
	}
}
