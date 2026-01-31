package utils_test

import (
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ChefBingbong/viem-go/utils"
)

var _ = Describe("Stringify", func() {
	Describe("Stringify", func() {
		It("should stringify simple values", func() {
			result, err := utils.Stringify(map[string]any{
				"name":  "test",
				"value": 42,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring(`"name":"test"`))
			Expect(result).To(ContainSubstring(`"value":42`))
		})

		It("should convert *big.Int to string", func() {
			result, err := utils.Stringify(map[string]any{
				"balance": big.NewInt(123456789012345),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring(`"balance":"123456789012345"`))
		})

		It("should handle nested big.Int values", func() {
			result, err := utils.Stringify(map[string]any{
				"outer": map[string]any{
					"inner": big.NewInt(999),
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring(`"inner":"999"`))
		})

		It("should handle arrays with big.Int", func() {
			result, err := utils.Stringify([]any{
				big.NewInt(1),
				big.NewInt(2),
				big.NewInt(3),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(`["1","2","3"]`))
		})

		It("should handle nil values", func() {
			result, err := utils.Stringify(map[string]any{
				"value": nil,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring(`"value":null`))
		})

		It("should handle structs with big.Int fields", func() {
			type TestStruct struct {
				Name    string   `json:"name"`
				Balance *big.Int `json:"balance"`
			}

			result, err := utils.Stringify(TestStruct{
				Name:    "test",
				Balance: big.NewInt(1000),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring(`"name":"test"`))
			Expect(result).To(ContainSubstring(`"balance":"1000"`))
		})
	})

	Describe("StringifyIndent", func() {
		It("should stringify with indentation", func() {
			result, err := utils.StringifyIndent(map[string]any{
				"value": big.NewInt(123),
			}, "", "  ")
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("\n"))
			Expect(result).To(ContainSubstring(`"value": "123"`))
		})
	})
})
