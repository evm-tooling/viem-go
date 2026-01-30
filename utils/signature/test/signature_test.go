package test

import (
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ChefBingbong/viem-go/utils/signature"
)

var _ = Describe("Signature", func() {
	Describe("ToPrefixedMessage", func() {
		It("should prefix a string message", func() {
			msg := signature.NewSignableMessage("hello world")
			result := signature.ToPrefixedMessage(msg)
			// \x19Ethereum Signed Message:\n11hello world
			Expect(result).To(HavePrefix("0x"))
			Expect(len(result)).To(BeNumerically(">", 4))
		})

		It("should prefix a raw hex message", func() {
			msg := signature.NewSignableMessageRawHex("0x68656c6c6f")
			result := signature.ToPrefixedMessage(msg)
			Expect(result).To(HavePrefix("0x"))
		})
	})

	Describe("HashMessage", func() {
		It("should hash a message with Ethereum prefix", func() {
			msg := signature.NewSignableMessage("hello world")
			hash := signature.HashMessage(msg)
			// This is the expected hash for "hello world" with Ethereum prefix
			Expect(hash).To(Equal("0xd9eba16ed0ecae432b71fe008c98cc872bb4cc214d3220a36f365326cf807d68"))
		})
	})

	Describe("ParseSignature", func() {
		It("should parse a valid 65-byte signature", func() {
			sigHex := "0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db81c"
			sig, err := signature.ParseSignature(sigHex)
			Expect(err).NotTo(HaveOccurred())
			Expect(sig.R).To(Equal("0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf"))
			Expect(sig.S).To(Equal("0x4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db8"))
			Expect(sig.YParity).To(Equal(1))
			Expect(sig.V.Int64()).To(Equal(int64(28)))
		})

		It("should handle yParity 0", func() {
			sigHex := "0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db81b"
			sig, err := signature.ParseSignature(sigHex)
			Expect(err).NotTo(HaveOccurred())
			Expect(sig.YParity).To(Equal(0))
			Expect(sig.V.Int64()).To(Equal(int64(27)))
		})

		It("should fail for invalid length", func() {
			_, err := signature.ParseSignature("0x1234")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("SerializeSignature", func() {
		It("should serialize a signature to hex", func() {
			sig := &signature.Signature{
				R:       "0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf",
				S:       "0x4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db8",
				YParity: 1,
			}
			hex, err := signature.SerializeSignature(sig)
			Expect(err).NotTo(HaveOccurred())
			Expect(hex).To(Equal("0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db81c"))
		})
	})

	Describe("CompactSignature", func() {
		It("should convert signature to compact and back", func() {
			sig := &signature.Signature{
				R:       "0x68a020a209d3d56c46f38cc50a33f704f4a9a10a59377f8dd762ac66910e9b90",
				S:       "0x7e865ad05c4035ab5792787d4a0297a43617ae897930a6fe4d822b8faea52064",
				YParity: 0,
			}

			compact, err := signature.SignatureToCompactSignature(sig)
			Expect(err).NotTo(HaveOccurred())
			Expect(compact.R).To(Equal(sig.R))

			// Convert back
			recovered, err := signature.CompactSignatureToSignature(compact)
			Expect(err).NotTo(HaveOccurred())
			Expect(recovered.R).To(Equal(sig.R))
			Expect(recovered.YParity).To(Equal(sig.YParity))
		})

		It("should handle yParity 1 in compact format", func() {
			sig := &signature.Signature{
				R:       "0x68a020a209d3d56c46f38cc50a33f704f4a9a10a59377f8dd762ac66910e9b90",
				S:       "0x7e865ad05c4035ab5792787d4a0297a43617ae897930a6fe4d822b8faea52064",
				YParity: 1,
			}

			compact, err := signature.SignatureToCompactSignature(sig)
			Expect(err).NotTo(HaveOccurred())

			recovered, err := signature.CompactSignatureToSignature(compact)
			Expect(err).NotTo(HaveOccurred())
			Expect(recovered.YParity).To(Equal(1))
		})
	})

	Describe("IsErc6492Signature", func() {
		It("should detect ERC-6492 signature", func() {
			erc6492Sig := "0x000000000000000000000000cafebabecafebabecafebabecafebabecafebabe000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000004deadbeef000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000041a461f509887bd19e312c0c58467ce8ff8e300d3c1a90b608a760c5b80318eaf15fe57c96f9175d6cd4daad4663763baa7e78836e067d0163e9a2ccf2ff753f5b1b000000000000000000000000000000000000000000000000000000000000006492649264926492649264926492649264926492649264926492649264926492"
			Expect(signature.IsErc6492Signature(erc6492Sig)).To(BeTrue())
		})

		It("should return false for regular signature", func() {
			regularSig := "0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db81c"
			Expect(signature.IsErc6492Signature(regularSig)).To(BeFalse())
		})
	})

	Describe("HashTypedData", func() {
		It("should hash typed data correctly", func() {
			typedData := signature.TypedDataDefinition{
				Domain: signature.TypedDataDomain{
					Name:              "Ether Mail",
					Version:           "1",
					ChainId:           big.NewInt(1),
					VerifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
				},
				Types: map[string][]signature.TypedDataField{
					"Person": {
						{Name: "name", Type: "string"},
						{Name: "wallet", Type: "address"},
					},
					"Mail": {
						{Name: "from", Type: "Person"},
						{Name: "to", Type: "Person"},
						{Name: "contents", Type: "string"},
					},
				},
				PrimaryType: "Mail",
				Message: map[string]any{
					"from": map[string]any{
						"name":   "Cow",
						"wallet": "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
					},
					"to": map[string]any{
						"name":   "Bob",
						"wallet": "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB",
					},
					"contents": "Hello, Bob!",
				},
			}

			hash, err := signature.HashTypedData(typedData)
			Expect(err).NotTo(HaveOccurred())
			Expect(hash).To(HavePrefix("0x"))
			Expect(len(hash)).To(Equal(66)) // 0x + 64 hex chars
		})
	})

	Describe("EncodeType", func() {
		It("should encode type string correctly", func() {
			types := map[string][]signature.TypedDataField{
				"Person": {
					{Name: "name", Type: "string"},
					{Name: "wallet", Type: "address"},
				},
				"Mail": {
					{Name: "from", Type: "Person"},
					{Name: "to", Type: "Person"},
					{Name: "contents", Type: "string"},
				},
			}

			encoded := signature.EncodeType("Mail", types)
			Expect(encoded).To(Equal("Mail(Person from,Person to,string contents)Person(string name,address wallet)"))
		})
	})
})
