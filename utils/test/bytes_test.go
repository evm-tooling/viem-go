package utils_test

import (
	"math/big"

	"github.com/ChefBingbong/viem-go/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bytes Utils", func() {
	// ==========================================================================
	// Int64 ↔ Bytes
	// ==========================================================================
	Context("Int64 conversions", func() {
		It("should convert int64 to bytes and back", func() {
			testCases := []struct {
				value    int64
				expected []byte
			}{
				{0, []byte{0, 0, 0, 0, 0, 0, 0, 0}},
				{1, []byte{0, 0, 0, 0, 0, 0, 0, 1}},
				{255, []byte{0, 0, 0, 0, 0, 0, 0, 255}},
				{256, []byte{0, 0, 0, 0, 0, 0, 1, 0}},
				{65535, []byte{0, 0, 0, 0, 0, 0, 255, 255}},
				{1000000, []byte{0, 0, 0, 0, 0, 15, 66, 64}},
				{9223372036854775807, []byte{127, 255, 255, 255, 255, 255, 255, 255}}, // max int64
				{-1, []byte{255, 255, 255, 255, 255, 255, 255, 255}},
				{-128, []byte{255, 255, 255, 255, 255, 255, 255, 128}},
			}

			for _, tc := range testCases {
				// Int → Bytes
				bytes := utils.IntToBytes(tc.value)
				Expect(bytes).To(Equal(tc.expected), "IntToBytes(%d)", tc.value)

				// Bytes → Int
				result := utils.BytesToInt(tc.expected)
				Expect(result).To(Equal(tc.value), "BytesToInt(%v)", tc.expected)
			}
		})

		It("should convert int64 to minimal bytes", func() {
			testCases := []struct {
				value    int64
				expected []byte
			}{
				{0, []byte{0}},
				{1, []byte{1}},
				{255, []byte{255}},
				{256, []byte{1, 0}},
				{65535, []byte{255, 255}},
				{65536, []byte{1, 0, 0}},
				{16777215, []byte{255, 255, 255}},
				{16777216, []byte{1, 0, 0, 0}},
			}

			for _, tc := range testCases {
				bytes := utils.IntToBytesMinimal(tc.value)
				Expect(bytes).To(Equal(tc.expected), "IntToBytesMinimal(%d)", tc.value)
			}
		})

		It("should handle short byte slices for BytesToInt", func() {
			// Single byte
			Expect(utils.BytesToInt([]byte{42})).To(Equal(int64(42)))
			// Two bytes
			Expect(utils.BytesToInt([]byte{1, 0})).To(Equal(int64(256)))
			// Four bytes
			Expect(utils.BytesToInt([]byte{0, 0, 0, 1})).To(Equal(int64(1)))
		})
	})

	// ==========================================================================
	// Uint64 ↔ Bytes
	// ==========================================================================
	Context("Uint64 conversions", func() {
		It("should convert uint64 to bytes and back", func() {
			testCases := []struct {
				value    uint64
				expected []byte
			}{
				{0, []byte{0, 0, 0, 0, 0, 0, 0, 0}},
				{1, []byte{0, 0, 0, 0, 0, 0, 0, 1}},
				{255, []byte{0, 0, 0, 0, 0, 0, 0, 255}},
				{256, []byte{0, 0, 0, 0, 0, 0, 1, 0}},
				{18446744073709551615, []byte{255, 255, 255, 255, 255, 255, 255, 255}}, // max uint64
			}

			for _, tc := range testCases {
				// Uint → Bytes
				bytes := utils.UintToBytes(tc.value)
				Expect(bytes).To(Equal(tc.expected), "UintToBytes(%d)", tc.value)

				// Bytes → Uint
				result := utils.BytesToUint(tc.expected)
				Expect(result).To(Equal(tc.value), "BytesToUint(%v)", tc.expected)
			}
		})
	})

	// ==========================================================================
	// BigInt ↔ Bytes
	// ==========================================================================
	Context("BigInt conversions", func() {
		It("should convert big.Int to bytes and back", func() {
			testCases := []struct {
				value    string // decimal string representation
				expected []byte
			}{
				{"0", []byte{}},
				{"1", []byte{1}},
				{"255", []byte{255}},
				{"256", []byte{1, 0}},
				{"65535", []byte{255, 255}},
				{"1000000000000000000", []byte{13, 224, 182, 179, 167, 100, 0, 0}}, // 1 ETH in wei
				// Large number: 2^256 - 1 (max uint256)
				{"115792089237316195423570985008687907853269984665640564039457584007913129639935",
					[]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
						255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}},
			}

			for _, tc := range testCases {
				value, _ := new(big.Int).SetString(tc.value, 10)

				// BigInt → Bytes
				bytes := utils.BigIntToBytes(value)
				Expect(bytes).To(Equal(tc.expected), "BigIntToBytes(%s)", tc.value)

				// Bytes → BigInt
				result := utils.BytesToBigInt(tc.expected)
				Expect(result.Cmp(value)).To(Equal(0), "BytesToBigInt(%v)", tc.expected)
			}
		})

		It("should pad big.Int to specified size", func() {
			value := big.NewInt(256) // 0x100

			// Pad to 4 bytes
			padded := utils.BigIntToBytesPadded(value, 4)
			Expect(padded).To(Equal([]byte{0, 0, 1, 0}))

			// Pad to 32 bytes (common for Ethereum)
			padded32 := utils.BigIntToBytesPadded(value, 32)
			Expect(len(padded32)).To(Equal(32))
			Expect(padded32[30]).To(Equal(byte(1)))
			Expect(padded32[31]).To(Equal(byte(0)))
		})

		It("should handle nil big.Int", func() {
			result := utils.BigIntToBytes(nil)
			Expect(result).To(Equal([]byte{0}))
		})
	})

	// ==========================================================================
	// Bool ↔ Bytes
	// ==========================================================================
	Context("Bool conversions", func() {
		It("should convert bool to bytes", func() {
			Expect(utils.BoolToBytes(true)).To(Equal([]byte{1}))
			Expect(utils.BoolToBytes(false)).To(Equal([]byte{0}))
		})

		It("should convert bytes to bool", func() {
			// True cases - any non-zero byte
			Expect(utils.BytesToBool([]byte{1})).To(BeTrue())
			Expect(utils.BytesToBool([]byte{255})).To(BeTrue())
			Expect(utils.BytesToBool([]byte{0, 0, 1})).To(BeTrue())
			Expect(utils.BytesToBool([]byte{0, 1, 0})).To(BeTrue())

			// False cases - all zeros
			Expect(utils.BytesToBool([]byte{0})).To(BeFalse())
			Expect(utils.BytesToBool([]byte{0, 0, 0})).To(BeFalse())
			Expect(utils.BytesToBool([]byte{})).To(BeFalse())
		})

		It("should roundtrip bool values", func() {
			Expect(utils.BytesToBool(utils.BoolToBytes(true))).To(BeTrue())
			Expect(utils.BytesToBool(utils.BoolToBytes(false))).To(BeFalse())
		})
	})

	// ==========================================================================
	// Hex ↔ Bytes
	// ==========================================================================
	Context("Hex conversions", func() {
		It("should convert bytes to hex with prefix", func() {
			testCases := []struct {
				bytes    []byte
				expected string
			}{
				{[]byte{}, "0x"},
				{[]byte{0}, "0x00"},
				{[]byte{1}, "0x01"},
				{[]byte{255}, "0xff"},
				{[]byte{0xde, 0xad, 0xbe, 0xef}, "0xdeadbeef"},
				{[]byte{0x00, 0x01, 0x02, 0x03}, "0x00010203"},
			}

			for _, tc := range testCases {
				result := utils.BytesToHex(tc.bytes)
				Expect(result).To(Equal(tc.expected), "BytesToHex(%v)", tc.bytes)
			}
		})

		It("should convert bytes to hex without prefix", func() {
			Expect(utils.BytesToHexUnprefixed([]byte{0xde, 0xad})).To(Equal("dead"))
			Expect(utils.BytesToHexUnprefixed([]byte{0x00, 0x01})).To(Equal("0001"))
		})

		It("should convert hex to bytes", func() {
			testCases := []struct {
				hex      string
				expected []byte
			}{
				{"0x", []byte{}},
				{"0x00", []byte{0}},
				{"0x01", []byte{1}},
				{"0xff", []byte{255}},
				{"0xdeadbeef", []byte{0xde, 0xad, 0xbe, 0xef}},
				{"deadbeef", []byte{0xde, 0xad, 0xbe, 0xef}},   // without prefix
				{"0xDEADBEEF", []byte{0xde, 0xad, 0xbe, 0xef}}, // uppercase
			}

			for _, tc := range testCases {
				result, err := utils.HexToBytes(tc.hex)
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal(tc.expected), "HexToBytes(%s)", tc.hex)
			}
		})

		It("should roundtrip hex values", func() {
			original := []byte{0xca, 0xfe, 0xba, 0xbe}
			hex := utils.BytesToHex(original)
			result, err := utils.HexToBytes(hex)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(original))
		})

		It("should return error for invalid hex", func() {
			_, err := utils.HexToBytes("0xGG")
			Expect(err).To(HaveOccurred())
		})

		It("should handle odd-length hex by padding", func() {
			// "0x123" (odd length) should be padded to "0x0123"
			result, err := utils.HexToBytes("0x123")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal([]byte{0x01, 0x23}))
		})
	})

	// ==========================================================================
	// Fluent API (ByteConverter)
	// ==========================================================================
	Context("Fluent API", func() {
		It("should chain conversions with FromBytes", func() {
			data := []byte{0, 0, 0, 0, 0, 0, 1, 0} // 256

			converter := utils.FromBytes(data)

			Expect(converter.ToInt()).To(Equal(int64(256)))
			Expect(converter.ToUint()).To(Equal(uint64(256)))
			Expect(converter.ToBigInt().Int64()).To(Equal(int64(256)))
			Expect(converter.ToBool()).To(BeTrue())
			Expect(converter.ToHex()).To(Equal("0x0000000000000100"))
			Expect(converter.ToBytes()).To(Equal(data))
		})

		It("should handle zero values", func() {
			converter := utils.FromBytes([]byte{0, 0, 0, 0})

			Expect(converter.ToInt()).To(Equal(int64(0)))
			Expect(converter.ToBool()).To(BeFalse())
		})
	})
})
