package errors

import "strings"

// EstimateGasExecutionError is returned when estimate gas execution fails.
type EstimateGasExecutionError struct {
	Cause        error
	MetaMessages []string
}

func (e *EstimateGasExecutionError) Error() string {
	base := "Estimate gas execution failed."
	if e.Cause != nil {
		base = e.Cause.Error()
	}
	if len(e.MetaMessages) > 0 {
		base += "\n\n" + strings.Join(e.MetaMessages, "\n")
	}
	return "EstimateGasExecutionError: " + base
}

func (e *EstimateGasExecutionError) Unwrap() error {
	return e.Cause
}
