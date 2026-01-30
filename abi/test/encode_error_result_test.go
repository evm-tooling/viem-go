package abi_test

import (
	"math/big"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ethereum/go-ethereum/common"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EncodeErrorResult", func() {
	// Test vectors from viem's encodeErrorResult.test.ts

	Context("basic errors", func() {
		It("should encode SoldOutError()", func() {
			jsonABI := []byte(`[{"inputs":[],"name":"SoldOutError","type":"error"}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			encoded, err := parsed.EncodeErrorResult("SoldOutError")
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x7f6df6bb"))
		})

		It("should encode AccessDeniedError(string)", func() {
			jsonABI := []byte(`[{"inputs":[{"name":"a","type":"string"}],"name":"AccessDeniedError","type":"error"}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			encoded, err := parsed.EncodeErrorResult("AccessDeniedError", "you do not have access ser")
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x83aa206e0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001a796f7520646f206e6f7420686176652061636365737320736572000000000000"))
		})
	})

	Context("tuple parameters", func() {
		It("should encode AccessDeniedError((uint256,bool,address,uint256))", func() {
			jsonABI := []byte(`[{"inputs":[{"components":[{"name":"weight","type":"uint256"},{"name":"voted","type":"bool"},{"name":"delegate","type":"address"},{"name":"vote","type":"uint256"}],"name":"voter","type":"tuple"}],"name":"AccessDeniedError","type":"error"}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

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

			encoded, err := parsed.EncodeErrorResult("AccessDeniedError", voter)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x0a1895610000000000000000000000000000000000000000000000000000000000010f2c0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000a5cc3c03994db5b0d9a5eedd10cabab0813678ac0000000000000000000000000000000000000000000000000000000000000029"))
		})
	})

	Context("error cases", func() {
		It("should return error for non-existent error", func() {
			jsonABI := []byte(`[{"inputs":[],"name":"SoldOutError","type":"error"}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			_, err = parsed.EncodeErrorResult("AccessDeniedError")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})

		It("should return error for error with no inputs but args provided", func() {
			jsonABI := []byte(`[{"name":"AccessDeniedError","type":"error"}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			_, err = parsed.EncodeErrorResult("AccessDeniedError", "some arg")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("takes no arguments"))
		})
	})
})
