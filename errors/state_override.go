package errors

import "fmt"

// AccountStateConflictError is returned when account state is set multiple times.
type AccountStateConflictError struct {
	Address string
}

func (e *AccountStateConflictError) Error() string {
	return fmt.Sprintf("State for account \"%s\" is set multiple times.", e.Address)
}

// StateAssignmentConflictError is returned when state and stateDiff are both set.
type StateAssignmentConflictError struct{}

func (e *StateAssignmentConflictError) Error() string {
	return "state and stateDiff are set on the same account."
}
