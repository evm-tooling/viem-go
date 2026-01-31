package kzg

import (
	"strings"
)

// SetupKzgParams contains the parameters for setting up KZG.
type SetupKzgParams struct {
	// BlobToKzgCommitment converts a blob to a KZG commitment.
	BlobToKzgCommitment func(blob []byte) ([]byte, error)
	// ComputeBlobKzgProof computes the KZG proof for a blob.
	ComputeBlobKzgProof func(blob []byte, commitment []byte) ([]byte, error)
	// LoadTrustedSetup loads the trusted setup from a file path.
	LoadTrustedSetup func(path string) error
}

// SetupKzg sets up and returns a KZG interface.
// It loads the trusted setup from the specified path before returning.
//
// If the trusted setup is already loaded, this function will still return
// a valid KZG interface (the error is ignored for this specific case).
//
// Example:
//
//	import "github.com/ethereum/c-kzg-4844/bindings/go"
//
//	kzg, err := SetupKzg(SetupKzgParams{
//		BlobToKzgCommitment: func(blob []byte) ([]byte, error) {
//			var b ckzg.Blob
//			copy(b[:], blob)
//			commitment, err := ckzg.BlobToKZGCommitment(&b)
//			if err != nil {
//				return nil, err
//			}
//			return commitment[:], nil
//		},
//		ComputeBlobKzgProof: func(blob, commitment []byte) ([]byte, error) {
//			var b ckzg.Blob
//			var c ckzg.Bytes48
//			copy(b[:], blob)
//			copy(c[:], commitment)
//			proof, err := ckzg.ComputeBlobKZGProof(&b, c)
//			if err != nil {
//				return nil, err
//			}
//			return proof[:], nil
//		},
//		LoadTrustedSetup: func(path string) error {
//			return ckzg.LoadTrustedSetupFile(path)
//		},
//	}, "/path/to/trusted_setup.txt")
func SetupKzg(params SetupKzgParams, trustedSetupPath string) (Kzg, error) {
	// Load the trusted setup
	if params.LoadTrustedSetup != nil {
		err := params.LoadTrustedSetup(trustedSetupPath)
		if err != nil {
			// Ignore "already loaded" errors
			if !strings.Contains(err.Error(), "trusted setup is already loaded") &&
				!strings.Contains(err.Error(), "already loaded") {
				return nil, err
			}
		}
	}

	// Return the KZG interface
	return DefineKzg(DefineKzgParams{
		BlobToKzgCommitment: params.BlobToKzgCommitment,
		ComputeBlobKzgProof: params.ComputeBlobKzgProof,
	}), nil
}

// SetupKzgWithLoader sets up KZG using a KzgWithSetup implementation.
// This is useful when you have a struct that implements both Kzg and SetupLoader.
func SetupKzgWithLoader(kzg KzgWithSetup, trustedSetupPath string) (Kzg, error) {
	err := kzg.LoadTrustedSetup(trustedSetupPath)
	if err != nil {
		// Ignore "already loaded" errors
		if !strings.Contains(err.Error(), "trusted setup is already loaded") &&
			!strings.Contains(err.Error(), "already loaded") {
			return nil, err
		}
	}
	return kzg, nil
}
