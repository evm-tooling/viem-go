package abi_test

import (
	"math/big"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ethereum/go-ethereum/common"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EncodeFunctionData", func() {
	// Test vectors from viem's encodeFunctionData.test.ts

	Context("basic functions", func() {
		It("should encode foo()", func() {
			jsonABI := []byte(`[{"inputs":[],"name":"foo","outputs":[],"stateMutability":"nonpayable","type":"function"}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			encoded, err := parsed.EncodeFunctionData("foo")
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0xc2985578"))
		})

		It("should encode bar(uint256)", func() {
			jsonABI := []byte(`[{"inputs":[{"name":"a","type":"uint256"}],"name":"bar","outputs":[],"stateMutability":"nonpayable","type":"function"}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			encoded, err := parsed.EncodeFunctionData("bar", big.NewInt(1))
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x0423a1320000000000000000000000000000000000000000000000000000000000000001"))
		})

		It("should encode transfer(address,uint256)", func() {
			jsonABI := []byte(`[{"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			to := common.HexToAddress("0x0000000000000000000000000000000000000000")
			amount := big.NewInt(69420)

			encoded, err := parsed.EncodeFunctionData("transfer", to, amount)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0xa9059cbb00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010f2c"))
		})
	})

	Context("tuple parameters", func() {
		It("should encode getVoter((uint256,bool,address,uint256))", func() {
			jsonABI := []byte(`[{"inputs":[{"components":[{"name":"weight","type":"uint256"},{"name":"voted","type":"bool"},{"name":"delegate","type":"address"},{"name":"vote","type":"uint256"}],"name":"voter","type":"tuple"}],"name":"getVoter","outputs":[],"stateMutability":"nonpayable","type":"function"}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			// Create the voter struct as a tuple
			voter := struct {
				Weight   *big.Int
				Voted    bool
				Delegate common.Address
				Vote     *big.Int
			}{
				Weight:   big.NewInt(69420),
				Voted:    true,
				Delegate: common.HexToAddress("0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC"),
				Vote:     big.NewInt(41),
			}

			encoded, err := parsed.EncodeFunctionData("getVoter", voter)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0xf37414670000000000000000000000000000000000000000000000000000000000010f2c0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000a5cc3c03994db5b0d9a5eedd10cabab0813678ac0000000000000000000000000000000000000000000000000000000000000029"))
		})
	})

	Context("error cases", func() {
		It("should return error for non-existent function", func() {
			jsonABI := []byte(`[{"inputs":[],"name":"foo","outputs":[],"stateMutability":"nonpayable","type":"function"}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			_, err = parsed.EncodeFunctionData("bar")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})
	})

	Context("aliases", func() {
		It("should work with Pack alias", func() {
			jsonABI := []byte(`[{"inputs":[],"name":"foo","outputs":[],"stateMutability":"nonpayable","type":"function"}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			encoded1, err := parsed.EncodeFunctionData("foo")
			Expect(err).ToNot(HaveOccurred())

			encoded2, err := parsed.Pack("foo")
			Expect(err).ToNot(HaveOccurred())

			Expect(encoded1).To(Equal(encoded2))
		})

		It("should work with EncodeCall alias", func() {
			jsonABI := []byte(`[{"inputs":[],"name":"foo","outputs":[],"stateMutability":"nonpayable","type":"function"}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			encoded1, err := parsed.EncodeFunctionData("foo")
			Expect(err).ToNot(HaveOccurred())

			encoded2, err := parsed.EncodeCall("foo")
			Expect(err).ToNot(HaveOccurred())

			Expect(encoded1).To(Equal(encoded2))
		})
	})
})
