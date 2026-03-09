package errors

import "testing"

func TestInvalidDomainError_Error(t *testing.T) {
	err := &InvalidDomainError{Domain: "invalid"}
	testErrorContains(t, err, []string{"Invalid domain", "EIP-712"})
}

func TestInvalidPrimaryTypeError_Error(t *testing.T) {
	err := &InvalidPrimaryTypeError{PrimaryType: "Foo", Types: "[\"Mail\",\"Person\"]"}
	testErrorContains(t, err, []string{"Foo", "primary type"})
}

func TestInvalidStructTypeError_Error(t *testing.T) {
	err := &InvalidStructTypeError{Type: "uint256"}
	testErrorContains(t, err, []string{"invalid", "Solidity type"})
}
