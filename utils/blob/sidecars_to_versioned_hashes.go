package blob

import (
	"github.com/ChefBingbong/viem-go/utils/kzg"
)

// SidecarsToVersionedHashes extracts versioned hashes from blob sidecars.
//
// Example:
//
//	sidecars, _ := ToBlobSidecars(params)
//	versionedHashes := SidecarsToVersionedHashes(sidecars, kzg.VersionedHashVersionKzg)
func SidecarsToVersionedHashes(sidecars []BlobSidecar, version byte) [][]byte {
	hashes := make([][]byte, len(sidecars))
	for i, sidecar := range sidecars {
		hashes[i] = CommitmentToVersionedHash(sidecar.Commitment, version)
	}
	return hashes
}

// SidecarsToVersionedHashesDefault uses the default KZG version (0x01).
func SidecarsToVersionedHashesDefault(sidecars []BlobSidecar) [][]byte {
	return SidecarsToVersionedHashes(sidecars, kzg.VersionedHashVersionKzg)
}

// SidecarsToVersionedHashesHex returns the versioned hashes as hex strings.
func SidecarsToVersionedHashesHex(sidecars []BlobSidecar, version byte) []string {
	hashes := make([]string, len(sidecars))
	for i, sidecar := range sidecars {
		hashes[i] = CommitmentToVersionedHashHex(sidecar.Commitment, version)
	}
	return hashes
}

// SidecarsHexToVersionedHashes computes versioned hashes from hex sidecars.
func SidecarsHexToVersionedHashes(sidecars []BlobSidecarHex, version byte) ([][]byte, error) {
	hashes := make([][]byte, len(sidecars))
	for i, sidecar := range sidecars {
		commitment, err := hexToBytes(sidecar.Commitment)
		if err != nil {
			return nil, err
		}
		hashes[i] = CommitmentToVersionedHash(commitment, version)
	}
	return hashes, nil
}

// SidecarsHexToVersionedHashesHex returns hex hashes from hex sidecars.
func SidecarsHexToVersionedHashesHex(sidecars []BlobSidecarHex, version byte) ([]string, error) {
	hashes := make([]string, len(sidecars))
	for i, sidecar := range sidecars {
		commitment, err := hexToBytes(sidecar.Commitment)
		if err != nil {
			return nil, err
		}
		hashes[i] = bytesToHex(CommitmentToVersionedHash(commitment, version))
	}
	return hashes, nil
}
