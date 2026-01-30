package abi_test

import (
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ethereum/go-ethereum/common"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// hexToBytes converts a hex string to bytes, stripping 0x prefix
func hexToBytes(s string) []byte {
	s = strings.TrimPrefix(s, "0x")
	b, _ := hex.DecodeString(s)
	return b
}

// bytesToHex converts bytes to hex string with 0x prefix
func bytesToHex(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

var _ = Describe("EncodeAbiParameters", func() {
	// Test vectors from viem's encodeAbiParameters.test.ts

	Context("static types", func() {
		It("should encode blank params", func() {
			encoded, err := abi.EncodeAbiParameters([]abi.AbiParam{}, []any{})
			Expect(err).ToNot(HaveOccurred())
			Expect(encoded).To(Equal([]byte{}))
		})

		It("should encode uint256", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Type: "uint256"}},
				[]any{big.NewInt(69420)},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x0000000000000000000000000000000000000000000000000000000000010f2c"))
		})

		It("should encode uint8", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "uint8"}},
				[]any{big.NewInt(32)},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x0000000000000000000000000000000000000000000000000000000000000020"))
		})

		It("should encode uint32", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "uint32"}},
				[]any{big.NewInt(69420)},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x0000000000000000000000000000000000000000000000000000000000010f2c"))
		})

		It("should encode int256", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "int256"}},
				[]any{big.NewInt(69420)},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x0000000000000000000000000000000000000000000000000000000000010f2c"))
		})

		It("should encode int256 negative (twos complement)", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "int256"}},
				[]any{big.NewInt(-69420)},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffef0d4"))
		})

		It("should encode int8", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "int8"}},
				[]any{big.NewInt(127)},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x000000000000000000000000000000000000000000000000000000000000007f"))
		})

		It("should encode int8 negative (twos complement)", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "int8"}},
				[]any{big.NewInt(-128)},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80"))
		})

		It("should encode int32", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "int32"}},
				[]any{big.NewInt(2147483647)},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x000000000000000000000000000000000000000000000000000000007fffffff"))
		})

		It("should encode int32 negative (twos complement)", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "int32"}},
				[]any{big.NewInt(-2147483648)},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffff80000000"))
		})

		It("should encode address", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "address"}},
				[]any{common.HexToAddress("0x14dC79964da2C08b23698B3D3cc7Ca32193d9955")},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x00000000000000000000000014dc79964da2c08b23698b3d3cc7ca32193d9955"))
		})

		It("should encode address from string", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "address"}},
				[]any{"0x14dC79964da2C08b23698B3D3cc7Ca32193d9955"},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x00000000000000000000000014dc79964da2c08b23698b3d3cc7ca32193d9955"))
		})

		It("should encode bool true", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "bool"}},
				[]any{true},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x0000000000000000000000000000000000000000000000000000000000000001"))
		})

		It("should encode bool false", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "bool"}},
				[]any{false},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x0000000000000000000000000000000000000000000000000000000000000000"))
		})

		It("should encode bytes8", func() {
			var b [8]byte
			copy(b[:], hexToBytes("0x0123456789abcdef"))
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "bytes8"}},
				[]any{b},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x0123456789abcdef000000000000000000000000000000000000000000000000"))
		})

		It("should encode bytes16", func() {
			var b [16]byte
			copy(b[:], hexToBytes("0x42069420694201023210101231415122"))
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "bytes16"}},
				[]any{b},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x4206942069420102321010123141512200000000000000000000000000000000"))
		})

		It("should encode uint256[3]", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "uint256[3]"}},
				[]any{[3]*big.Int{big.NewInt(69420), big.NewInt(42069), big.NewInt(420420420)}},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x0000000000000000000000000000000000000000000000000000000000010f2c000000000000000000000000000000000000000000000000000000000000a45500000000000000000000000000000000000000000000000000000000190f1b44"))
		})

		It("should encode int256[3] with negative", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "int256[3]"}},
				[]any{[3]*big.Int{big.NewInt(69420), big.NewInt(-42069), big.NewInt(420420420)}},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x0000000000000000000000000000000000000000000000000000000000010f2cffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff5bab00000000000000000000000000000000000000000000000000000000190f1b44"))
		})

		It("should encode address[2]", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "address[2]"}},
				[]any{[2]common.Address{
					common.HexToAddress("0xc961145a54C96E3aE9bAA048c4F4D6b04C13916b"),
					common.HexToAddress("0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC"),
				}},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x000000000000000000000000c961145a54c96e3ae9baa048c4f4d6b04c13916b000000000000000000000000a5cc3c03994db5b0d9a5eedd10cabab0813678ac"))
		})

		It("should encode bool[2]", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "bool[2]"}},
				[]any{[2]bool{true, false}},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000"))
		})

		It("should encode multiple params (uint,bool,address)", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{
					{Name: "xIn", Type: "uint256"},
					{Name: "yIn", Type: "bool"},
					{Name: "zIn", Type: "address"},
				},
				[]any{
					big.NewInt(420),
					true,
					common.HexToAddress("0xc961145a54C96E3aE9bAA048c4F4D6b04C13916b"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x00000000000000000000000000000000000000000000000000000000000001a40000000000000000000000000000000000000000000000000000000000000001000000000000000000000000c961145a54c96e3ae9baa048c4f4d6b04c13916b"))
		})
	})

	Context("dynamic types", func() {
		It("should encode string", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xOut", Type: "string"}},
				[]any{"wagmi"},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000057761676d69000000000000000000000000000000000000000000000000000000"))
		})

		It("should encode string with uint and bool", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{
					{Name: "xIn", Type: "string"},
					{Name: "yIn", Type: "uint256"},
					{Name: "zIn", Type: "bool"},
				},
				[]any{"wagmi", big.NewInt(420), true},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000001a4000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000057761676d69000000000000000000000000000000000000000000000000000000"))
		})

		It("should encode bytes", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "bytes"}},
				[]any{hexToBytes("0x042069")},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000030420690000000000000000000000000000000000000000000000000000000000"))
		})

		It("should encode empty bytes", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Type: "bytes"}},
				[]any{[]byte{}},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000"))
		})

		It("should encode uint256[]", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "uint256[]"}},
				[]any{[]*big.Int{big.NewInt(420), big.NewInt(69), big.NewInt(22), big.NewInt(55)}},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000001a4000000000000000000000000000000000000000000000000000000000000004500000000000000000000000000000000000000000000000000000000000000160000000000000000000000000000000000000000000000000000000000000037"))
		})

		It("should encode empty uint256[]", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "uint256[]"}},
				[]any{[]*big.Int{}},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000"))
		})

		It("should encode string[2]", func() {
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Name: "xIn", Type: "string[2]"}},
				[]any{[2]string{"wagmi", "viem"}},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000057761676d6900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000047669656d00000000000000000000000000000000000000000000000000000000"))
		})
	})

	Context("tuple types", func() {
		It("should encode struct (uint256,bool,address)", func() {
			// go-ethereum expects a struct for tuple types
			type FooStruct struct {
				X *big.Int
				Y bool
				Z common.Address
			}
			encoded, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{
					{
						Name: "fooIn",
						Type: "tuple",
						Components: []abi.AbiParam{
							{Name: "x", Type: "uint256"},
							{Name: "y", Type: "bool"},
							{Name: "z", Type: "address"},
						},
					},
				},
				[]any{
					FooStruct{
						X: big.NewInt(420),
						Y: true,
						Z: common.HexToAddress("0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC"),
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesToHex(encoded)).To(Equal("0x00000000000000000000000000000000000000000000000000000000000001a40000000000000000000000000000000000000000000000000000000000000001000000000000000000000000a5cc3c03994db5b0d9a5eedd10cabab0813678ac"))
		})
	})

	Context("error cases", func() {
		It("should return error for params/values length mismatch", func() {
			_, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Type: "uint256"}},
				[]any{big.NewInt(1), big.NewInt(2)},
			)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("length mismatch"))
		})

		It("should return error for invalid address", func() {
			_, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Type: "address"}},
				[]any{"0x111"}, // too short
			)
			Expect(err).To(HaveOccurred())
		})

		It("should return error for invalid bool", func() {
			_, err := abi.EncodeAbiParameters(
				[]abi.AbiParam{{Type: "bool"}},
				[]any{"true"}, // string instead of bool
			)
			Expect(err).To(HaveOccurred())
		})
	})
})
