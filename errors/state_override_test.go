package errors

import "testing"

func TestAccountStateConflictError_Error(t *testing.T) {
	err := &AccountStateConflictError{Address: "0x1234"}
	testErrorContains(t, err, []string{"set multiple times", "0x1234"})
}

func TestStateAssignmentConflictError_Error(t *testing.T) {
	err := &StateAssignmentConflictError{}
	testErrorContains(t, err, []string{"state and stateDiff"})
}
