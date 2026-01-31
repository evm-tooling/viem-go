package kzg

import "errors"

var (
	// ErrKzgNotInitialized is returned when KZG functions are called before setup.
	ErrKzgNotInitialized = errors.New("kzg not initialized")

	// ErrTrustedSetupAlreadyLoaded is returned when trying to load setup twice.
	ErrTrustedSetupAlreadyLoaded = errors.New("trusted setup is already loaded")

	// ErrInvalidBlobSize is returned when a blob has an invalid size.
	ErrInvalidBlobSize = errors.New("invalid blob size")

	// ErrInvalidCommitmentSize is returned when a commitment has an invalid size.
	ErrInvalidCommitmentSize = errors.New("invalid commitment size")

	// ErrEmptyBlob is returned when trying to process empty data.
	ErrEmptyBlob = errors.New("blob data is empty")

	// ErrBlobSizeTooLarge is returned when data exceeds the maximum blob size.
	ErrBlobSizeTooLarge = errors.New("blob size too large")
)
