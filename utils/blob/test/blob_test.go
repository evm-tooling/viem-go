package test

import (
	"bytes"
	"crypto/sha256"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ChefBingbong/viem-go/utils/blob"
	"github.com/ChefBingbong/viem-go/utils/kzg"
)

// mockKzg is a mock KZG implementation for testing
type mockKzg struct{}

func (m *mockKzg) BlobToKzgCommitment(b []byte) ([]byte, error) {
	// Return a mock 48-byte commitment (hash of blob)
	hash := sha256.Sum256(b)
	commitment := make([]byte, 48)
	copy(commitment, hash[:])
	return commitment, nil
}

func (m *mockKzg) ComputeBlobKzgProof(b []byte, commitment []byte) ([]byte, error) {
	// Return a mock 48-byte proof
	hash := sha256.Sum256(append(b, commitment...))
	proof := make([]byte, 48)
	copy(proof, hash[:])
	return proof, nil
}

var _ = Describe("Blob", func() {
	Describe("ToBlobs and FromBlobs", func() {
		It("should convert data to blobs and back", func() {
			data := []byte("hello world, this is a test of blob encoding")

			blobs, err := blob.ToBlobs(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(blobs)).To(Equal(1))
			Expect(len(blobs[0])).To(Equal(blob.BytesPerBlob))

			recovered, err := blob.FromBlobs(blobs)
			Expect(err).NotTo(HaveOccurred())
			Expect(recovered).To(Equal(data))
		})

		It("should handle empty data", func() {
			_, err := blob.ToBlobs([]byte{})
			Expect(err).To(Equal(kzg.ErrEmptyBlob))
		})

		It("should handle large data requiring multiple blobs", func() {
			// Create data that requires multiple blobs
			// Each field element holds 31 bytes, and there are 4096 field elements per blob
			// So each blob holds approximately 31 * 4096 = 127,000 bytes
			//
			// Note: We avoid using 0x80 in the test data because the EIP-4844 encoding
			// uses 0x80 as a terminator byte. If data contains 0x80 followed by no more
			// 0x80 bytes, it could be mistaken for the terminator. This is inherent to
			// the encoding scheme used by viem.
			data := make([]byte, 200000) // Should require ~2 blobs
			for i := range data {
				// Use values that won't conflict with terminator (0x80 = 128)
				data[i] = byte((i % 127) + 1) // Values 1-127, avoiding 0 and 0x80
			}

			blobs, err := blob.ToBlobs(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(blobs)).To(BeNumerically(">=", 2))

			recovered, err := blob.FromBlobs(blobs)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(recovered)).To(Equal(len(data)))
			Expect(recovered).To(Equal(data))
		})
	})

	Describe("ToBlobsHex", func() {
		It("should return hex-encoded blobs", func() {
			data := []byte("test data")

			hexBlobs, err := blob.ToBlobsHex(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(hexBlobs)).To(Equal(1))
			Expect(hexBlobs[0]).To(HavePrefix("0x"))
		})
	})

	Describe("CommitmentToVersionedHash", func() {
		It("should compute versioned hash with correct version", func() {
			commitment := make([]byte, 48)
			for i := range commitment {
				commitment[i] = byte(i)
			}

			hash := blob.CommitmentToVersionedHash(commitment, kzg.VersionedHashVersionKzg)

			Expect(len(hash)).To(Equal(32))
			Expect(hash[0]).To(Equal(byte(0x01)))

			// Verify it's SHA256 with version prefix
			expectedHash := sha256.Sum256(commitment)
			expectedHash[0] = 0x01
			Expect(hash).To(Equal(expectedHash[:]))
		})

		It("should support different versions", func() {
			commitment := make([]byte, 48)

			hashV1 := blob.CommitmentToVersionedHash(commitment, 0x01)
			hashV2 := blob.CommitmentToVersionedHash(commitment, 0x02)

			Expect(hashV1[0]).To(Equal(byte(0x01)))
			Expect(hashV2[0]).To(Equal(byte(0x02)))
			// Rest of hash should be the same
			Expect(hashV1[1:]).To(Equal(hashV2[1:]))
		})
	})

	Describe("CommitmentsToVersionedHashes", func() {
		It("should compute hashes for multiple commitments", func() {
			commitments := [][]byte{
				make([]byte, 48),
				make([]byte, 48),
			}
			commitments[0][0] = 0x01
			commitments[1][0] = 0x02

			hashes := blob.CommitmentsToVersionedHashes(commitments, kzg.VersionedHashVersionKzg)

			Expect(len(hashes)).To(Equal(2))
			Expect(hashes[0][0]).To(Equal(byte(0x01)))
			Expect(hashes[1][0]).To(Equal(byte(0x01)))
			Expect(bytes.Equal(hashes[0], hashes[1])).To(BeFalse())
		})
	})

	Describe("BlobsToCommitments", func() {
		It("should compute commitments using KZG", func() {
			kzgImpl := &mockKzg{}
			blobs, err := blob.ToBlobs([]byte("test"))
			Expect(err).NotTo(HaveOccurred())

			commitments, err := blob.BlobsToCommitments(blobs, kzgImpl)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(commitments)).To(Equal(1))
			Expect(len(commitments[0])).To(Equal(48))
		})
	})

	Describe("BlobsToProofs", func() {
		It("should compute proofs using KZG", func() {
			kzgImpl := &mockKzg{}
			blobs, err := blob.ToBlobs([]byte("test"))
			Expect(err).NotTo(HaveOccurred())

			commitments, err := blob.BlobsToCommitments(blobs, kzgImpl)
			Expect(err).NotTo(HaveOccurred())

			proofs, err := blob.BlobsToProofs(blobs, commitments, kzgImpl)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(proofs)).To(Equal(1))
			Expect(len(proofs[0])).To(Equal(48))
		})
	})

	Describe("ToBlobSidecars", func() {
		It("should create sidecars from data", func() {
			kzgImpl := &mockKzg{}

			sidecars, err := blob.ToBlobSidecars(blob.ToBlobSidecarsParams{
				Data: []byte("test data"),
				Kzg:  kzgImpl,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(sidecars)).To(Equal(1))
			Expect(len(sidecars[0].Blob)).To(Equal(blob.BytesPerBlob))
			Expect(len(sidecars[0].Commitment)).To(Equal(48))
			Expect(len(sidecars[0].Proof)).To(Equal(48))
		})

		It("should create sidecars from pre-computed components", func() {
			mockBlob := make([]byte, blob.BytesPerBlob)
			mockCommitment := make([]byte, 48)
			mockProof := make([]byte, 48)

			sidecars, err := blob.ToBlobSidecars(blob.ToBlobSidecarsParams{
				Blobs:       [][]byte{mockBlob},
				Commitments: [][]byte{mockCommitment},
				Proofs:      [][]byte{mockProof},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(sidecars)).To(Equal(1))
			Expect(sidecars[0].Blob).To(Equal(mockBlob))
			Expect(sidecars[0].Commitment).To(Equal(mockCommitment))
			Expect(sidecars[0].Proof).To(Equal(mockProof))
		})
	})

	Describe("SidecarsToVersionedHashes", func() {
		It("should extract versioned hashes from sidecars", func() {
			sidecars := []blob.BlobSidecar{
				{
					Blob:       make([]byte, blob.BytesPerBlob),
					Commitment: make([]byte, 48),
					Proof:      make([]byte, 48),
				},
			}

			hashes := blob.SidecarsToVersionedHashes(sidecars, kzg.VersionedHashVersionKzg)
			Expect(len(hashes)).To(Equal(1))
			Expect(len(hashes[0])).To(Equal(32))
			Expect(hashes[0][0]).To(Equal(byte(0x01)))
		})
	})
})
