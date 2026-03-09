package errors

import "fmt"

// InvalidDomainError is returned when EIP-712 domain is invalid.
type InvalidDomainError struct {
	Domain string
}

func (e *InvalidDomainError) Error() string {
	msg := fmt.Sprintf("Invalid domain \"%s\".", e.Domain)
	msg += "\nMust be a valid EIP-712 domain."
	return msg
}

// InvalidPrimaryTypeError is returned when primary type is invalid.
type InvalidPrimaryTypeError struct {
	PrimaryType string
	Types       string
}

func (e *InvalidPrimaryTypeError) Error() string {
	msg := fmt.Sprintf("Invalid primary type `%s` must be one of `%s`.", e.PrimaryType, e.Types)
	msg += "\nCheck that the primary type is a key in `types`."
	return msg
}

// InvalidStructTypeError is returned when struct type is invalid.
type InvalidStructTypeError struct {
	Type string
}

func (e *InvalidStructTypeError) Error() string {
	msg := fmt.Sprintf("Struct type \"%s\" is invalid.", e.Type)
	msg += "\nStruct type must not be a Solidity type."
	return msg
}
