package crypto_test

import (
	"bytes"
	"testing/quick"

	json "github.com/goccy/go-json"

	"github.com/renproject/surge"

	"github.com/ChefBingbong/viem-go/crypto"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Signatures", func() {
	Context("when marshaling and then unmarshaling using binary", func() {
		It("should equal itself", func() {
			f := func(data [65]byte) bool {
				sig := crypto.Signature(data)
				marshaled, err := surge.ToBinary(sig)
				Expect(err).ToNot(HaveOccurred())
				unmarshaled := crypto.Signature{}
				err = surge.FromBinary(&unmarshaled, marshaled)
				Expect(err).ToNot(HaveOccurred())
				Expect(sig.Equal(&unmarshaled)).To(BeTrue())
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})

	Context("when unmarshaling random bytes using binary", func() {
		It("should equal reutrn an error", func() {
			f := func(data []byte) bool {
				if len(data) >= 65 {
					return true
				}
				unmarshaled := crypto.Signature{}
				err := surge.FromBinary(&unmarshaled, data)
				Expect(err).To(HaveOccurred())
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})

	Context("when marshaling and then unmarshaling using JSON", func() {
		It("should equal itself", func() {
			f := func(data [65]byte) bool {
				sig := crypto.Signature(data)
				marshaled, err := sig.MarshalJSON()
				Expect(err).ToNot(HaveOccurred())
				unmarshaled := crypto.Signature{}
				err = unmarshaled.UnmarshalJSON(marshaled)
				Expect(err).ToNot(HaveOccurred())
				Expect(sig.Equal(&unmarshaled)).To(BeTrue())
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})

		It("should equal its string representation", func() {
			f := func(data [65]byte) bool {
				sig := crypto.Signature(data)
				got, err := sig.MarshalJSON()
				Expect(err).ToNot(HaveOccurred())
				expected, err := json.Marshal(sig.String())
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
				unmarshaled := crypto.Signature{}
				err := unmarshaled.UnmarshalJSON(data)
				Expect(err).To(HaveOccurred())
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})
})

var _ = Describe("Signatories", func() {
	Context("when marshaling and then unmarshaling using binary", func() {
		It("should equal itself", func() {
			f := func(data [32]byte) bool {
				sig := crypto.Signatory(data)
				marshaled, err := surge.ToBinary(sig)
				Expect(err).ToNot(HaveOccurred())
				unmarshaled := crypto.Signatory{}
				err = surge.FromBinary(&unmarshaled, marshaled)
				Expect(err).ToNot(HaveOccurred())
				Expect(sig.Equal(&unmarshaled)).To(BeTrue())
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
				unmarshaled := crypto.Signatory{}
				err := surge.FromBinary(&unmarshaled, data)
				Expect(err).To(HaveOccurred())
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})

	Context("when marshaling and then unmarshaling using JSON", func() {
		It("should equal itself", func() {
			f := func(data [32]byte) bool {
				sig := crypto.Signatory(data)
				marshaled, err := sig.MarshalJSON()
				Expect(err).ToNot(HaveOccurred())
				unmarshaled := crypto.Signatory{}
				err = unmarshaled.UnmarshalJSON(marshaled)
				Expect(err).ToNot(HaveOccurred())
				Expect(sig.Equal(&unmarshaled)).To(BeTrue())
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})

		It("should equal its string representation", func() {
			f := func(data [32]byte) bool {
				sig := crypto.Signatory(data)
				got, err := sig.MarshalJSON()
				Expect(err).ToNot(HaveOccurred())
				expected, err := json.Marshal(sig.String())
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
				unmarshaled := crypto.Signatory{}
				err := unmarshaled.UnmarshalJSON(data)
				Expect(err).To(HaveOccurred())
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})
})
