package crypto_test

import (
	"encoding/json"
	"testing/quick"

	"github.com/ChefBingbong/viem-go/crypto"
	"github.com/renproject/surge"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Private keys", func() {
	Context("when signing hashes", func() {
		Context("when verifying signatures", func() {
			It("should return the expected pubkey", func() {
				f := func(data []byte) bool {
					hash := crypto.NewHash(data)
					privKey := crypto.NewPrivKey()
					sig, err := privKey.Sign(&hash)
					Expect(err).ToNot(HaveOccurred())
					expected := privKey.Signatory()
					got, err := sig.Signatory(&hash)
					Expect(err).ToNot(HaveOccurred())
					Expect(got).To(Equal(expected))
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})
		})
	})

	Context("when marshal and then unmarshaling using binary", func() {
		It("should equal itself", func() {
			f := func() bool {
				privKey := crypto.NewPrivKey()
				marshaled, err := surge.ToBinary(privKey)
				Expect(err).ToNot(HaveOccurred())
				unmarshaled := crypto.PrivKey{}
				err = surge.FromBinary(&unmarshaled, marshaled)
				Expect(err).ToNot(HaveOccurred())
				Expect(privKey.D.Cmp(unmarshaled.D)).To(Equal(0))
				Expect(privKey.X.Cmp(unmarshaled.X)).To(Equal(0))
				Expect(privKey.Y.Cmp(unmarshaled.Y)).To(Equal(0))
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})

	Context("when marshal and then unmarshaling using JSON", func() {
		It("should equal itself", func() {
			f := func() bool {
				privKey := crypto.NewPrivKey()
				marshaled, err := json.Marshal(privKey)
				Expect(err).ToNot(HaveOccurred())
				unmarshaled := crypto.PrivKey{}
				err = json.Unmarshal(marshaled, &unmarshaled)
				Expect(err).ToNot(HaveOccurred())
				Expect(privKey.D.Cmp(unmarshaled.D)).To(Equal(0))
				Expect(privKey.X.Cmp(unmarshaled.X)).To(Equal(0))
				Expect(privKey.Y.Cmp(unmarshaled.Y)).To(Equal(0))
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})

	Context("when unmarshaling random bytes using binary", func() {
		It("should equal return an error", func() {
			f := func(data []byte) bool {
				if len(data) >= 32 {
					return true
				}
				unmarshaled := crypto.PrivKey{}
				err := surge.FromBinary(&unmarshaled, data)
				Expect(err).To(HaveOccurred())
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})
})
