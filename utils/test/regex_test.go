package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ChefBingbong/viem-go/utils"
)

var _ = Describe("Regex", func() {
	Describe("ArrayRegex", func() {
		It("should match dynamic array types", func() {
			Expect(utils.IsArrayType("uint256[]")).To(BeTrue())
			Expect(utils.IsArrayType("bytes32[]")).To(BeTrue())
			Expect(utils.IsArrayType("address[]")).To(BeTrue())
		})

		It("should match fixed-size array types", func() {
			Expect(utils.IsArrayType("uint256[10]")).To(BeTrue())
			Expect(utils.IsArrayType("bytes32[5]")).To(BeTrue())
			Expect(utils.IsArrayType("address[100]")).To(BeTrue())
		})

		It("should not match non-array types", func() {
			Expect(utils.IsArrayType("uint256")).To(BeFalse())
			Expect(utils.IsArrayType("bytes32")).To(BeFalse())
			Expect(utils.IsArrayType("address")).To(BeFalse())
		})
	})

	Describe("BytesRegex", func() {
		It("should match dynamic bytes", func() {
			Expect(utils.IsBytesType("bytes")).To(BeTrue())
		})

		It("should match fixed bytes types", func() {
			Expect(utils.IsBytesType("bytes1")).To(BeTrue())
			Expect(utils.IsBytesType("bytes16")).To(BeTrue())
			Expect(utils.IsBytesType("bytes32")).To(BeTrue())
		})

		It("should not match invalid bytes types", func() {
			Expect(utils.IsBytesType("bytes0")).To(BeFalse())
			Expect(utils.IsBytesType("bytes33")).To(BeFalse())
			Expect(utils.IsBytesType("bytes256")).To(BeFalse())
		})
	})

	Describe("IntegerRegex", func() {
		It("should match uint types", func() {
			Expect(utils.IsIntegerType("uint")).To(BeTrue())
			Expect(utils.IsIntegerType("uint8")).To(BeTrue())
			Expect(utils.IsIntegerType("uint256")).To(BeTrue())
		})

		It("should match int types", func() {
			Expect(utils.IsIntegerType("int")).To(BeTrue())
			Expect(utils.IsIntegerType("int8")).To(BeTrue())
			Expect(utils.IsIntegerType("int256")).To(BeTrue())
		})

		It("should not match invalid sizes", func() {
			Expect(utils.IsIntegerType("uint7")).To(BeFalse())
			Expect(utils.IsIntegerType("uint300")).To(BeFalse())
			Expect(utils.IsIntegerType("int15")).To(BeFalse())
		})
	})

	Describe("ParseArrayType", func() {
		It("should parse dynamic arrays", func() {
			base, size := utils.ParseArrayType("uint256[]")
			Expect(base).To(Equal("uint256"))
			Expect(size).To(Equal(""))
		})

		It("should parse fixed-size arrays", func() {
			base, size := utils.ParseArrayType("bytes32[10]")
			Expect(base).To(Equal("bytes32"))
			Expect(size).To(Equal("10"))
		})

		It("should return empty for non-arrays", func() {
			base, size := utils.ParseArrayType("uint256")
			Expect(base).To(Equal(""))
			Expect(size).To(Equal(""))
		})
	})

	Describe("ParseIntegerType", func() {
		It("should parse uint types", func() {
			unsigned, bits := utils.ParseIntegerType("uint256")
			Expect(unsigned).To(BeTrue())
			Expect(bits).To(Equal(256))
		})

		It("should parse int types", func() {
			unsigned, bits := utils.ParseIntegerType("int128")
			Expect(unsigned).To(BeFalse())
			Expect(bits).To(Equal(128))
		})

		It("should default to 256 bits", func() {
			unsigned, bits := utils.ParseIntegerType("uint")
			Expect(unsigned).To(BeTrue())
			Expect(bits).To(Equal(256))
		})

		It("should return 0 for invalid types", func() {
			_, bits := utils.ParseIntegerType("address")
			Expect(bits).To(Equal(0))
		})
	})
})
