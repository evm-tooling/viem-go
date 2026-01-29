package abi_test

import (
	"github.com/ChefBingbong/viem-go/abi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ABI Parser", func() {
	var erc20ABI = []byte(`[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`)

	Context("when parsing a valid ABI", func() {
		It("should parse without error", func() {
			parsed, err := abi.Parse(erc20ABI)
			Expect(err).ToNot(HaveOccurred())
			Expect(parsed).ToNot(BeNil())
		})

		It("should contain the expected functions", func() {
			parsed, err := abi.Parse(erc20ABI)
			Expect(err).ToNot(HaveOccurred())

			Expect(parsed.HasFunction("name")).To(BeTrue())
			Expect(parsed.HasFunction("symbol")).To(BeTrue())
			Expect(parsed.HasFunction("decimals")).To(BeTrue())
			Expect(parsed.HasFunction("totalSupply")).To(BeTrue())
			Expect(parsed.HasFunction("balanceOf")).To(BeTrue())
			Expect(parsed.HasFunction("transfer")).To(BeTrue())
		})

		It("should contain the expected events", func() {
			parsed, err := abi.Parse(erc20ABI)
			Expect(err).ToNot(HaveOccurred())

			Expect(parsed.HasEvent("Transfer")).To(BeTrue())
		})

		It("should return function names", func() {
			parsed, err := abi.Parse(erc20ABI)
			Expect(err).ToNot(HaveOccurred())

			names := parsed.FunctionNames()
			Expect(names).To(ContainElement("name"))
			Expect(names).To(ContainElement("transfer"))
		})
	})

	Context("when getting a function", func() {
		It("should return function details", func() {
			parsed, err := abi.Parse(erc20ABI)
			Expect(err).ToNot(HaveOccurred())

			fn, err := parsed.GetFunction("transfer")
			Expect(err).ToNot(HaveOccurred())
			Expect(fn.Name).To(Equal("transfer"))
			Expect(len(fn.Inputs)).To(Equal(2))
			Expect(fn.Inputs[0].Type).To(Equal("address"))
			Expect(fn.Inputs[1].Type).To(Equal("uint256"))
		})

		It("should return error for non-existent function", func() {
			parsed, err := abi.Parse(erc20ABI)
			Expect(err).ToNot(HaveOccurred())

			_, err = parsed.GetFunction("nonexistent")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when parsing an invalid ABI", func() {
		It("should return an error", func() {
			_, err := abi.Parse([]byte(`invalid json`))
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when using MustParse", func() {
		It("should panic on invalid ABI", func() {
			Expect(func() {
				abi.MustParse([]byte(`invalid`))
			}).To(Panic())
		})

		It("should not panic on valid ABI", func() {
			Expect(func() {
				abi.MustParse(erc20ABI)
			}).ToNot(Panic())
		})
	})
})
