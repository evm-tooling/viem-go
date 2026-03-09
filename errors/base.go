package errors

import (
	"errors"
	"fmt"
	"strings"
)

// BaseError provides common fields for viem-style errors.
// Errors may embed BaseError or use it as a reference for formatting.
type BaseError struct {
	Name         string
	ShortMessage string
	Details      string
	DocsPath     string
	MetaMessages []string
	Cause        error
}

// Error implements the error interface.
func (e *BaseError) Error() string {
	var parts []string
	if e.ShortMessage != "" {
		parts = append(parts, e.ShortMessage)
	} else {
		parts = append(parts, "An error occurred.")
	}
	if len(e.MetaMessages) > 0 {
		parts = append(parts, "")
		parts = append(parts, e.MetaMessages...)
	}
	if e.DocsPath != "" {
		parts = append(parts, "", fmt.Sprintf("Docs: https://viem.sh%s", e.DocsPath))
	}
	if e.Details != "" {
		parts = append(parts, "", fmt.Sprintf("Details: %s", e.Details))
	}
	return strings.TrimSpace(strings.Join(parts, "\n"))
}

// Unwrap returns the underlying cause for error chaining.
func (e *BaseError) Unwrap() error {
	return e.Cause
}

// unwrapper is the interface for errors that wrap another error.
type unwrapper interface {
	Unwrap() error
}

// Walk walks the cause chain and returns the first error for which fn returns true.
// If fn is nil, returns the leaf error (end of chain).
// If no error matches fn, returns nil.
func Walk(err error, fn func(err error) bool) error {
	if err == nil {
		return nil
	}
	if fn != nil && fn(err) {
		return err
	}
	var u unwrapper
	if errors.As(err, &u) {
		return Walk(u.Unwrap(), fn)
	}
	if fn != nil {
		return nil
	}
	return err
}
