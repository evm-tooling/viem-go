package abi_test

import (
	"math/big"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ethereum/go-ethereum/common"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EncodePacked", func() {
	// Test vectors from viem's encodePacked.test.ts

	Context("address", func() {
		It("should encode address", func() {
			encoded, err := abi.EncodePacked(
				[]string{"address"},
				[]any{common.HexToAddress("0x14dC79964da2C08b23698B3D3cc7Ca32193d9955")},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x14dc79964da2c08b23698b3d3cc7ca32193d9955"))
		})

		It("should encode address from string", func() {
			encoded, err := abi.EncodePacked(
				[]string{"address"},
				[]any{"0x14dC79964da2C08b23698B3D3cc7Ca32193d9955"},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x14dc79964da2c08b23698b3d3cc7ca32193d9955"))
		})
	})

	Context("string", func() {
		It("should encode string", func() {
			encoded, err := abi.EncodePacked(
				[]string{"string"},
				[]any{"wagmi"},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x7761676d69"))
		})

		It("should encode empty string", func() {
			encoded, err := abi.EncodePacked(
				[]string{"string"},
				[]any{""},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(encoded).To(Equal([]byte{}))
		})
	})

	Context("bool", func() {
		It("should encode bool true", func() {
			encoded, err := abi.EncodePacked(
				[]string{"bool"},
				[]any{true},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x01"))
		})

		It("should encode bool false", func() {
			encoded, err := abi.EncodePacked(
				[]string{"bool"},
				[]any{false},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x00"))
		})
	})

	Context("integers", func() {
		It("should encode uint8", func() {
			encoded, err := abi.EncodePacked(
				[]string{"uint8"},
				[]any{uint8(32)},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x20"))
		})

		It("should encode uint16", func() {
			encoded, err := abi.EncodePacked(
				[]string{"uint16"},
				[]any{uint16(420)},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x01a4"))
		})

		It("should encode uint32", func() {
			encoded, err := abi.EncodePacked(
				[]string{"uint32"},
				[]any{uint32(69420)},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x00010f2c"))
		})

		It("should encode uint256", func() {
			encoded, err := abi.EncodePacked(
				[]string{"uint256"},
				[]any{big.NewInt(69420)},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x0000000000000000000000000000000000000000000000000000000000010f2c"))
		})

		It("should encode int8 negative", func() {
			encoded, err := abi.EncodePacked(
				[]string{"int8"},
				[]any{int8(-1)},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0xff"))
		})

		It("should encode int32 negative", func() {
			encoded, err := abi.EncodePacked(
				[]string{"int32"},
				[]any{int32(-69420)},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0xfffef0d4"))
		})
	})

	Context("bytes", func() {
		It("should encode bytes", func() {
			encoded, err := abi.EncodePacked(
				[]string{"bytes"},
				[]any{hexToBytes("0x1234")},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x1234"))
		})

		It("should encode bytes4", func() {
			encoded, err := abi.EncodePacked(
				[]string{"bytes4"},
				[]any{hexToBytes("0x12345678")},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x12345678"))
		})

		It("should encode bytes32", func() {
			encoded, err := abi.EncodePacked(
				[]string{"bytes32"},
				[]any{hexToBytes("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"))
		})
	})

	Context("multiple types", func() {
		It("should encode address and uint256", func() {
			encoded, err := abi.EncodePacked(
				[]string{"address", "uint256"},
				[]any{
					common.HexToAddress("0x14dC79964da2C08b23698B3D3cc7Ca32193d9955"),
					big.NewInt(420),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x14dc79964da2c08b23698b3d3cc7ca32193d995500000000000000000000000000000000000000000000000000000000000001a4"))
		})

		It("should encode uint8 and uint8", func() {
			encoded, err := abi.EncodePacked(
				[]string{"uint8", "uint8"},
				[]any{uint8(1), uint8(2)},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x0102"))
		})

		It("should encode string and uint256", func() {
			encoded, err := abi.EncodePacked(
				[]string{"string", "uint256"},
				[]any{"wagmi", big.NewInt(420)},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x7761676d6900000000000000000000000000000000000000000000000000000000000001a4"))
		})

		It("should encode bool, address, string", func() {
			encoded, err := abi.EncodePacked(
				[]string{"bool", "address", "string"},
				[]any{
					true,
					common.HexToAddress("0x14dC79964da2C08b23698B3D3cc7Ca32193d9955"),
					"hello",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x0114dc79964da2c08b23698b3d3cc7ca32193d995568656c6c6f"))
		})
	})

	Context("error cases", func() {
		It("should return error for length mismatch", func() {
			_, err := abi.EncodePacked(
				[]string{"uint256"},
				[]any{big.NewInt(1), big.NewInt(2)},
			)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("length mismatch"))
		})

		It("should return error for invalid address", func() {
			_, err := abi.EncodePacked(
				[]string{"address"},
				[]any{"0x123"}, // invalid
			)
			Expect(err).To(HaveOccurred())
		})

		It("should return error for bytes size mismatch", func() {
			_, err := abi.EncodePacked(
				[]string{"bytes4"},
				[]any{hexToBytes("0x12")}, // only 1 byte
			)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("size mismatch"))
		})
	})
})
