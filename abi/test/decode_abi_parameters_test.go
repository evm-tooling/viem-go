package abi_test

import (
	"math/big"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ethereum/go-ethereum/common"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DecodeAbiParameters", func() {
	// Test vectors from viem's decodeAbiParameters.test.ts

	Context("static types", func() {
		It("should decode blank params", func() {
			decoded, err := abi.DecodeAbiParameters([]abi.AbiParam{}, []byte{})
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(Equal([]any{}))
		})

		It("should decode uint256", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "uint256"}},
				hexToBytes("0x0000000000000000000000000000000000000000000000000000000000010f2c"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			Expect(decoded[0].(*big.Int).Cmp(big.NewInt(69420))).To(Equal(0))
		})

		It("should decode uint8", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "uint8"}},
				hexToBytes("0x0000000000000000000000000000000000000000000000000000000000000020"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			// uint8 should return as int64 (small integer)
			Expect(decoded[0]).To(BeEquivalentTo(32))
		})

		It("should decode uint32", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "uint32"}},
				hexToBytes("0x0000000000000000000000000000000000000000000000000000000000010f2c"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			Expect(decoded[0]).To(BeEquivalentTo(69420))
		})

		It("should decode int256", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "int256"}},
				hexToBytes("0x0000000000000000000000000000000000000000000000000000000000010f2c"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			Expect(decoded[0].(*big.Int).Cmp(big.NewInt(69420))).To(Equal(0))
		})

		It("should decode int256 negative (twos complement)", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "int256"}},
				hexToBytes("0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffef0d4"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			Expect(decoded[0].(*big.Int).Cmp(big.NewInt(-69420))).To(Equal(0))
		})

		It("should decode int8", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "int8"}},
				hexToBytes("0x000000000000000000000000000000000000000000000000000000000000007f"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			Expect(decoded[0]).To(BeEquivalentTo(127))
		})

		It("should decode int8 negative (twos complement)", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "int8"}},
				hexToBytes("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			Expect(decoded[0]).To(BeEquivalentTo(-128))
		})

		It("should decode int32", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "int32"}},
				hexToBytes("0x000000000000000000000000000000000000000000000000000000007fffffff"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			Expect(decoded[0]).To(BeEquivalentTo(2147483647))
		})

		It("should decode int32 negative (twos complement)", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "int32"}},
				hexToBytes("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffff80000000"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			Expect(decoded[0]).To(BeEquivalentTo(-2147483648))
		})

		It("should decode address", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "address"}},
				hexToBytes("0x00000000000000000000000014dc79964da2c08b23698b3d3cc7ca32193d9955"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			addr := decoded[0].(common.Address)
			Expect(addr.Hex()).To(Equal("0x14dC79964da2C08b23698B3D3cc7Ca32193d9955"))
		})

		It("should decode bool true", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "bool"}},
				hexToBytes("0x0000000000000000000000000000000000000000000000000000000000000001"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			Expect(decoded[0]).To(BeTrue())
		})

		It("should decode bool false", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "bool"}},
				hexToBytes("0x0000000000000000000000000000000000000000000000000000000000000000"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			Expect(decoded[0]).To(BeFalse())
		})

		It("should decode bytes8", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "bytes8"}},
				hexToBytes("0x0123456789abcdef000000000000000000000000000000000000000000000000"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			Expect(decoded[0]).To(Equal("0x0123456789abcdef"))
		})

		It("should decode bytes16", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "bytes16"}},
				hexToBytes("0x0123456789abcdef0123456789abcdef00000000000000000000000000000000"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			Expect(decoded[0]).To(Equal("0x0123456789abcdef0123456789abcdef"))
		})

		It("should decode uint256[3]", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "uint256[3]"}},
				hexToBytes("0x0000000000000000000000000000000000000000000000000000000000010f2c000000000000000000000000000000000000000000000000000000000000a45500000000000000000000000000000000000000000000000000000000190f1b44"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			arr := decoded[0].([]any)
			Expect(arr).To(HaveLen(3))
			Expect(arr[0].(*big.Int).Cmp(big.NewInt(69420))).To(Equal(0))
			Expect(arr[1].(*big.Int).Cmp(big.NewInt(42069))).To(Equal(0))
			Expect(arr[2].(*big.Int).Cmp(big.NewInt(420420420))).To(Equal(0))
		})

		It("should decode int256[3] with negative", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "int256[3]"}},
				hexToBytes("0x0000000000000000000000000000000000000000000000000000000000010f2cffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff5bab00000000000000000000000000000000000000000000000000000000190f1b44"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			arr := decoded[0].([]any)
			Expect(arr).To(HaveLen(3))
			Expect(arr[0].(*big.Int).Cmp(big.NewInt(69420))).To(Equal(0))
			Expect(arr[1].(*big.Int).Cmp(big.NewInt(-42069))).To(Equal(0))
			Expect(arr[2].(*big.Int).Cmp(big.NewInt(420420420))).To(Equal(0))
		})

		It("should decode address[2]", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "address[2]"}},
				hexToBytes("0x000000000000000000000000c961145a54c96e3ae9baa048c4f4d6b04c13916b000000000000000000000000a5cc3c03994db5b0d9a5eedd10cabab0813678ac"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			arr := decoded[0].([]any)
			Expect(arr).To(HaveLen(2))
			Expect(arr[0].(common.Address).Hex()).To(Equal("0xc961145a54C96E3aE9bAA048c4F4D6b04C13916b"))
			Expect(arr[1].(common.Address).Hex()).To(Equal("0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC"))
		})

		It("should decode bool[2]", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "bool[2]"}},
				hexToBytes("0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			arr := decoded[0].([]any)
			Expect(arr).To(HaveLen(2))
			Expect(arr[0]).To(BeTrue())
			Expect(arr[1]).To(BeFalse())
		})

		It("should decode multiple params (uint,bool,address)", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{
					{Name: "xIn", Type: "uint256"},
					{Name: "yIn", Type: "bool"},
					{Name: "zIn", Type: "address"},
				},
				hexToBytes("0x00000000000000000000000000000000000000000000000000000000000001a40000000000000000000000000000000000000000000000000000000000000001000000000000000000000000c961145a54c96e3ae9baa048c4f4d6b04c13916b"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(3))
			Expect(decoded[0].(*big.Int).Cmp(big.NewInt(420))).To(Equal(0))
			Expect(decoded[1]).To(BeTrue())
			Expect(decoded[2].(common.Address).Hex()).To(Equal("0xc961145a54C96E3aE9bAA048c4F4D6b04C13916b"))
		})
	})

	Context("dynamic types", func() {
		It("should decode string", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "string"}},
				hexToBytes("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000057761676d69000000000000000000000000000000000000000000000000000000"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			Expect(decoded[0]).To(Equal("wagmi"))
		})

		It("should decode bytes", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "bytes"}},
				hexToBytes("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000030420690000000000000000000000000000000000000000000000000000000000"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			Expect(decoded[0].([]byte)).To(Equal(hexToBytes("0x042069")))
		})

		It("should decode empty bytes", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "bytes"}},
				hexToBytes("0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			Expect(decoded[0].([]byte)).To(Equal([]byte{}))
		})

		It("should decode uint256[]", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "uint256[]"}},
				hexToBytes("0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000001a4000000000000000000000000000000000000000000000000000000000000004500000000000000000000000000000000000000000000000000000000000000160000000000000000000000000000000000000000000000000000000000000037"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			arr := decoded[0].([]any)
			Expect(arr).To(HaveLen(4))
			Expect(arr[0].(*big.Int).Cmp(big.NewInt(420))).To(Equal(0))
			Expect(arr[1].(*big.Int).Cmp(big.NewInt(69))).To(Equal(0))
			Expect(arr[2].(*big.Int).Cmp(big.NewInt(22))).To(Equal(0))
			Expect(arr[3].(*big.Int).Cmp(big.NewInt(55))).To(Equal(0))
		})

		It("should decode empty uint256[]", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "uint256[]"}},
				hexToBytes("0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			arr := decoded[0].([]any)
			Expect(arr).To(HaveLen(0))
		})

		It("should decode string[2]", func() {
			decoded, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "string[2]"}},
				hexToBytes("0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000057761676d6900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000047669656d00000000000000000000000000000000000000000000000000000000"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))
			arr := decoded[0].([]any)
			Expect(arr).To(HaveLen(2))
			Expect(arr[0]).To(Equal("wagmi"))
			Expect(arr[1]).To(Equal("viem"))
		})
	})

	Context("tuple types", func() {
		It("should decode struct (uint256,bool,address)", func() {
			decoded, err := abi.DecodeAbiParameters(
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
				hexToBytes("0x00000000000000000000000000000000000000000000000000000000000001a40000000000000000000000000000000000000000000000000000000000000001000000000000000000000000a5cc3c03994db5b0d9a5eedd10cabab0813678ac"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded).To(HaveLen(1))

			result := decoded[0].(map[string]any)
			Expect(result["x"].(*big.Int).Cmp(big.NewInt(420))).To(Equal(0))
			Expect(result["y"]).To(BeTrue())
			Expect(result["z"].(common.Address).Hex()).To(Equal("0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC"))
		})
	})

	Context("error cases", func() {
		It("should return error for zero data with params", func() {
			_, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "uint256"}},
				[]byte{},
			)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("zero data"))
		})

		It("should return error for data too small", func() {
			_, err := abi.DecodeAbiParameters(
				[]abi.AbiParam{{Type: "uint256"}},
				[]byte{0x01, 0x02, 0x03}, // Only 3 bytes
			)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("too small"))
		})
	})

	Context("roundtrip encode/decode", func() {
		It("should roundtrip uint256", func() {
			params := []abi.AbiParam{{Type: "uint256"}}
			original := big.NewInt(69420)

			encoded, err := abi.EncodeAbiParameters(params, []any{original})
			Expect(err).ToNot(HaveOccurred())

			decoded, err := abi.DecodeAbiParameters(params, encoded)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded[0].(*big.Int).Cmp(original)).To(Equal(0))
		})

		It("should roundtrip string", func() {
			params := []abi.AbiParam{{Type: "string"}}
			original := "wagmi"

			encoded, err := abi.EncodeAbiParameters(params, []any{original})
			Expect(err).ToNot(HaveOccurred())

			decoded, err := abi.DecodeAbiParameters(params, encoded)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded[0]).To(Equal(original))
		})

		It("should roundtrip multiple params", func() {
			params := []abi.AbiParam{
				{Type: "uint256"},
				{Type: "bool"},
				{Type: "address"},
			}
			addr := common.HexToAddress("0xc961145a54C96E3aE9bAA048c4F4D6b04C13916b")

			encoded, err := abi.EncodeAbiParameters(params, []any{big.NewInt(420), true, addr})
			Expect(err).ToNot(HaveOccurred())

			decoded, err := abi.DecodeAbiParameters(params, encoded)
			Expect(err).ToNot(HaveOccurred())
			Expect(decoded[0].(*big.Int).Cmp(big.NewInt(420))).To(Equal(0))
			Expect(decoded[1]).To(BeTrue())
			Expect(decoded[2].(common.Address)).To(Equal(addr))
		})
	})
})
