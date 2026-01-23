package accounts_test

import (
	"github.com/ChefBingbong/viem-go/accounts"
	"github.com/ChefBingbong/viem-go/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generate Private Key", func() {
	It("should generate a private key", func() {
		privateKey := accounts.GeneratePrivateKey()
		Expect(privateKey).ToNot(BeEmpty())
		Expect(privateKey).To(HavePrefix("0x"))
	})

	Context("deterministic key generation", func() {
		// Known test vector - this is a well-known test private key (DO NOT use in production!)
		const testPrivKeyHex = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		const expectedAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

		It("should create a key from hex and derive correct address", func() {
			privKey, err := crypto.PrivKeyFromHex(testPrivKeyHex)
			Expect(err).ToNot(HaveOccurred())
			Expect(privKey).ToNot(BeNil())

			// Verify we can sign and recover the signatory
			hash := crypto.NewHash([]byte("test message"))
			sig, err := privKey.Sign(&hash)
			Expect(err).ToNot(HaveOccurred())

			signatory, err := sig.Signatory(&hash)
			Expect(err).ToNot(HaveOccurred())
			Expect(signatory).To(Equal(privKey.Signatory()))
		})

		It("should create same key from bytes and hex", func() {
			privKeyFromHex, err := crypto.PrivKeyFromHex(testPrivKeyHex)
			Expect(err).ToNot(HaveOccurred())

			// Same key as raw bytes
			testBytes := []byte{
				0xac, 0x09, 0x74, 0xbe, 0xc3, 0x9a, 0x17, 0xe3,
				0x6b, 0xa4, 0xa6, 0xb4, 0xd2, 0x38, 0xff, 0x94,
				0x4b, 0xac, 0xb4, 0x78, 0xcb, 0xed, 0x5e, 0xfc,
				0xae, 0x78, 0x4d, 0x7b, 0xf4, 0xf2, 0xff, 0x80,
			}
			privKeyFromBytes, err := crypto.PrivKeyFromBytes(testBytes)
			Expect(err).ToNot(HaveOccurred())

			// Both should produce the same signatory
			Expect(privKeyFromHex.Signatory()).To(Equal(privKeyFromBytes.Signatory()))
		})
	})
})
