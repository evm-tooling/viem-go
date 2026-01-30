package abi_test

import (
	"math/big"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ethereum/go-ethereum/common"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PrepareEncodeFunctionData", func() {
	Context("basic preparation", func() {
		It("should prepare transfer function", func() {
			jsonABI := []byte(`[{"type":"function","name":"transfer","inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"name":"","type":"bool"}]}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			prepared, err := parsed.PrepareEncodeFunctionData("transfer",
				common.HexToAddress("0x1234567890123456789012345678901234567890"),
				big.NewInt(1000),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(prepared.FunctionName).To(Equal("transfer"))
			Expect(prepared.FunctionSelector).To(Equal([4]byte{0xa9, 0x05, 0x9c, 0xbb}))
			Expect(prepared.Abi).To(HaveLen(1))
			Expect(prepared.Abi[0].Name).To(Equal("transfer"))
		})

		It("should prepare function with no args", func() {
			jsonABI := []byte(`[{"type":"function","name":"totalSupply","inputs":[],"outputs":[{"name":"","type":"uint256"}]}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			prepared, err := parsed.PrepareEncodeFunctionData("totalSupply")
			Expect(err).ToNot(HaveOccurred())
			Expect(prepared.FunctionName).To(Equal("totalSupply"))
		})

		It("should return error for non-existent function", func() {
			jsonABI := []byte(`[{"type":"function","name":"transfer","inputs":[]}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			_, err = parsed.PrepareEncodeFunctionData("nonexistent")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})
	})

	Context("function with args", func() {
		// Note: go-ethereum's ABI parser doesn't support true function overloading
		// (same name, different signatures). It keeps only one function per name.
		It("should prepare function that takes args", func() {
			jsonABI := []byte(`[
				{"type":"function","name":"transfer","inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[]}
			]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			// Prepare transfer with args
			prepared, err := parsed.PrepareEncodeFunctionData("transfer",
				common.HexToAddress("0x1234567890123456789012345678901234567890"),
				big.NewInt(123),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(prepared.Abi[0].Inputs).To(HaveLen(2))
		})
	})

	Context("PrepareEncodeFunctionDataBySelector", func() {
		It("should prepare by selector", func() {
			jsonABI := []byte(`[{"type":"function","name":"transfer","inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"name":"","type":"bool"}]}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			selector := [4]byte{0xa9, 0x05, 0x9c, 0xbb}
			prepared, err := parsed.PrepareEncodeFunctionDataBySelector(selector)
			Expect(err).ToNot(HaveOccurred())
			Expect(prepared.FunctionName).To(Equal("transfer"))
		})

		It("should return error for unknown selector", func() {
			jsonABI := []byte(`[{"type":"function","name":"transfer","inputs":[]}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			unknownSelector := [4]byte{0xde, 0xad, 0xbe, 0xef}
			_, err = parsed.PrepareEncodeFunctionDataBySelector(unknownSelector)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("EncodeWithPrepared", func() {
		It("should encode using prepared data", func() {
			jsonABI := []byte(`[{"type":"function","name":"transfer","inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"name":"","type":"bool"}]}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			to := common.HexToAddress("0x1234567890123456789012345678901234567890")
			amount := big.NewInt(1000)

			prepared, err := parsed.PrepareEncodeFunctionData("transfer", to, amount)
			Expect(err).ToNot(HaveOccurred())

			// Encode using prepared
			encoded, err := parsed.EncodeWithPrepared(prepared, to, amount)
			Expect(err).ToNot(HaveOccurred())

			// Compare with direct encoding
			direct, err := parsed.EncodeFunctionData("transfer", to, amount)
			Expect(err).ToNot(HaveOccurred())

			Expect(encoded).To(Equal(direct))
		})

		It("should return error for nil prepared data", func() {
			jsonABI := []byte(`[{"type":"function","name":"transfer","inputs":[]}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			_, err = parsed.EncodeWithPrepared(nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid"))
		})
	})
})
