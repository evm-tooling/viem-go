package abi_test

import (
	"github.com/ChefBingbong/viem-go/abi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ABI Selector", func() {
	Context("when computing function selectors", func() {
		It("should compute transfer selector correctly", func() {
			selector := abi.ComputeSelector("transfer(address,uint256)")
			// Known selector for ERC20 transfer
			Expect(selector).To(Equal([4]byte{0xa9, 0x05, 0x9c, 0xbb}))
		})

		It("should compute balanceOf selector correctly", func() {
			selector := abi.ComputeSelector("balanceOf(address)")
			// Known selector for ERC20 balanceOf
			Expect(selector).To(Equal([4]byte{0x70, 0xa0, 0x82, 0x31}))
		})

		It("should compute approve selector correctly", func() {
			selector := abi.ComputeSelector("approve(address,uint256)")
			// Known selector for ERC20 approve
			Expect(selector).To(Equal([4]byte{0x09, 0x5e, 0xa7, 0xb3}))
		})
	})

	Context("when computing event topics", func() {
		It("should compute Transfer event topic correctly", func() {
			topic := abi.ComputeEventTopic("Transfer(address,address,uint256)")
			expectedHex := "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
			Expect(topic.Hex()).To(Equal(expectedHex))
		})

		It("should compute Approval event topic correctly", func() {
			topic := abi.ComputeEventTopic("Approval(address,address,uint256)")
			expectedHex := "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"
			Expect(topic.Hex()).To(Equal(expectedHex))
		})
	})

	Context("when converting selectors to hex", func() {
		It("should convert selector to hex string", func() {
			selector := [4]byte{0xa9, 0x05, 0x9c, 0xbb}
			hex := abi.SelectorToHex(selector)
			Expect(hex).To(Equal("0xa9059cbb"))
		})
	})

	Context("when converting hex to selector", func() {
		It("should convert hex string with 0x prefix", func() {
			selector, err := abi.HexToSelector("0xa9059cbb")
			Expect(err).ToNot(HaveOccurred())
			Expect(selector).To(Equal([4]byte{0xa9, 0x05, 0x9c, 0xbb}))
		})

		It("should convert hex string without prefix", func() {
			selector, err := abi.HexToSelector("a9059cbb")
			Expect(err).ToNot(HaveOccurred())
			Expect(selector).To(Equal([4]byte{0xa9, 0x05, 0x9c, 0xbb}))
		})

		It("should return error for invalid length", func() {
			_, err := abi.HexToSelector("a905")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when building function signatures", func() {
		It("should build signature correctly", func() {
			sig := abi.BuildFunctionSignature("transfer", []string{"address", "uint256"})
			Expect(sig).To(Equal("transfer(address,uint256)"))
		})

		It("should handle no parameters", func() {
			sig := abi.BuildFunctionSignature("totalSupply", nil)
			Expect(sig).To(Equal("totalSupply()"))
		})
	})

	Context("when parsing function signatures", func() {
		It("should parse simple signature", func() {
			name, types, err := abi.ParseFunctionSignature("transfer(address,uint256)")
			Expect(err).ToNot(HaveOccurred())
			Expect(name).To(Equal("transfer"))
			Expect(types).To(Equal([]string{"address", "uint256"}))
		})

		It("should parse signature with no params", func() {
			name, types, err := abi.ParseFunctionSignature("totalSupply()")
			Expect(err).ToNot(HaveOccurred())
			Expect(name).To(Equal("totalSupply"))
			Expect(types).To(BeEmpty())
		})

		It("should parse signature with tuple", func() {
			name, types, err := abi.ParseFunctionSignature("foo((address,uint256),bool)")
			Expect(err).ToNot(HaveOccurred())
			Expect(name).To(Equal("foo"))
			Expect(types).To(Equal([]string{"(address,uint256)", "bool"}))
		})
	})

	Context("when checking standard selectors", func() {
		It("should identify transfer selector", func() {
			selector := abi.ComputeSelector("transfer(address,uint256)")
			name, found := abi.IsStandardSelector(selector)
			Expect(found).To(BeTrue())
			Expect(name).To(Equal("transfer"))
		})

		It("should identify balanceOf selector", func() {
			selector := abi.ComputeSelector("balanceOf(address)")
			name, found := abi.IsStandardSelector(selector)
			Expect(found).To(BeTrue())
			Expect(name).To(Equal("balanceOf"))
		})
	})
})
