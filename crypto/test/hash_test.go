package crypto_test

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"testing"
	"testing/quick"

	json "github.com/goccy/go-json"

	"github.com/renproject/surge"

	"github.com/ChefBingbong/viem-go/crypto"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hashes", func() {
	Context("when hashing data", func() {
		It("should return the SHA2 256-bit hash of the data", func() {
			f := func(data []byte) bool {
				expected := crypto.Hash(sha256.Sum256(data))
				got := crypto.NewHash(data)
				Expect(got.Equal(&expected)).To(BeTrue())
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})

	Context("when marshaling and then unmarshaling using binary", func() {
		It("should equal itself", func() {
			f := func(data []byte) bool {
				hash := crypto.NewHash(data)
				marshaled, err := surge.ToBinary(hash)
				Expect(err).ToNot(HaveOccurred())
				unmarshaled := crypto.Hash{}
				err = surge.FromBinary(&unmarshaled, marshaled)
				Expect(err).ToNot(HaveOccurred())
				Expect(hash.Equal(&unmarshaled)).To(BeTrue())
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})

	Context("when unmarshaling random bytes using binary", func() {
		It("should equal reutrn an error", func() {
			f := func(data []byte) bool {
				if len(data) >= 32 {
					return true
				}
				unmarshaled := crypto.Hash{}
				err := surge.FromBinary(&unmarshaled, data)
				Expect(err).To(HaveOccurred())
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})

	Context("when marshaling and then unmarshaling using JSON", func() {
		It("should equal itself", func() {
			f := func(data []byte) bool {
				hash := crypto.NewHash(data)
				marshaled, err := hash.MarshalJSON()
				Expect(err).ToNot(HaveOccurred())
				unmarshaled := crypto.Hash{}
				err = unmarshaled.UnmarshalJSON(marshaled)
				Expect(err).ToNot(HaveOccurred())
				Expect(hash.Equal(&unmarshaled)).To(BeTrue())
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})

		It("should equal its string representation", func() {
			f := func(data []byte) bool {
				hash := crypto.NewHash(data)
				got, err := hash.MarshalJSON()
				Expect(err).ToNot(HaveOccurred())
				expected, err := json.Marshal(hash.String())
				Expect(err).ToNot(HaveOccurred())
				Expect(bytes.Equal(got, expected)).To(BeTrue())
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})

	Context("when unmarshaling random bytes using JSON", func() {
		It("should equal reutrn an error", func() {
			f := func(data []byte) bool {
				unmarshaled := crypto.Hash{}
				err := unmarshaled.UnmarshalJSON(data)
				Expect(err).To(HaveOccurred())
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})

	Context("when computing the merkle hash", func() {
		Context("when using the safe and unsafe implementations", func() {
			Context("when computing the merkle hash of zero hashes", func() {
				It("should return the empty hash", func() {
					f := func() bool {
						hashes := make([]crypto.Hash, 0)
						rootHash := crypto.NewMerkleHash(hashes)
						safeRootHash := crypto.NewMerkleHashSafe(hashes)
						Expect(rootHash).To(Equal(crypto.Hash{}))
						Expect(safeRootHash).To(Equal(crypto.Hash{}))
						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})

			Context("when computing the merkle hash of one hash", func() {
				It("should return the input hash", func() {
					f := func() bool {
						hashes := make([]crypto.Hash, 1)
						for i := range hashes {
							rand.Read(hashes[i][:])
						}
						rootHash := crypto.NewMerkleHash(hashes)
						safeRootHash := crypto.NewMerkleHashSafe(hashes)
						Expect(rootHash).To(Equal(hashes[0]))
						Expect(safeRootHash).To(Equal(hashes[0]))
						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})

			Context("when computing the merkle hash of two hashes", func() {
				It("should return the same merkle hash", func() {
					f := func() bool {
						hashes := make([]crypto.Hash, 2)
						for i := range hashes {
							rand.Read(hashes[i][:])
						}
						expectedHash := crypto.Hash(sha256.Sum256(append(hashes[0][:], hashes[1][:]...)))
						rootHash := crypto.NewMerkleHash(hashes)
						safeRootHash := crypto.NewMerkleHashSafe(hashes)
						Expect(rootHash).To(Equal(expectedHash))
						Expect(safeRootHash).To(Equal(expectedHash))
						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})

			Context("when computing the merkle hash of three hash", func() {
				It("should return the same merkle hash", func() {
					f := func() bool {
						hashes := make([]crypto.Hash, 3)
						for i := range hashes {
							rand.Read(hashes[i][:])
						}
						expectedHash := crypto.Hash(sha256.Sum256(append(hashes[1][:], hashes[2][:]...)))
						expectedHash = crypto.Hash(sha256.Sum256(append(hashes[0][:], expectedHash[:]...)))
						rootHash := crypto.NewMerkleHash(hashes)
						safeRootHash := crypto.NewMerkleHashSafe(hashes)
						Expect(rootHash).To(Equal(expectedHash))
						Expect(safeRootHash).To(Equal(expectedHash))
						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})

			It("should return the same merkle hash", func() {
				f := func(n uint) bool {
					n = n % 1000
					hashes := make([]crypto.Hash, n)
					for i := range hashes {
						rand.Read(hashes[i][:])
					}
					rootHash := crypto.NewMerkleHash(hashes)
					safeRootHash := crypto.NewMerkleHashSafe(hashes)
					Expect(rootHash).To(Equal(safeRootHash))
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})
		})

		Context("when using signatories", func() {
			It("should equal the hash equivalent", func() {
				f := func(n uint) bool {
					n = n % 1000
					hashes := make([]crypto.Hash, n)
					for i := range hashes {
						rand.Read(hashes[i][:])
					}
					signatories := make([]crypto.Signatory, n)
					for i := range signatories {
						signatories[i] = crypto.Signatory(hashes[i])
					}
					expected := crypto.NewMerkleHash(hashes)
					got := crypto.NewMerkleHashFromSignatories(signatories)
					Expect(expected).To(Equal(got))
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})
		})
	})
})

func BenchmarkMerkleHash(b *testing.B) {
	hashes := make([]crypto.Hash, b.N)
	for i := range hashes {
		rand.Read(hashes[i][:])
	}
	b.ResetTimer()
	b.ReportAllocs()
	crypto.NewMerkleHash(hashes)
}

func BenchmarkMerkleHashSafe(b *testing.B) {
	hashes := make([]crypto.Hash, b.N)
	for i := range hashes {
		rand.Read(hashes[i][:])
	}
	b.ResetTimer()
	b.ReportAllocs()
	crypto.NewMerkleHashSafe(hashes)
}

func BenchmarkMerkleHash1000(b *testing.B) {
	hashes := make([]crypto.Hash, 1000)
	for i := range hashes {
		rand.Read(hashes[i][:])
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		crypto.NewMerkleHash(hashes)
	}
}

func BenchmarkMerkleHashSafe1000(b *testing.B) {
	hashes := make([]crypto.Hash, 1000)
	for i := range hashes {
		rand.Read(hashes[i][:])
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		crypto.NewMerkleHashSafe(hashes)
	}
}
