package kzg

// Kzg defines the interface for KZG cryptographic operations.
// This is used for EIP-4844 blob transactions.
//
// Implementations of this interface can be provided by external KZG libraries
// such as c-kzg-4844 or go-kzg-4844.
type Kzg interface {
	// BlobToKzgCommitment converts a blob to a KZG commitment.
	// The blob should be exactly 131072 bytes (4096 field elements * 32 bytes).
	// Returns a 48-byte commitment.
	BlobToKzgCommitment(blob []byte) ([]byte, error)

	// ComputeBlobKzgProof computes the KZG proof for a blob given its commitment.
	// The blob should be exactly 131072 bytes.
	// The commitment should be exactly 48 bytes.
	// Returns a 48-byte proof.
	ComputeBlobKzgProof(blob []byte, commitment []byte) ([]byte, error)
}

// KzgFunctions contains the KZG function implementations.
// This is used when you have separate functions instead of a struct implementing Kzg.
type KzgFunctions struct {
	BlobToKzgCommitmentFn   func(blob []byte) ([]byte, error)
	ComputeBlobKzgProofFn   func(blob []byte, commitment []byte) ([]byte, error)
}

// BlobToKzgCommitment implements Kzg interface.
func (k *KzgFunctions) BlobToKzgCommitment(blob []byte) ([]byte, error) {
	if k.BlobToKzgCommitmentFn == nil {
		return nil, ErrKzgNotInitialized
	}
	return k.BlobToKzgCommitmentFn(blob)
}

// ComputeBlobKzgProof implements Kzg interface.
func (k *KzgFunctions) ComputeBlobKzgProof(blob []byte, commitment []byte) ([]byte, error) {
	if k.ComputeBlobKzgProofFn == nil {
		return nil, ErrKzgNotInitialized
	}
	return k.ComputeBlobKzgProofFn(blob, commitment)
}

// SetupLoader defines an interface for loading trusted setup.
type SetupLoader interface {
	// LoadTrustedSetup loads the KZG trusted setup from a file path.
	LoadTrustedSetup(path string) error
}

// KzgWithSetup combines the Kzg interface with setup loading capability.
type KzgWithSetup interface {
	Kzg
	SetupLoader
}

// Constants for EIP-4844

const (
	// BytesPerFieldElement is the number of bytes in a BLS scalar field element.
	BytesPerFieldElement = 32

	// FieldElementsPerBlob is the number of field elements in a blob.
	FieldElementsPerBlob = 4096

	// BytesPerBlob is the total number of bytes in a blob.
	BytesPerBlob = BytesPerFieldElement * FieldElementsPerBlob // 131072 bytes

	// BlobsPerTransaction is the maximum number of blobs per transaction.
	BlobsPerTransaction = 6

	// BytesPerCommitment is the size of a KZG commitment.
	BytesPerCommitment = 48

	// BytesPerProof is the size of a KZG proof.
	BytesPerProof = 48

	// VersionedHashVersionKzg is the version byte for KZG versioned hashes.
	VersionedHashVersionKzg = 0x01
)

// MaxBytesPerTransaction is the maximum data bytes per transaction.
// This accounts for the zero byte prefix and terminator byte.
var MaxBytesPerTransaction = BytesPerBlob*BlobsPerTransaction -
	1 - // terminator byte (0x80)
	1*FieldElementsPerBlob*BlobsPerTransaction // zero byte (0x00) per field element
