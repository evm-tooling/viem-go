package errors

import (
	"strings"
	"testing"
)

func TestOffchainLookupError_Error(t *testing.T) {
	err := &OffchainLookupError{Cause: &BaseError{ShortMessage: "fetch failed"}}
	testErrorContains(t, err, []string{"fetch failed"})
}

func TestOffchainLookupResponseMalformedError_Error(t *testing.T) {
	err := &OffchainLookupResponseMalformedError{URL: "https://gateway.io", Result: "invalid"}
	testErrorContains(t, err, []string{"malformed", "hex value"})
}

func TestOffchainLookupSenderMismatchError_Error(t *testing.T) {
	err := &OffchainLookupSenderMismatchError{Sender: "0xaaa", To: "0xbbb"}
	got := err.Error()
	if !strings.Contains(got, "does not match") || !strings.Contains(got, "0xaaa") {
		t.Errorf("Error() = %q, want to contain mismatch and addresses", got)
	}
}
