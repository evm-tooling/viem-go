package errors

import "testing"

func TestBlobSizeTooLargeError_Error(t *testing.T) {
	err := &BlobSizeTooLargeError{MaxSize: 128, Size: 256}
	testErrorContains(t, err, []string{"too large", "128", "256"})
}

func TestEmptyBlobError_Error(t *testing.T) {
	err := &EmptyBlobError{}
	testErrorContains(t, err, []string{"must not be empty"})
}

func TestInvalidVersionedHashSizeError_Error(t *testing.T) {
	err := &InvalidVersionedHashSizeError{Hash: "0x1234", Size: 16}
	testErrorContains(t, err, []string{"size is invalid", "32", "16"})
}

func TestInvalidVersionedHashVersionError_Error(t *testing.T) {
	err := &InvalidVersionedHashVersionError{Hash: "0xabcd", Version: 2}
	testErrorContains(t, err, []string{"version is invalid", "Expected"})
}
