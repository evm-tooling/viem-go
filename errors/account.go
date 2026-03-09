package errors

import "fmt"

// AccountNotFoundError is returned when no account is provided to an action that requires one.
type AccountNotFoundError struct {
	DocsPath string
}

func (e *AccountNotFoundError) Error() string {
	msg := "Could not find an Account to execute with this Action."
	msg += "\nPlease provide an Account with the `account` argument on the Action, or by supplying an `account` to the Client."
	if e.DocsPath != "" {
		msg += fmt.Sprintf("\nDocs: https://viem.sh%s", e.DocsPath)
	}
	return msg
}

// AccountTypeNotSupportedError is returned when account type is not supported.
type AccountTypeNotSupportedError struct {
	Type         string
	DocsPath     string
	MetaMessages []string
}

func (e *AccountTypeNotSupportedError) Error() string {
	msg := fmt.Sprintf("Account type \"%s\" is not supported.", e.Type)
	for _, m := range e.MetaMessages {
		msg += "\n" + m
	}
	if e.DocsPath != "" {
		msg += fmt.Sprintf("\nDocs: https://viem.sh%s", e.DocsPath)
	}
	return msg
}
