package errors

import "fmt"

// InvalidDecimalNumberError is returned when value is not a valid decimal number.
type InvalidDecimalNumberError struct {
	Value string
}

func (e *InvalidDecimalNumberError) Error() string {
	return fmt.Sprintf("Number `%s` is not a valid decimal number.", e.Value)
}
