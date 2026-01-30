package abi_test

import (
	"math/big"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ethereum/go-ethereum/common"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DecodeErrorResult", func() {
	// Test vectors from viem's decodeErrorResult.test.ts

	Context("basic errors", func() {
		It("should decode SoldOutError()", func() {
			jsonABI := []byte(`[{"inputs":[],"name":"SoldOutError","type":"error"}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			result, err := parsed.DecodeErrorResult(hexToBytes("0x7f6df6bb"))
			Expect(err).ToNot(HaveOccurred())
			Expect(result.ErrorName).To(Equal("SoldOutError"))
			// Args can be nil or empty slice depending on go-ethereum version
			Expect(result.Args == nil || len(result.Args) == 0).To(BeTrue())
		})

		It("should decode AccessDeniedError(string)", func() {
			jsonABI := []byte(`[{"inputs":[{"name":"a","type":"string"}],"name":"AccessDeniedError","type":"error"}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			result, err := parsed.DecodeErrorResult(hexToBytes("0x83aa206e0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001a796f7520646f206e6f7420686176652061636365737320736572000000000000"))
			Expect(err).ToNot(HaveOccurred())
			Expect(result.ErrorName).To(Equal("AccessDeniedError"))
			Expect(result.Args).To(HaveLen(1))
			Expect(result.Args[0]).To(Equal("you do not have access ser"))
		})
	})

	Context("tuple parameters", func() {
		It("should decode AccessDeniedError((uint256,bool,address,uint256))", func() {
			jsonABI := []byte(`[{"inputs":[{"components":[{"name":"weight","type":"uint256"},{"name":"voted","type":"bool"},{"name":"delegate","type":"address"},{"name":"vote","type":"uint256"}],"name":"voter","type":"tuple"}],"name":"AccessDeniedError","type":"error"}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			result, err := parsed.DecodeErrorResult(hexToBytes("0x0a1895610000000000000000000000000000000000000000000000000000000000010f2c0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000a5cc3c03994db5b0d9a5eedd10cabab0813678ac0000000000000000000000000000000000000000000000000000000000000029"))
			Expect(err).ToNot(HaveOccurred())
			Expect(result.ErrorName).To(Equal("AccessDeniedError"))
			Expect(result.Args).To(HaveLen(1))
		})
	})

	Context("standard errors", func() {
		It("should decode Error(string)", func() {
			// Using nil ABI to test standard error decoding
			result, err := abi.DecodeErrorResultWithoutABI(hexToBytes("0x08c379a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000047465737400000000000000000000000000000000000000000000000000000000"))
			Expect(err).ToNot(HaveOccurred())
			Expect(result.ErrorName).To(Equal("Error"))
			Expect(result.Args).To(HaveLen(1))
			Expect(result.Args[0]).To(Equal("test"))
		})

		It("should decode Panic(uint256)", func() {
			result, err := abi.DecodeErrorResultWithoutABI(hexToBytes("0x4e487b710000000000000000000000000000000000000000000000000000000000000001"))
			Expect(err).ToNot(HaveOccurred())
			Expect(result.ErrorName).To(Equal("Panic"))
			Expect(result.Args).To(HaveLen(1))
			Expect(result.Args[0].(*big.Int).Cmp(big.NewInt(1))).To(Equal(0))
		})
	})

	Context("error cases", func() {
		It("should return error for zero data", func() {
			jsonABI := []byte(`[{"inputs":[],"name":"SoldOutError","type":"error"}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			_, err = parsed.DecodeErrorResult([]byte{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("zero data"))
		})

		It("should return error for unknown selector", func() {
			jsonABI := []byte(`[{"inputs":[],"name":"SoldOutError","type":"error"}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			_, err = parsed.DecodeErrorResult(hexToBytes("0xa3741467"))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unknown error selector"))
		})
	})

	Context("roundtrip encode/decode", func() {
		It("should roundtrip SoldOutError()", func() {
			jsonABI := []byte(`[{"inputs":[],"name":"SoldOutError","type":"error"}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			encoded, err := parsed.EncodeErrorResult("SoldOutError")
			Expect(err).ToNot(HaveOccurred())

			result, err := parsed.DecodeErrorResult(encoded)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.ErrorName).To(Equal("SoldOutError"))
		})

		It("should roundtrip AccessDeniedError(string)", func() {
			jsonABI := []byte(`[{"inputs":[{"name":"a","type":"string"}],"name":"AccessDeniedError","type":"error"}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			encoded, err := parsed.EncodeErrorResult("AccessDeniedError", "test message")
			Expect(err).ToNot(HaveOccurred())

			result, err := parsed.DecodeErrorResult(encoded)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.ErrorName).To(Equal("AccessDeniedError"))
			Expect(result.Args[0]).To(Equal("test message"))
		})

		It("should roundtrip with tuple", func() {
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

			result, err := parsed.DecodeErrorResult(encoded)
			Expect(err).ToNot(HaveOccurred())
			Expect(result.ErrorName).To(Equal("AccessDeniedError"))
		})
	})
})
