package errors

import "fmt"

// SiweInvalidMessageFieldError is returned when SIWE message field is invalid.
type SiweInvalidMessageFieldError struct {
	Field        string
	DocsPath     string
	MetaMessages []string
}

func (e *SiweInvalidMessageFieldError) Error() string {
	msg := fmt.Sprintf("Invalid Sign-In with Ethereum message field \"%s\".", e.Field)
	for _, m := range e.MetaMessages {
		msg += "\n" + m
	}
	if e.DocsPath != "" {
		msg += fmt.Sprintf("\nDocs: https://viem.sh%s", e.DocsPath)
	}
	return msg
}
