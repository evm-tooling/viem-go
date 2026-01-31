package blob

import (
	"github.com/ChefBingbong/viem-go/utils/kzg"
)

// BlobSidecar represents a blob with its commitment and proof.
type BlobSidecar struct {
	Blob       []byte `json:"blob"`
	Commitment []byte `json:"commitment"`
	Proof      []byte `json:"proof"`
}

// BlobSidecarHex represents a blob sidecar with hex-encoded values.
type BlobSidecarHex struct {
	Blob       string `json:"blob"`
	Commitment string `json:"commitment"`
	Proof      string `json:"proof"`
}

// Re-export constants from kzg package for convenience
const (
	BytesPerFieldElement = kzg.BytesPerFieldElement
	FieldElementsPerBlob = kzg.FieldElementsPerBlob
	BytesPerBlob         = kzg.BytesPerBlob
	BlobsPerTransaction  = kzg.BlobsPerTransaction
	BytesPerCommitment   = kzg.BytesPerCommitment
	BytesPerProof        = kzg.BytesPerProof
)

// MaxBytesPerTransaction is re-exported from kzg package.
var MaxBytesPerTransaction = kzg.MaxBytesPerTransaction
