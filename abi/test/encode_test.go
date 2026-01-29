package abi_test

import (
	"math/big"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ethereum/go-ethereum/common"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ABI Encode", func() {
	var erc20ABI = []byte(`[{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"type":"function"}]`)

	Context("when encoding a function call", func() {
		It("should encode transfer call correctly", func() {
			parsed, err := abi.Parse(erc20ABI)
			Expect(err).ToNot(HaveOccurred())

			to := common.HexToAddress("0x1234567890123456789012345678901234567890")
			amount := big.NewInt(1000)

			calldata, err := parsed.EncodeCall("transfer", to, amount)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(calldata)).To(BeNumerically(">", 4))

			// Check selector
			Expect(calldata[:4]).To(Equal([]byte{0xa9, 0x05, 0x9c, 0xbb}))
		})

		It("should return error for non-existent method", func() {
			parsed, err := abi.Parse(erc20ABI)
			Expect(err).ToNot(HaveOccurred())

			_, err = parsed.EncodeCall("nonexistent", common.Address{}, big.NewInt(0))
			Expect(err).To(HaveOccurred())
		})

		It("should return error for wrong argument types", func() {
			parsed, err := abi.Parse(erc20ABI)
			Expect(err).ToNot(HaveOccurred())

			// Passing wrong types
			_, err = parsed.EncodeCall("transfer", "not an address", "not a number")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when encoding arguments only", func() {
		It("should encode without selector", func() {
			parsed, err := abi.Parse(erc20ABI)
			Expect(err).ToNot(HaveOccurred())

			to := common.HexToAddress("0x1234567890123456789012345678901234567890")
			amount := big.NewInt(1000)

			args, err := parsed.EncodeArgs("transfer", to, amount)
			Expect(err).ToNot(HaveOccurred())

			// Should be 64 bytes (32 for address, 32 for uint256)
			Expect(len(args)).To(Equal(64))
		})
	})

	Context("when using Pack alias", func() {
		It("should behave same as EncodeCall", func() {
			parsed, err := abi.Parse(erc20ABI)
			Expect(err).ToNot(HaveOccurred())

			to := common.HexToAddress("0x1234567890123456789012345678901234567890")
			amount := big.NewInt(1000)

			calldata1, err := parsed.EncodeCall("transfer", to, amount)
			Expect(err).ToNot(HaveOccurred())

			calldata2, err := parsed.Pack("transfer", to, amount)
			Expect(err).ToNot(HaveOccurred())

			Expect(calldata1).To(Equal(calldata2))
		})
	})
})
