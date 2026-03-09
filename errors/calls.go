package errors

import "fmt"

// BundleFailedError is returned when call bundle fails.
type BundleFailedError struct {
	StatusCode int
}

func (e *BundleFailedError) Error() string {
	return fmt.Sprintf("Call bundle failed with status: %d", e.StatusCode)
}
