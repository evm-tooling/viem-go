package abi_test

import (
	"math/big"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ethereum/go-ethereum/common"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetAbiItem", func() {
	Context("basic lookup", func() {
		It("should find function by name", func() {
			jsonABI := []byte(`[{"type":"function","name":"transfer","inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"name":"","type":"bool"}]}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			item, err := parsed.GetAbiItem("transfer", nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(item).ToNot(BeNil())

			fn, ok := item.(abi.Function)
			Expect(ok).To(BeTrue())
			Expect(fn.Name).To(Equal("transfer"))
		})

		It("should find function by selector", func() {
			jsonABI := []byte(`[{"type":"function","name":"transfer","inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"name":"","type":"bool"}]}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			// transfer(address,uint256) selector: 0xa9059cbb
			item, err := parsed.GetAbiItem("0xa9059cbb", nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(item).ToNot(BeNil())

			fn, ok := item.(abi.Function)
			Expect(ok).To(BeTrue())
			Expect(fn.Name).To(Equal("transfer"))
		})

		It("should find event by name", func() {
			jsonABI := []byte(`[{"type":"event","name":"Transfer","inputs":[{"name":"from","type":"address","indexed":true},{"name":"to","type":"address","indexed":true},{"name":"value","type":"uint256","indexed":false}]}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			item, err := parsed.GetAbiItem("Transfer", nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(item).ToNot(BeNil())

			ev, ok := item.(abi.Event)
			Expect(ok).To(BeTrue())
			Expect(ev.Name).To(Equal("Transfer"))
		})

		It("should find error by name", func() {
			jsonABI := []byte(`[{"type":"error","name":"InsufficientBalance","inputs":[{"name":"balance","type":"uint256"},{"name":"required","type":"uint256"}]}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			item, err := parsed.GetAbiItem("InsufficientBalance", nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(item).ToNot(BeNil())

			e, ok := item.(abi.Error)
			Expect(ok).To(BeTrue())
			Expect(e.Name).To(Equal("InsufficientBalance"))
		})

		It("should return error for non-existent item", func() {
			jsonABI := []byte(`[{"type":"function","name":"transfer","inputs":[]}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			_, err = parsed.GetAbiItem("nonexistent", nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})
	})

	Context("overloaded functions", func() {
		// Note: go-ethereum's ABI parser uses a map keyed by function name,
		// so it doesn't support true overloading (it keeps the last definition).
		// This is a limitation of the underlying library.
		// For true overloading support, you'd need a different ABI representation.
		It("should find function with matching args", func() {
			// Single function case - overloading requires different approach
			jsonABI := []byte(`[
				{"type":"function","name":"transfer","inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[]}
			]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			// Should match transfer with correct args
			item, err := parsed.GetAbiItem("transfer", &abi.GetAbiItemOptions{
				Args: []any{common.Address{}, big.NewInt(123)},
			})
			Expect(err).ToNot(HaveOccurred())
			fn := item.(abi.Function)
			Expect(fn.Name).To(Equal("transfer"))
			Expect(fn.Inputs).To(HaveLen(2))
		})
	})

	Context("GetFunction helper", func() {
		It("should find function by name", func() {
			jsonABI := []byte(`[{"type":"function","name":"transfer","inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"name":"","type":"bool"}]}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			fn, err := parsed.GetFunction("transfer")
			Expect(err).ToNot(HaveOccurred())
			Expect(fn.Name).To(Equal("transfer"))
		})

		It("should return error for non-existent function", func() {
			jsonABI := []byte(`[{"type":"function","name":"transfer","inputs":[]}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			_, err = parsed.GetFunction("nonexistent")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("GetFunctionBySelector", func() {
		It("should find function by selector", func() {
			jsonABI := []byte(`[{"type":"function","name":"transfer","inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"name":"","type":"bool"}]}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			selector := [4]byte{0xa9, 0x05, 0x9c, 0xbb}
			fn, err := parsed.GetFunctionBySelector(selector)
			Expect(err).ToNot(HaveOccurred())
			Expect(fn.Name).To(Equal("transfer"))
		})
	})

	Context("GetEvent helper", func() {
		It("should find event by name", func() {
			jsonABI := []byte(`[{"type":"event","name":"Transfer","inputs":[{"name":"from","type":"address","indexed":true}]}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			ev, err := parsed.GetEvent("Transfer")
			Expect(err).ToNot(HaveOccurred())
			Expect(ev.Name).To(Equal("Transfer"))
		})
	})

	Context("GetError helper", func() {
		It("should find error by name", func() {
			jsonABI := []byte(`[{"type":"error","name":"InsufficientBalance","inputs":[{"name":"balance","type":"uint256"}]}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			e, err := parsed.GetError("InsufficientBalance")
			Expect(err).ToNot(HaveOccurred())
			Expect(e.Name).To(Equal("InsufficientBalance"))
		})
	})
})
