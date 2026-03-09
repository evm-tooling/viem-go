package errors

// EstimateGasExecutionError is returned when estimate gas execution fails.
type EstimateGasExecutionError struct {
	Cause        error
	MetaMessages []string
}

func (e *EstimateGasExecutionError) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return "Estimate gas execution failed."
}

func (e *EstimateGasExecutionError) Unwrap() error {
	return e.Cause
}
