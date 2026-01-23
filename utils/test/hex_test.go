package utils_test

import (
	"math/big"

	"github.com/ChefBingbong/viem-go/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hex Utils", func() {
	// ==========================================================================
	// Hex → Int64
	// ==========================================================================
	Context("Hex to Int64 conversions", func() {
		It("should convert hex to int64", func() {
			testCases := []struct {
				hex      string
				expected int64
			}{
				{"0x0", 0},
				{"0x1", 1},
				{"0xff", 255},
				{"0xFF", 255}, // uppercase
				{"0x100", 256},
				{"0xffff", 65535},
				{"0xf4240", 1000000},
				{"0x7fffffffffffffff", 9223372036854775807}, // max int64
				{"ff", 255},                                 // without prefix
			}

			for _, tc := range testCases {
				result, err := utils.HexToInt(tc.hex)
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal(tc.expected), "HexToInt(%s)", tc.hex)
			}
		})

		It("should convert int64 to hex", func() {
			testCases := []struct {
				value    int64
				expected string
			}{
				{0, "0x0"},
				{1, "0x1"},
				{255, "0xff"},
				{256, "0x100"},
				{65535, "0xffff"},
				{1000000, "0xf4240"},
				{-1, "-0x1"},
				{-255, "-0xff"},
			}

			for _, tc := range testCases {
				result := utils.IntToHex(tc.value)
				Expect(result).To(Equal(tc.expected), "IntToHex(%d)", tc.value)
			}
		})

		It("should roundtrip int64 values", func() {
			values := []int64{0, 1, 255, 256, 65535, 1000000, 9223372036854775807}
			for _, v := range values {
				hex := utils.IntToHex(v)
				result, err := utils.HexToInt(hex)
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal(v))
			}
		})
	})

	// ==========================================================================
	// Hex → Uint64
	// ==========================================================================
	Context("Hex to Uint64 conversions", func() {
		It("should convert hex to uint64", func() {
			testCases := []struct {
				hex      string
				expected uint64
			}{
				{"0x0", 0},
				{"0x1", 1},
				{"0xff", 255},
				{"0xffffffffffffffff", 18446744073709551615}, // max uint64
			}

			for _, tc := range testCases {
				result, err := utils.HexToUint(tc.hex)
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal(tc.expected), "HexToUint(%s)", tc.hex)
			}
		})

		It("should convert uint64 to hex", func() {
			testCases := []struct {
				value    uint64
				expected string
			}{
				{0, "0x0"},
				{1, "0x1"},
				{255, "0xff"},
				{18446744073709551615, "0xffffffffffffffff"},
			}

			for _, tc := range testCases {
				result := utils.UintToHex(tc.value)
				Expect(result).To(Equal(tc.expected), "UintToHex(%d)", tc.value)
			}
		})
	})

	// ==========================================================================
	// Hex → BigInt
	// ==========================================================================
	Context("Hex to BigInt conversions", func() {
		It("should convert hex to big.Int", func() {
			testCases := []struct {
				hex      string
				expected string // decimal representation
			}{
				{"0x0", "0"},
				{"0x1", "1"},
				{"0xff", "255"},
				{"0xde0b6b3a7640000", "1000000000000000000"},                                                           // 1 ETH in wei
				{"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "115792089237316195423570985008687907853269984665640564039457584007913129639935"}, // max uint256
			}

			for _, tc := range testCases {
				result, err := utils.HexToBigInt(tc.hex)
				Expect(err).ToNot(HaveOccurred())
				Expect(result.String()).To(Equal(tc.expected), "HexToBigInt(%s)", tc.hex)
			}
		})

		It("should convert big.Int to hex", func() {
			testCases := []struct {
				decimal  string
				expected string
			}{
				{"0", "0x0"},
				{"1", "0x1"},
				{"255", "0xff"},
				{"256", "0x100"},
				{"1000000000000000000", "0xde0b6b3a7640000"}, // 1 ETH in wei
			}

			for _, tc := range testCases {
				value, _ := new(big.Int).SetString(tc.decimal, 10)
				result := utils.BigIntToHex(value)
				Expect(result).To(Equal(tc.expected), "BigIntToHex(%s)", tc.decimal)
			}
		})

		It("should handle nil big.Int", func() {
			result := utils.BigIntToHex(nil)
			Expect(result).To(Equal("0x0"))
		})

		It("should return error for invalid hex", func() {
			_, err := utils.HexToBigInt("0xGGGG")
			Expect(err).To(HaveOccurred())
		})
	})

	// ==========================================================================
	// Hex → Bool
	// ==========================================================================
	Context("Hex to Bool conversions", func() {
		It("should convert hex to bool", func() {
			// True cases
			result, err := utils.HexToBool("0x1")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeTrue())

			result, err = utils.HexToBool("0xff")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeTrue())

			result, err = utils.HexToBool("0x0001")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeTrue())

			// False cases
			result, err = utils.HexToBool("0x0")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeFalse())

			result, err = utils.HexToBool("0x00")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeFalse())

			result, err = utils.HexToBool("0x0000")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeFalse())
		})

		It("should convert bool to hex", func() {
			Expect(utils.BoolToHex(true)).To(Equal("0x1"))
			Expect(utils.BoolToHex(false)).To(Equal("0x0"))
		})
	})

	// ==========================================================================
	// Hex → Bytes
	// ==========================================================================
	Context("Hex to Bytes conversions", func() {
		It("should handle odd-length hex strings", func() {
			// "0x1" should become [0x01]
			result, err := utils.HexToBytes("0x1")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal([]byte{0x01}))

			// "0xfff" should become [0x0f, 0xff]
			result, err = utils.HexToBytes("0xfff")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal([]byte{0x0f, 0xff}))
		})
	})

	// ==========================================================================
	// Helper Functions
	// ==========================================================================
	Context("Helper functions", func() {
		It("should validate hex strings", func() {
			// Valid
			Expect(utils.IsValidHex("0x0")).To(BeTrue())
			Expect(utils.IsValidHex("0xff")).To(BeTrue())
			Expect(utils.IsValidHex("0xDeadBeef")).To(BeTrue())
			Expect(utils.IsValidHex("deadbeef")).To(BeTrue())

			// Invalid
			Expect(utils.IsValidHex("")).To(BeFalse())
			Expect(utils.IsValidHex("0x")).To(BeFalse())
			Expect(utils.IsValidHex("0xGG")).To(BeFalse())
			Expect(utils.IsValidHex("hello")).To(BeFalse())
		})

		It("should pad hex strings", func() {
			// Pad to 4 bytes (8 hex chars)
			Expect(utils.PadHex("0x1", 4)).To(Equal("0x00000001"))
			Expect(utils.PadHex("0xff", 4)).To(Equal("0x000000ff"))
			Expect(utils.PadHex("0x12345678", 4)).To(Equal("0x12345678"))

			// Pad to 32 bytes (64 hex chars) - common for Ethereum
			padded := utils.PadHex("0x1", 32)
			Expect(len(padded)).To(Equal(66)) // "0x" + 64 chars
			Expect(padded[len(padded)-1:]).To(Equal("1"))
		})
	})

	// ==========================================================================
	// Fluent API (HexConverter)
	// ==========================================================================
	Context("Fluent API", func() {
		It("should chain conversions with FromHex", func() {
			converter := utils.FromHex("0x100") // 256

			intVal, err := converter.ToInt()
			Expect(err).ToNot(HaveOccurred())
			Expect(intVal).To(Equal(int64(256)))

			uintVal, err := converter.ToUint()
			Expect(err).ToNot(HaveOccurred())
			Expect(uintVal).To(Equal(uint64(256)))

			bigVal, err := converter.ToBigInt()
			Expect(err).ToNot(HaveOccurred())
			Expect(bigVal.Int64()).To(Equal(int64(256)))

			boolVal, err := converter.ToBool()
			Expect(err).ToNot(HaveOccurred())
			Expect(boolVal).To(BeTrue())

			bytesVal, err := converter.ToBytes()
			Expect(err).ToNot(HaveOccurred())
			Expect(bytesVal).To(Equal([]byte{0x01, 0x00}))

			Expect(converter.String()).To(Equal("0x100"))
		})

		It("should handle zero values", func() {
			converter := utils.FromHex("0x0")

			intVal, _ := converter.ToInt()
			Expect(intVal).To(Equal(int64(0)))

			boolVal, _ := converter.ToBool()
			Expect(boolVal).To(BeFalse())
		})
	})
})
