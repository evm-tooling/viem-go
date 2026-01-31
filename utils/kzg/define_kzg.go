package kzg

// DefineKzgParams contains the parameters for defining a KZG interface.
type DefineKzgParams struct {
	// BlobToKzgCommitment converts a blob to a KZG commitment.
	BlobToKzgCommitment func(blob []byte) ([]byte, error)
	// ComputeBlobKzgProof computes the KZG proof for a blob.
	ComputeBlobKzgProof func(blob []byte, commitment []byte) ([]byte, error)
}

// DefineKzg creates a Kzg implementation from individual functions.
// This is useful when working with external KZG libraries that expose
// separate functions rather than a unified interface.
//
// Example:
//
//	import "github.com/ethereum/c-kzg-4844/bindings/go"
//
//	kzg := DefineKzg(DefineKzgParams{
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
//	})
func DefineKzg(params DefineKzgParams) Kzg {
	return &KzgFunctions{
		BlobToKzgCommitmentFn: params.BlobToKzgCommitment,
		ComputeBlobKzgProofFn: params.ComputeBlobKzgProof,
	}
}

// DefineKzgFromInterface wraps an existing Kzg implementation.
// This is useful when you already have a struct implementing the Kzg interface.
func DefineKzgFromInterface(kzg Kzg) Kzg {
	return kzg
}
