package errors

import "fmt"

// FilterTypeNotSupportedError is returned when filter type is not supported.
type FilterTypeNotSupportedError struct {
	Type string
}

func (e *FilterTypeNotSupportedError) Error() string {
	return fmt.Sprintf("Filter type \"%s\" is not supported.", e.Type)
}
