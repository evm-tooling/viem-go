package test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ChefBingbong/viem-go/utils/kzg"
)

var _ = Describe("Kzg", func() {
	Describe("Constants", func() {
		It("should have correct blob constants", func() {
			Expect(kzg.BytesPerFieldElement).To(Equal(32))
			Expect(kzg.FieldElementsPerBlob).To(Equal(4096))
			Expect(kzg.BytesPerBlob).To(Equal(131072))
			Expect(kzg.BlobsPerTransaction).To(Equal(6))
			Expect(kzg.BytesPerCommitment).To(Equal(48))
			Expect(kzg.BytesPerProof).To(Equal(48))
			Expect(kzg.VersionedHashVersionKzg).To(Equal(0x01))
		})

		It("should calculate max bytes per transaction", func() {
			// MaxBytesPerTransaction = BytesPerBlob * BlobsPerTransaction - 1 - FieldElementsPerBlob * BlobsPerTransaction
			expected := 131072*6 - 1 - 4096*6
			Expect(kzg.MaxBytesPerTransaction).To(Equal(expected))
		})
	})

	Describe("DefineKzg", func() {
		It("should create KZG interface from functions", func() {
			callCount := 0
			kzgImpl := kzg.DefineKzg(kzg.DefineKzgParams{
				BlobToKzgCommitment: func(blob []byte) ([]byte, error) {
					callCount++
					return make([]byte, 48), nil
				},
				ComputeBlobKzgProof: func(blob, commitment []byte) ([]byte, error) {
					callCount++
					return make([]byte, 48), nil
				},
			})

			_, err := kzgImpl.BlobToKzgCommitment(make([]byte, 131072))
			Expect(err).NotTo(HaveOccurred())
			Expect(callCount).To(Equal(1))

			_, err = kzgImpl.ComputeBlobKzgProof(make([]byte, 131072), make([]byte, 48))
			Expect(err).NotTo(HaveOccurred())
			Expect(callCount).To(Equal(2))
		})

		It("should return error when functions not set", func() {
			kzgImpl := kzg.DefineKzg(kzg.DefineKzgParams{})

			_, err := kzgImpl.BlobToKzgCommitment(nil)
			Expect(err).To(Equal(kzg.ErrKzgNotInitialized))

			_, err = kzgImpl.ComputeBlobKzgProof(nil, nil)
			Expect(err).To(Equal(kzg.ErrKzgNotInitialized))
		})
	})

	Describe("SetupKzg", func() {
		It("should setup KZG with trusted setup", func() {
			setupLoaded := false
			kzgImpl, err := kzg.SetupKzg(kzg.SetupKzgParams{
				BlobToKzgCommitment: func(blob []byte) ([]byte, error) {
					return make([]byte, 48), nil
				},
				ComputeBlobKzgProof: func(blob, commitment []byte) ([]byte, error) {
					return make([]byte, 48), nil
				},
				LoadTrustedSetup: func(path string) error {
					setupLoaded = true
					return nil
				},
			}, "/path/to/setup.txt")

			Expect(err).NotTo(HaveOccurred())
			Expect(kzgImpl).NotTo(BeNil())
			Expect(setupLoaded).To(BeTrue())
		})

		It("should ignore already loaded error", func() {
			kzgImpl, err := kzg.SetupKzg(kzg.SetupKzgParams{
				BlobToKzgCommitment: func(blob []byte) ([]byte, error) {
					return make([]byte, 48), nil
				},
				ComputeBlobKzgProof: func(blob, commitment []byte) ([]byte, error) {
					return make([]byte, 48), nil
				},
				LoadTrustedSetup: func(path string) error {
					return kzg.ErrTrustedSetupAlreadyLoaded
				},
			}, "/path/to/setup.txt")

			Expect(err).NotTo(HaveOccurred())
			Expect(kzgImpl).NotTo(BeNil())
		})
	})
})
