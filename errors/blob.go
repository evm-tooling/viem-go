package errors

import "fmt"

// VersionedHashVersionKzg is the expected version for KZG blob hashes.
const VersionedHashVersionKzg = 1

// BlobSizeTooLargeError is returned when blob size exceeds maximum.
type BlobSizeTooLargeError struct {
	MaxSize int
	Size    int
}

func (e *BlobSizeTooLargeError) Error() string {
	return fmt.Sprintf("Blob size is too large.\nMax: %d bytes\nGiven: %d bytes", e.MaxSize, e.Size)
}

// EmptyBlobError is returned when blob data is empty.
type EmptyBlobError struct{}

func (e *EmptyBlobError) Error() string {
	return "Blob data must not be empty."
}

// InvalidVersionedHashSizeError is returned when versioned hash size is invalid.
type InvalidVersionedHashSizeError struct {
	Hash   string
	Size   int
}

func (e *InvalidVersionedHashSizeError) Error() string {
	return fmt.Sprintf("Versioned hash \"%s\" size is invalid.\nExpected: 32\nReceived: %d", e.Hash, e.Size)
}

// InvalidVersionedHashVersionError is returned when versioned hash version is invalid.
type InvalidVersionedHashVersionError struct {
	Hash    string
	Version int
}

func (e *InvalidVersionedHashVersionError) Error() string {
	return fmt.Sprintf("Versioned hash \"%s\" version is invalid.\nExpected: %d\nReceived: %d", e.Hash, VersionedHashVersionKzg, e.Version)
}
