package contract_test

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/contract"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Contract", func() {
	var erc20ABI = []byte(`[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"type":"function"}]`)

	Context("when creating a contract without a client", func() {
		It("should create a contract with valid ABI", func() {
			addr := common.HexToAddress("0x1234567890123456789012345678901234567890")

			// Note: This will work for ABI parsing but client will be nil
			cont, err := contract.NewContract(addr, erc20ABI, nil)
			cont.Read(context.Background(), "balanceOf", common.HexToAddress("0x1234567890123456789012345678901234567890"))
			Expect(err).ToNot(HaveOccurred())
			Expect(cont).ToNot(BeNil())
			Expect(cont.Address()).To(Equal(addr))
		})

		It("should fail with invalid ABI", func() {
			addr := common.HexToAddress("0x1234567890123456789012345678901234567890")

			_, err := contract.NewContract(addr, []byte(`invalid`), nil)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when checking contract functions", func() {
		It("should report available functions", func() {
			addr := common.HexToAddress("0x1234567890123456789012345678901234567890")
			cont, err := contract.NewContract(addr, erc20ABI, nil)
			Expect(err).ToNot(HaveOccurred())

			Expect(cont.HasFunction("name")).To(BeTrue())
			Expect(cont.HasFunction("transfer")).To(BeTrue())
			Expect(cont.HasFunction("nonexistent")).To(BeFalse())
		})

		It("should list function names", func() {
			addr := common.HexToAddress("0x1234567890123456789012345678901234567890")
			cont, err := contract.NewContract(addr, erc20ABI, nil)
			Expect(err).ToNot(HaveOccurred())

			names := cont.FunctionNames()
			Expect(names).To(ContainElement("name"))
			Expect(names).To(ContainElement("symbol"))
			Expect(names).To(ContainElement("transfer"))
		})

		It("should get function details", func() {
			addr := common.HexToAddress("0x1234567890123456789012345678901234567890")
			cont, err := contract.NewContract(addr, erc20ABI, nil)
			Expect(err).ToNot(HaveOccurred())

			fn, err := cont.GetFunction("transfer")
			Expect(err).ToNot(HaveOccurred())
			Expect(fn.Name).To(Equal("transfer"))
			Expect(len(fn.Inputs)).To(Equal(2))
		})
	})

	Context("when encoding calls", func() {
		It("should encode function calls", func() {
			addr := common.HexToAddress("0x1234567890123456789012345678901234567890")
			cont, err := contract.NewContract(addr, erc20ABI, nil)
			Expect(err).ToNot(HaveOccurred())

			calldata, err := cont.EncodeCall("balanceOf", common.HexToAddress("0x0"))
			Expect(err).ToNot(HaveOccurred())
			Expect(len(calldata)).To(BeNumerically(">", 4))

			// Check selector for balanceOf(address)
			Expect(calldata[:4]).To(Equal([]byte{0x70, 0xa0, 0x82, 0x31}))
		})
	})

	Context("when cloning a contract", func() {
		It("should create a copy with new address", func() {
			addr1 := common.HexToAddress("0x1111111111111111111111111111111111111111")
			addr2 := common.HexToAddress("0x2222222222222222222222222222222222222222")

			cont, err := contract.NewContract(addr1, erc20ABI, nil)
			Expect(err).ToNot(HaveOccurred())

			cloned := cont.Clone(addr2)
			Expect(cloned.Address()).To(Equal(addr2))
			Expect(cont.Address()).To(Equal(addr1)) // Original unchanged
		})
	})

	Context("when using MustNewContract", func() {
		It("should panic on invalid ABI", func() {
			addr := common.HexToAddress("0x1234567890123456789012345678901234567890")

			Expect(func() {
				contract.MustNewContract(addr, []byte(`invalid`), nil)
			}).To(Panic())
		})

		It("should not panic on valid ABI", func() {
			addr := common.HexToAddress("0x1234567890123456789012345678901234567890")

			Expect(func() {
				contract.MustNewContract(addr, erc20ABI, nil)
			}).ToNot(Panic())
		})
	})
})
