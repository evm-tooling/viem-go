package abi_test

import (
	"math/big"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ethereum/go-ethereum/common"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FormatAbiItemWithArgs", func() {
	Context("basic formatting", func() {
		It("should format function with address and uint256 args", func() {
			fn := abi.Function{
				Name: "transfer",
				Inputs: []abi.Parameter{
					{Name: "to", Type: "address"},
					{Name: "amount", Type: "uint256"},
				},
			}

			args := []any{
				common.HexToAddress("0x1234567890123456789012345678901234567890"),
				big.NewInt(1000000000000000000),
			}

			result := abi.FormatAbiItemWithArgs(fn, args, nil)
			Expect(result).To(Equal("transfer(0x1234567890123456789012345678901234567890, 1000000000000000000)"))
		})

		It("should format with parameter names", func() {
			fn := abi.Function{
				Name: "transfer",
				Inputs: []abi.Parameter{
					{Name: "to", Type: "address"},
					{Name: "amount", Type: "uint256"},
				},
			}

			args := []any{
				common.HexToAddress("0x1234567890123456789012345678901234567890"),
				big.NewInt(1000),
			}

			opts := &abi.FormatAbiItemWithArgsOptions{
				IncludeFunctionName: true,
				IncludeName:         true,
			}

			result := abi.FormatAbiItemWithArgs(fn, args, opts)
			Expect(result).To(Equal("transfer(to: 0x1234567890123456789012345678901234567890, amount: 1000)"))
		})

		It("should format without function name", func() {
			fn := abi.Function{
				Name: "transfer",
				Inputs: []abi.Parameter{
					{Name: "to", Type: "address"},
				},
			}

			args := []any{
				common.HexToAddress("0x1234567890123456789012345678901234567890"),
			}

			opts := &abi.FormatAbiItemWithArgsOptions{
				IncludeFunctionName: false,
				IncludeName:         false,
			}

			result := abi.FormatAbiItemWithArgs(fn, args, opts)
			Expect(result).To(Equal("(0x1234567890123456789012345678901234567890)"))
		})

		It("should format event", func() {
			ev := abi.Event{
				Name: "Transfer",
				Inputs: []abi.Parameter{
					{Name: "from", Type: "address"},
					{Name: "to", Type: "address"},
					{Name: "value", Type: "uint256"},
				},
			}

			args := []any{
				common.HexToAddress("0x1111111111111111111111111111111111111111"),
				common.HexToAddress("0x2222222222222222222222222222222222222222"),
				big.NewInt(1000),
			}

			result := abi.FormatAbiItemWithArgs(ev, args, nil)
			Expect(result).To(ContainSubstring("Transfer"))
			Expect(result).To(ContainSubstring("0x1111111111111111111111111111111111111111"))
			Expect(result).To(ContainSubstring("0x2222222222222222222222222222222222222222"))
			Expect(result).To(ContainSubstring("1000"))
		})

		It("should format error", func() {
			e := abi.Error{
				Name: "InsufficientBalance",
				Inputs: []abi.Parameter{
					{Name: "balance", Type: "uint256"},
					{Name: "required", Type: "uint256"},
				},
			}

			args := []any{
				big.NewInt(100),
				big.NewInt(1000),
			}

			result := abi.FormatAbiItemWithArgs(e, args, nil)
			Expect(result).To(Equal("InsufficientBalance(100, 1000)"))
		})

		It("should format with bool args", func() {
			fn := abi.Function{
				Name: "setApproval",
				Inputs: []abi.Parameter{
					{Name: "approved", Type: "bool"},
				},
			}

			result1 := abi.FormatAbiItemWithArgs(fn, []any{true}, nil)
			Expect(result1).To(Equal("setApproval(true)"))

			result2 := abi.FormatAbiItemWithArgs(fn, []any{false}, nil)
			Expect(result2).To(Equal("setApproval(false)"))
		})

		It("should format with string args", func() {
			fn := abi.Function{
				Name: "setName",
				Inputs: []abi.Parameter{
					{Name: "name", Type: "string"},
				},
			}

			result := abi.FormatAbiItemWithArgs(fn, []any{"hello world"}, nil)
			Expect(result).To(Equal(`setName("hello world")`))
		})

		It("should format with bytes args", func() {
			fn := abi.Function{
				Name: "execute",
				Inputs: []abi.Parameter{
					{Name: "data", Type: "bytes"},
				},
			}

			result := abi.FormatAbiItemWithArgs(fn, []any{[]byte{0x12, 0x34}}, nil)
			Expect(result).To(Equal("execute(0x1234)"))
		})

		It("should format with nil arg", func() {
			fn := abi.Function{
				Name: "test",
				Inputs: []abi.Parameter{
					{Name: "value", Type: "uint256"},
				},
			}

			result := abi.FormatAbiItemWithArgs(fn, []any{nil}, nil)
			Expect(result).To(Equal("test(null)"))
		})

		It("should format empty function", func() {
			fn := abi.Function{
				Name:   "noArgs",
				Inputs: nil,
			}

			result := abi.FormatAbiItemWithArgs(fn, []any{}, nil)
			Expect(result).To(Equal("noArgs()"))
		})
	})

	Context("using ABI method", func() {
		It("should format via ABI method", func() {
			jsonABI := []byte(`[{"type":"function","name":"transfer","inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}]}]`)
			parsed, err := abi.Parse(jsonABI)
			Expect(err).ToNot(HaveOccurred())

			args := []any{
				common.HexToAddress("0x1234567890123456789012345678901234567890"),
				big.NewInt(1000),
			}

			result, err := parsed.FormatFunctionCallWithArgs("transfer", args, nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(ContainSubstring("transfer"))
			Expect(result).To(ContainSubstring("1000"))
		})
	})
})
